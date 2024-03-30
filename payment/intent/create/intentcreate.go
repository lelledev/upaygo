package apppaymentintentcreate

import (
	"errors"

	appamount "github.com/lelledaniele/upaygo/amount"
	appconfig "github.com/lelledaniele/upaygo/config"
	appcustomer "github.com/lelledaniele/upaygo/customer"
	apperror "github.com/lelledaniele/upaygo/error"
	apppaymentintent "github.com/lelledaniele/upaygo/payment/intent"
	apppaymentsource "github.com/lelledaniele/upaygo/payment/source"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// Create creates an intent in Stripe and returns it as an instance of Intent
func Create(a appamount.Amount, p apppaymentsource.Source, c appcustomer.Customer) (apppaymentintent.Intent, error) {
	if a == nil || p == nil {
		return nil, errors.New("impossible to create a payment intent without required parameters")
	}

	sck, e := appconfig.GetStripeAPIConfigByCurrency(a.GetCurrency().GetISO4217())
	if e != nil {
		return nil, e
	}

	stripe.Key = sck.GetSK()

	ic := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(a.GetAmount())),
		Currency:           stripe.String(a.GetCurrency().GetISO4217()),
		PaymentMethod:      stripe.String(p.GetGatewayReference()),
		ConfirmationMethod: stripe.String("manual"),
		CaptureMethod:      stripe.String("manual"),
		OffSession:         stripe.Bool(true),
		Confirm:            stripe.Bool(true),
	}

	if c != nil {
		ic.Customer = stripe.String(c.GetGatewayReference())
	}

	intent, e := paymentintent.New(ic)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
