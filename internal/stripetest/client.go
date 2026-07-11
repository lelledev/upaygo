package stripetest

import (
	"testing"

	appconfig "github.com/lelledev/upaygo/config"

	"github.com/stripe/stripe-go/v82"
)

// Client returns a *stripe.Client for iso4217, or fatals the test.
// Use this in Stripe integration tests when obtaining a client for fixtures
// or t.Cleanup so the error handling stays consistent.
func Client(t *testing.T, iso4217 string) *stripe.Client {
	t.Helper()
	sc, err := appconfig.ClientForCurrency(iso4217)
	if err != nil {
		t.Fatalf("impossible to create Stripe client for cleanup: %v", err)
	}
	return sc
}
