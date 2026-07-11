package appcustomer

import (
	"context"
	"fmt"

	apperror "github.com/lelledev/upaygo/error"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"

	"github.com/stripe/stripe-go/v82"
)

func NewStripe(email string, ac appcurrency.Currency) (Customer, error) {
	if email == "" || ac == nil {
		return nil, apperror.Invalid("impossible to create a Stripe customer without required parameters")
	}

	sc, e := appconfig.GetStripeClientByCurrency(ac.GetISO4217())
	if e != nil {
		return nil, fmt.Errorf("impossible to get the Stripe client to create the customer: %w", e)
	}

	params := &stripe.CustomerCreateParams{
		Email: new(email),
	}
	cus, e := sc.V1Customers.Create(context.Background(), params)
	if e != nil {
		return nil, fmt.Errorf("create Stripe customer: %w", e)
	}

	return &c{
		R:     cus.ID,
		Email: cus.Email,
	}, nil
}
