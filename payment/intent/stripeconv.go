package apppaymentintent

import (
	"time"

	appamount "github.com/lelledaniele/upaygo/amount"
	appcustomer "github.com/lelledaniele/upaygo/customer"
	apppaymentsource "github.com/lelledaniele/upaygo/payment/source"

	"github.com/stripe/stripe-go/v76"
)

func FromStripeToAppIntent(intent stripe.PaymentIntent) Intent {
	var nat string
	na := intent.NextAction
	if na != nil {
		nat = string(na.Type)
	}

	var ps apppaymentsource.Source
	if intent.PaymentMethod != nil {
		ps = apppaymentsource.New(intent.PaymentMethod.ID)
	}

	a, _ := appamount.New(int(intent.Amount), intent.Currency)

	var cus appcustomer.Customer
	if intent.Customer != nil {
		cus = appcustomer.New(intent.Customer.ID, intent.Customer.Email)
	}

	return &i{
		R:    intent.ID,
		CM:   string(intent.ConfirmationMethod),
		NA:   nat,
		OFFS: intent.SetupFutureUsage == "off_session",
		CT:   time.Unix(intent.Created, 0),
		C:    cus,
		PS:   ps,
		A:    a,
		S:    status{S: string(intent.Status)},
	}
}
