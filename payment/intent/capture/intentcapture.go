package apppaymentintentcapture

import (
	"context"
	"fmt"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Capture gets the intent id from c Stripe account and capture it
func Capture(id string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if id == "" || c == nil {
		return nil, apperror.Invalid("impossible to capture the payment intent without required parameters")
	}

	sc, e := appconfig.GetStripeClientByCurrency(c.GetISO4217())
	if e != nil {
		return nil, fmt.Errorf("impossible to get the Stripe client to capture the payment intent: %w", e)
	}

	intent, e := sc.V1PaymentIntents.Capture(context.Background(), id, nil)
	if e != nil {
		return nil, fmt.Errorf("capture payment intent %s: %w", id, e)
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
