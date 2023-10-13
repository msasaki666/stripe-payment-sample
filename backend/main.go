package main

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"

	"github.com/labstack/echo/v4"
	"github.com/msasaki666/backend/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	// stripeのテスト用のキー
	stripe.Key = "sk_test_09l3shTSTKHYCzzZZsiLl2vA"

	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=app_database sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	boil.SetDB(db)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/payment_stripes", func(c echo.Context) error {
		p, err := models.PaymentStripes().All(context.Background(), db)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		return c.JSONPretty(http.StatusOK, p, "  ")
	})
	e.POST("/create-checkout-session", createCheckoutSession)
	e.Any("/webhook", echo.WrapHandler(http.HandlerFunc(handleWebhook)))
	e.Logger.Fatal(e.Start(":1323"))
}

func createCheckoutSession(c echo.Context) (err error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("jpy"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("T-shirt"),
					},
					UnitAmount: stripe.Int64(2000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:3000/success"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
		Locale:     stripe.String("auto"),
	}

	s, _ := session.New(params)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, s.URL)
}
