package apppaymentintentget

import (
	"context"
	"fmt"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Get gets the gf intent from c Stripe account and returns it as an instance of i
func Get(gf string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if gf == "" || c == nil {
		return nil, apperror.Invalid("impossible to get the payment intent without required parameters")
	}

	sc, e := appconfig.GetStripeClientByCurrency(c.GetISO4217())
	if e != nil {
		return nil, fmt.Errorf("impossible to get the Stripe client to get the payment intent: %w", e)
	}

	intent, e := sc.V1PaymentIntents.Retrieve(context.Background(), gf, nil)
	if e != nil {
		return nil, fmt.Errorf("get payment intent %s: %w", gf, e)
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
