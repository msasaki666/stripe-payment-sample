package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/msasaki666/backend/models"
)

func main() {
	ctx := context.Background()
	db, err := getDB()
	if err != nil {
		panic(err.Error())
	}

	customer := models.Customer{
		IDOnStripe: "cus_OoPK667Vsk7yyh",
	}
	if err := customer.Insert(ctx, db, boil.Infer()); err != nil {
		panic(err.Error())
	}
	one_time_product := models.StripeProduct{
		Name:       "ドロップイン1h",
		IDOnStripe: "prod_Oj7DpcILYrwsep",
	}

	if err := one_time_product.Insert(ctx, db, boil.Infer()); err != nil {
		panic(err.Error())
	}
	recurring_product := models.StripeProduct{
		Name:       "スタンダードプラン",
		IDOnStripe: "prod_Oj7dEuhMDKwZgE",
	}

	if err := recurring_product.Insert(ctx, db, boil.Infer()); err != nil {
		panic(err.Error())
	}
	one_time_price := models.StripePrice{
		StripeProductID: one_time_product.ID,
		IDOnStripe:      "price_1NveyfAj9ehS6HaZQ1SPdae2",
		Type:            "one_time",
	}
	if err := one_time_price.Insert(ctx, db, boil.Infer()); err != nil {
		panic(err.Error())
	}
	recurring_price := models.StripePrice{
		StripeProductID: recurring_product.ID,
		IDOnStripe:      "price_1NvfOfAj9ehS6HaZXMIwK6do",
		Type:            "recurring",
	}
	if err := recurring_price.Insert(ctx, db, boil.Infer()); err != nil {
		panic(err.Error())
	}

	// seeder := seeds.Seeder{}
	// seeder.MinCustomersToSeed = 1
	// seeder.MinStripeProductsToSeed = 1
	// seeder.MinStripePricesToSeed = 1
	// seeder.MinStripeRecurringsToSeed = 0

	// if err := seeder.Run(ctx, db); err != nil {
	// 	panic(err.Error())
	// }
}

func getDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=app_database sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}
