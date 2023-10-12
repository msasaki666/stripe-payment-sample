package main

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/msasaki666/backend/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
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
	e.Logger.Fatal(e.Start(":1323"))
}
