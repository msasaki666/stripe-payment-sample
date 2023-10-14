package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/msasaki666/backend/seeds"
)

func main() {
	ctx := context.Background()
	db, err := getDB()
	if err != nil {
		panic(err.Error())
	}
	seeder := seeds.Seeder{}
	seeder.MinCustomersToSeed = 1
	seeder.MinStripeProductsToSeed = 1
	seeder.MinStripePricesToSeed = 1
	seeder.MinStripeRecurringsToSeed = 0

	if err := seeder.Run(ctx, db); err != nil {
		panic(err.Error())
	}
}
func getDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=app_database sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}
