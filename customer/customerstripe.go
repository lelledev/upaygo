package appcustomer

import (
	"context"
	"errors"

	apperror "github.com/lelledev/upaygo/error"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"

	"github.com/stripe/stripe-go/v82"
)

func NewStripe(email string, ac appcurrency.Currency) (Customer, error) {
	if email == "" || ac == nil {
		return nil, errors.New("impossible to create a Stripe customer without required parameters")
	}

	sc, e := appconfig.GetStripeClientByCurrency(ac.GetISO4217())
	if e != nil {
		return nil, e
	}

	params := &stripe.CustomerCreateParams{
		Email: new(email),
	}
	cus, e := sc.V1Customers.Create(context.Background(), params)
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
