// Package appstripetest provides shared Stripe fixtures for integration tests:
// create a payment intent for a test, cancel it and delete customers as
// best-effort cleanup. It goes through the same per-currency stripe.Client
// used by production code, so tests never touch the deprecated package-level key.
package appstripetest

import (
	"context"

	appconfig "github.com/lelledev/upaygo/config"

	"github.com/stripe/stripe-go/v82"
)

// NewIntent creates a payment intent on the c currency Stripe account
func NewIntent(c string, params *stripe.PaymentIntentCreateParams) (*stripe.PaymentIntent, error) {
	sc, e := appconfig.GetStripeClientByCurrency(c)
	if e != nil {
		return nil, e
	}

	return sc.V1PaymentIntents.Create(context.Background(), params)
}

// CancelIntent cancels the id payment intent on the c currency Stripe account,
// ignoring errors (best-effort cleanup)
func CancelIntent(c string, id string) {
	sc, e := appconfig.GetStripeClientByCurrency(c)
	if e != nil {
		return
	}

	_, _ = sc.V1PaymentIntents.Cancel(context.Background(), id, nil)
}

// DeleteCustomer deletes the id customer on the c currency Stripe account,
// ignoring errors (best-effort cleanup)
func DeleteCustomer(c string, id string) {
	sc, e := appconfig.GetStripeClientByCurrency(c)
	if e != nil {
		return
	}

	_, _ = sc.V1Customers.Delete(context.Background(), id, nil)
}
