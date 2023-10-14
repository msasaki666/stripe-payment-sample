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
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/msasaki666/backend/internal/renv"
	"github.com/msasaki666/backend/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type config struct {
	StripeApiKey    string           `env:"STRIPE_API_KEY"`
	FrontEndBaseURL string           `env:"FRONT_END_BASE_URL"`
	GoEnv           renv.Environment `env:"GO_ENV"`
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

	boil.SetDB(db)
	e := echo.New()
	corsConfig, err := createCorsConfig(&cfg)
	if err != nil {
		panic(err.Error())
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/payment_stripes", func(c echo.Context) error {
		p, err := models.PaymentStripes().All(context.Background(), db)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		return c.JSON(http.StatusOK, p)
	})
	e.GET("/products", handleGetProducts)
	e.POST("/products/create-checkout-session", generateCreateCheckoutSessionHandler(&cfg))
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

		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, s.URL)
	}
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
	// TODO: cusomer idがnullの場合は、Stripeの顧客を作成し、保存する
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// 料金IDを指定することで、既存の商品のセッションを開始できる
				Price:    stripe.String(stripePriceID),
				Quantity: stripe.Int64(quantity),
			},
		},
		// TODO: 追加する。ユーザーモデルから取得する
		// Customer:  stripe.String("cus_J0Z0Z0Z0Z0Z0Z0Z0Z0Z0Z0Z0Z"),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Locale:     stripe.String("auto"),
		ExpiresAt:  stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),
	}

	s, err := session.New(params)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func calcQuwntity() int64 {
	return 1
}

type ProductSample struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	StripePriceID string `json:"stripe_price_id"`
}

func handleGetProducts(c echo.Context) error {
	// TODO: レコード例
	products := []ProductSample{
		{
			ID:            1,
			Name:          "ドロップイン",
			StripePriceID: "price_1NveyfAj9ehS6HaZQ1SPdae2",
		},
		{
			ID:            2,
			Name:          "スタンダードプラン",
			StripePriceID: "price_1NvfOfAj9ehS6HaZXMIwK6do",
		},
	}

	return c.JSON(http.StatusOK, products)
}
