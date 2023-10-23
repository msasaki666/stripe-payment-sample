package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/caarlos0/env/v9"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
	bsession "github.com/stripe/stripe-go/v76/billingportal/session"
	csession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/msasaki666/backend/internal/renv"
	"github.com/msasaki666/backend/models"
)

type config struct {
	StripeApiKey    string           `env:"STRIPE_API_KEY"`
	FrontEndBaseURL string           `env:"FRONT_END_BASE_URL"`
	GoEnv           renv.Environment `env:"GO_ENV"`
	DB              *sql.DB
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if e, ok := err.(*env.AggregateError); ok {
		for _, er := range e.Errors {
			switch v := er.(type) {
			case env.ParseError:
				// handle it
			case env.NotStructPtrError:
				// handle it
			case env.NoParserError:
				// handle it
			case env.NoSupportedTagOptionError:
				// handle it
			default:
				fmt.Printf("Unknown error type %v", v)
			}
		}
	}
	stripe.Key = cfg.StripeApiKey

	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=app_database sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	cfg.DB = db
	e := echo.New()
	corsConfig, err := createCorsConfig(&cfg)
	if err != nil {
		panic(err.Error())
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/customers", func(c echo.Context) error {
		customers, err := models.Customers().All(context.Background(), db)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		return c.JSON(http.StatusOK, customers)
	})
	e.GET("/products", generateGetProductsHandler(&cfg))
	e.POST("/products/create-checkout-session", generateCreateCheckoutSessionHandler(&cfg))
	e.POST("/create-customer-portal-session", generateCreateCustomerPortalSessionHandler(&cfg))
	e.Any("/webhook", echo.WrapHandler(http.HandlerFunc(handleWebhook)))
	e.Logger.Fatal(e.Start(":1323"))
}

type createCheckoutSessionHandlerParams struct {
	StripePriceID string `form:"stripe_price_id"`
}

func generateCreateCheckoutSessionHandler(cfg *config) echo.HandlerFunc {
	return func(c echo.Context) error {
		p := createCheckoutSessionHandlerParams{}
		if err := c.Bind(&p); err != nil {
			return err
		}

		s, err := createCheckoutSession(cfg, p.StripePriceID, calcQuwntity())

		// TODO: 購入済みの時のハンドリング
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, s.URL)
	}
}

type alreadyPurachasedError struct {
}

func (e *alreadyPurachasedError) Error() string {
	return "already purchased"
}

func createCheckoutSession(cfg *config, stripePriceID string, quantity int64) (*stripe.CheckoutSession, error) {
	successURL, err := url.JoinPath(cfg.FrontEndBaseURL, "success")
	if err != nil {
		return nil, err
	}
	cancelURL, err := url.JoinPath(cfg.FrontEndBaseURL, "cancel")
	if err != nil {
		return nil, err
	}

	price, err := models.
		StripePrices(models.StripePriceWhere.IDOnStripe.EQ(stripePriceID)).
		One(context.Background(), cfg.DB)
	if err != nil {
		return nil, err
	}

	// TODO: cusomer idがnullの場合は、Stripeの顧客を作成し、保存する
	customer, err := models.Customers().One(context.Background(), cfg.DB)
	if err != nil {
		return nil, err
	}

	// NOTE: 既に購入済みか確認する
	subscriptions := subscription.List(&stripe.SubscriptionListParams{
		Customer: stripe.String(customer.IDOnStripe),
		Price:    stripe.String(stripePriceID),
	})
	if subscriptions.Next() {
		return nil, &alreadyPurachasedError{}
	}

	sessionType, err := judgeCheckoutSessionType(price.Type)
	if err != nil {
		return nil, err
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(sessionType),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// 料金IDを指定することで、既存の商品のセッションを開始できる
				Price:    stripe.String(stripePriceID),
				Quantity: stripe.Int64(quantity),
			},
		},
		Customer:   stripe.String(customer.IDOnStripe),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Locale:     stripe.String("auto"),
		ExpiresAt:  stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),
	}

	s, err := csession.New(params)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func generateCreateCustomerPortalSessionHandler(cfg *config) echo.HandlerFunc {
	return func(c echo.Context) error {
		returnURL, err := url.JoinPath(cfg.FrontEndBaseURL, "account")
		if err != nil {
			return err
		}

		// TODO: ログインユーザーから取得する
		customer, err := models.Customers().One(context.Background(), cfg.DB)
		if err != nil {
			return err
		}

		params := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(customer.IDOnStripe),
			ReturnURL: stripe.String(returnURL),
		}
		s, err := bsession.New(params)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, s.URL)
	}
}

func calcQuwntity() int64 {
	return 1
}

func judgeCheckoutSessionType(t string) (string, error) {
	switch t {
	case "one_time":
		return "payment", nil
	case "recurring":
		return "subscription", nil
	default:
		return "", fmt.Errorf("invalid type: %s", t)
	}
}

func generateGetProductsHandler(cfg *config) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := models.StripeProducts(
			qm.Load(models.StripeProductRels.StripePrices),
		).All(context.Background(), cfg.DB)
		if err != nil {
			return err
		}
		type ProductsResponse struct {
			*models.StripeProduct
			Price *models.StripePrice `json:"price"`
		}

		pr := lo.Map(products, func(p *models.StripeProduct, i int) *ProductsResponse {
			price, _ := lo.Find(p.R.StripePrices, func(price *models.StripePrice) bool {
				return price.StripeProductID == p.ID
			})
			return &ProductsResponse{
				StripeProduct: p,
				Price:         price,
			}
		})

		return c.JSON(http.StatusOK, pr)
	}
}
