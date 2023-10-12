package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"

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
	e.GET("/users", func(c echo.Context) error {
		num, err := models.Users().Count(context.Background(), db)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		return c.String(http.StatusOK, fmt.Sprint(num))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
