package appconfig

import (
	"fmt"

	"github.com/stripe/stripe-go/v82"
)

// ClientForCurrency returns a stripe.Client authenticated with the secret key
// configured for iso4217 (or the default key if the currency is not listed).
// Prefer this over a package-level global secret so concurrent requests with
// different currencies do not race on shared API credentials.
func ClientForCurrency(iso4217 string) (*stripe.Client, error) {
	cfg, err := GetStripeAPIConfigByCurrency(iso4217)
	if err != nil {
		return nil, fmt.Errorf("load Stripe configuration for %q: %w", iso4217, err)
	}

	return stripe.NewClient(cfg.GetSK()), nil
}

// ClientForCurrencyOp is like ClientForCurrency but wraps lookup failures with
// an operation name for clearer call-site diagnostics.
func ClientForCurrencyOp(iso4217, op string) (*stripe.Client, error) {
	sc, err := ClientForCurrency(iso4217)
	if err != nil {
		return nil, fmt.Errorf("create Stripe client for %s currency %q: %w", op, iso4217, err)
	}
	return sc, nil
}
