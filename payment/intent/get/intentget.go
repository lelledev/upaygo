package apppaymentintentget

import (
	"context"
	"errors"
	"fmt"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
)

// Get gets the gf intent from c Stripe account and returns it as an instance of i
func Get(gf string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if gf == "" || c == nil {
		return nil, errors.New("impossible to get the payment intent without required parameters")
	}

	currency := c.GetISO4217()
	sc, e := appconfig.ClientForCurrency(currency)
	if e != nil {
		return nil, fmt.Errorf("create Stripe client for payment intent get currency %q: %w", currency, e)
	}

	intent, e := sc.V1PaymentIntents.Retrieve(context.Background(), gf, nil)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
