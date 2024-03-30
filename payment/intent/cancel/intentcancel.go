package apppaymentintentcancel

import (
	"errors"

	appconfig "github.com/lelledaniele/upaygo/config"
	appcurrency "github.com/lelledaniele/upaygo/currency"
	apperror "github.com/lelledaniele/upaygo/error"
	apppaymentintent "github.com/lelledaniele/upaygo/payment/intent"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// Cancel gets the intent id from c Stripe account and cancel it
func Cancel(id string, c appcurrency.Currency) (apppaymentintent.Intent, error) {
	if id == "" || c == nil {
		return nil, errors.New("impossible to cancel the payment intent without required parameters")
	}

	sck, e := appconfig.GetStripeAPIConfigByCurrency(c.GetISO4217())
	if e != nil {
		return nil, e
	}

	stripe.Key = sck.GetSK()

	intent, e := paymentintent.Cancel(id, nil)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
