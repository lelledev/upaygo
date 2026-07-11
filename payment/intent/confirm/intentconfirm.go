package apppaymentintentconfirm

import (
	"context"
	"errors"
	"fmt"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Confirm gets the intent id from c Stripe account and confirm it
func Confirm(id string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if id == "" || c == nil {
		return nil, errors.New("impossible to confirm the payment intent without required parameters")
	}

	currency := c.GetISO4217()
	sc, e := appconfig.ClientForCurrency(currency)
	if e != nil {
		return nil, fmt.Errorf("create Stripe client for payment intent confirm currency %q: %w", currency, e)
	}

	intent, e := sc.V1PaymentIntents.Confirm(context.Background(), id, nil)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
