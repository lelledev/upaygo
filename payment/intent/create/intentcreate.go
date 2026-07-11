package apppaymentintentcreate

import (
	"context"
	"errors"

	appamount "github.com/lelledev/upaygo/amount"
	appconfig "github.com/lelledev/upaygo/config"
	appcustomer "github.com/lelledev/upaygo/customer"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintent "github.com/lelledev/upaygo/payment/intent"
	apppaymentsource "github.com/lelledev/upaygo/payment/source"

	"github.com/stripe/stripe-go/v82"
)

// Create creates an intent in Stripe and returns it as an instance of Intent
func Create(a appamount.Amount, p apppaymentsource.Source, c appcustomer.Customer) (apppaymentintent.Intent, error) {
	if a == nil || p == nil {
		return nil, errors.New("impossible to create a payment intent without required parameters")
	}

	sc, e := appconfig.ClientForCurrency(a.GetCurrency().GetISO4217())
	if e != nil {
		return nil, e
	}

	ic := &stripe.PaymentIntentCreateParams{
		Amount:             new(int64(a.GetAmount())),
		Currency:           new(a.GetCurrency().GetISO4217()),
		PaymentMethod:      new(p.GetGatewayReference()),
		SetupFutureUsage:   new("off_session"),
		ConfirmationMethod: new("manual"),
		CaptureMethod:      new("manual"),
		// Explicit types avoid Dashboard redirect methods (no return_url on this
		// API). payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	if c != nil {
		// With SetupFutureUsage set, Stripe attaches the payment method to the
		// customer after confirmation (SavePaymentMethod was removed in stripe-go v72+).
		ic.Customer = new(c.GetGatewayReference())
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), ic)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return apppaymentintent.FromStripeToAppIntent(*intent), nil
}
