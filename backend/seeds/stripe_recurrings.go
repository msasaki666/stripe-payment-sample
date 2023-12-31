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
	stripeRecurringColumnsWithDefault = []string{"id"}
	stripeRecurringDBTypes            = map[string]string{`ID`: `bigint`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`, `Interval`: `text`, `IntervalCount`: `smallint`, `StripePriceID`: `bigint`}
)

// defaultRandomStripeRecurring creates a random model.StripeRecurring
// Used when RandomStripeRecurring is not set in the Seeder
func defaultRandomStripeRecurring() (*models.StripeRecurring, error) {
	o := &models.StripeRecurring{}
	seed := randomize.NewSeed()
	err := randomize.Struct(seed, o, stripeRecurringDBTypes, true, stripeRecurringColumnsWithDefault...)

	return o, err
}

func (s Seeder) seedStripeRecurrings(ctx context.Context, exec boil.ContextExecutor) error {
	fmt.Println("Adding StripeRecurrings")
	StripeRecurringsToAdd := s.MinStripeRecurringsToSeed

	randomFunc := s.RandomStripeRecurring
	if randomFunc == nil {
		randomFunc = defaultRandomStripeRecurring
	}

	for i := 0; i < StripeRecurringsToAdd; i++ {
		// create model
		o, err := randomFunc()
		if err != nil {
			return fmt.Errorf("unable to get Random StripeRecurring: %w", err)
		}

		// insert model
		if err := o.Insert(ctx, exec, boil.Infer()); err != nil {
			return fmt.Errorf("unable to insert StripeRecurring: %w", err)
		}
	}

	// run afterAdd
	if s.AfterStripeRecurringsAdded != nil {
		if err := s.AfterStripeRecurringsAdded(ctx); err != nil {
			return fmt.Errorf("error running AfterStripeRecurringsAdded: %w", err)
		}
	}

	fmt.Println("Finished adding StripeRecurrings")
	return nil
}

// These packages are needed in SOME models
// This is to prevent errors in those that do not need it
var _ = math.E
var _ = queries.Query{}

// This is to force strconv to be used. Without it, it causes an error because strconv is imported by ALL the drivers
var _ = strconv.IntSize

// stripeRecurring is here to prevent erros due to driver "BasedOnType" imports.
type stripeRecurring struct {
	ID            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Interval      string
	IntervalCount int16
	StripePriceID int64
}
