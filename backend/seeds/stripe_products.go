// Code generated by SQLBoiler boilingseed-0.1.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package seeds

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	models "github.com/msasaki666/backend/models"
	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var (
	stripeProductColumnsWithDefault = []string{"id"}
	stripeProductDBTypes            = map[string]string{`ID`: `bigint`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`, `Name`: `text`, `IDOnStripe`: `text`}
)

// defaultRandomStripeProduct creates a random model.StripeProduct
// Used when RandomStripeProduct is not set in the Seeder
func defaultRandomStripeProduct() (*models.StripeProduct, error) {
	o := &models.StripeProduct{}
	seed := randomize.NewSeed()
	err := randomize.Struct(seed, o, stripeProductDBTypes, true, stripeProductColumnsWithDefault...)

	return o, err
}

func (s Seeder) seedStripeProducts(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("Adding StripeProducts")
	StripeProductsToAdd := s.MinStripeProductsToSeed

	randomFunc := s.RandomStripeProduct
	if randomFunc == nil {
		randomFunc = defaultRandomStripeProduct
	}

	for i := 0; i < StripeProductsToAdd; i++ {
		// create model
		o, err := randomFunc()
		if err != nil {
			return fmt.Errorf("unable to get Random StripeProduct: %w", err)
		}

		// insert model
		if err := o.Insert(ctx, exec, boil.Infer()); err != nil {
			return fmt.Errorf("unable to insert StripeProduct: %w", err)
		}
	}

	// run afterAdd
	if s.AfterStripeProductsAdded != nil {
		if err := s.AfterStripeProductsAdded(ctx); err != nil {
			return fmt.Errorf("error running AfterStripeProductsAdded: %w", err)
		}
	}

	fmt.Println("Finished adding StripeProducts")
	return nil
}

// These packages are needed in SOME models
// This is to prevent errors in those that do not need it
var _ = math.E
var _ = queries.Query{}

// This is to force strconv to be used. Without it, it causes an error because strconv is imported by ALL the drivers
var _ = strconv.IntSize

// stripeProduct is here to prevent erros due to driver "BasedOnType" imports.
type stripeProduct struct {
	ID         int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	IDOnStripe string
}
