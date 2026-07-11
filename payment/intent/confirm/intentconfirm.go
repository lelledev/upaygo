package apppaymentintentconfirm

import (
	"context"
	"fmt"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Confirm gets the intent id from c Stripe account and confirm it
func Confirm(id string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if id == "" || c == nil {
		return nil, apperror.Invalid("impossible to confirm the payment intent without required parameters")
	}

	sc, e := appconfig.GetStripeClientByCurrency(c.GetISO4217())
	if e != nil {
		return nil, fmt.Errorf("impossible to get the Stripe client to confirm the payment intent: %w", e)
	}

	intent, e := sc.V1PaymentIntents.Confirm(context.Background(), id, nil)
	if e != nil {
		return nil, fmt.Errorf("confirm payment intent %s: %w", id, e)
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
