package appconfig

import (
	"fmt"

	"github.com/stripe/stripe-go/v82"
)

// GetStripeClientByCurrency returns a Stripe client bound to the secret key
// configured for c currency (or the default keys if c is not configured).
// The key is bound to the client instead of the deprecated package-level key,
// so concurrent requests with different currencies cannot race on it.
func GetStripeClientByCurrency(c string) (*stripe.Client, error) {
	sck, e := GetStripeAPIConfigByCurrency(c)
	if e != nil {
		return nil, fmt.Errorf("impossible to get the Stripe API configuration for %v currency: %w", c, e)
	}

	return stripe.NewClient(sck.GetSK()), nil
}
