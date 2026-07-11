package apppaymentintentcancel

import (
	"context"
	"errors"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Cancel gets the intent id from c Stripe account and cancel it
func Cancel(id string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if id == "" || c == nil {
		return nil, errors.New("impossible to cancel the payment intent without required parameters")
	}

	sc, e := appconfig.ClientForCurrencyOp(c.GetISO4217(), "payment intent cancel")
	if e != nil {
		return nil, e
	}

	intent, e := sc.V1PaymentIntents.Cancel(context.Background(), id, nil)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
