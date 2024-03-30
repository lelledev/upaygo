package appcustomer

import (
	"errors"

	apperror "github.com/lelledaniele/upaygo/error"

	"github.com/stripe/stripe-go/v76/customer"

	appconfig "github.com/lelledaniele/upaygo/config"
	appcurrency "github.com/lelledaniele/upaygo/currency"

	"github.com/stripe/stripe-go/v76"
)

func NewStripe(email string, ac appcurrency.Currency) (Customer, error) {
	if email == "" || ac == nil {
		return nil, errors.New("impossible to create a Stripe customer without required parameters")
	}

	sck, e := appconfig.GetStripeAPIConfigByCurrency(ac.GetISO4217())
	if e != nil {
		return nil, e
	}

	stripe.Key = sck.GetSK()

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	cus, e := customer.New(params)
	if e != nil {
		m, es := apperror.GetStripeErrorMessage(e)
		if es == nil {
			return nil, errors.New(m)
		}

		return nil, e
	}

	return &c{
		R:     cus.ID,
		Email: cus.Email,
	}, nil
}
