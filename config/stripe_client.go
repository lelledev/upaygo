package appconfig

import "github.com/stripe/stripe-go/v82"

// ClientForCurrency returns a stripe.Client authenticated with the secret key
// configured for iso4217 (or the default key if the currency is not listed).
// Prefer this over a package-level global secret so concurrent requests with
// different currencies do not race on shared API credentials.
func ClientForCurrency(iso4217 string) (*stripe.Client, error) {
	cfg, err := GetStripeAPIConfigByCurrency(iso4217)
	if err != nil {
		return nil, err
	}

	return stripe.NewClient(cfg.GetSK()), nil
}
