//go:build unit
// +build unit

package apppaymentintent_test

import (
	"testing"

	apppaymentintent "github.com/lelledev/upaygo/payment/intent"

	"github.com/stripe/stripe-go"
)

func TestFromStripeToAppIntent(t *testing.T) {
	pmID := "card_xxx"
	pm := stripe.PaymentMethod{
		ID: pmID,
	}

	nat := "stripe_sdk"
	na := stripe.PaymentIntentNextAction{
		Type: stripe.PaymentIntentNextActionType(nat),
	}

	cusID := "cus_xxx"
	cus := stripe.Customer{
		ID: cusID,
	}

	intent := stripe.PaymentIntent{
		ID:            "intent_xxx",
		NextAction:    &na,
		PaymentMethod: &pm,
		Customer:      &cus,
	}

	appintent := apppaymentintent.FromStripeToAppIntent(intent)

	if appintent.GetNextAction() != nat {
		t.Errorf("incorrect next_action intent domain conversion, got: %v want: %v", appintent.GetNextAction(), nat)
	}

	if appintent.GetSource().GetGatewayReference() != pmID {
		t.Errorf("incorrect payment_method intent domain conversion, got: %v want: %v", appintent.GetSource().GetGatewayReference(), pmID)
	}

	if appintent.GetCustomer().GetGatewayReference() != cusID {
		t.Errorf("incorrect customer intent domain conversion, got: %v want: %v", appintent.GetCustomer().GetGatewayReference(), cusID)
	}
}

func TestFromStripeToAppIntentWithoutNextAction(t *testing.T) {
	pmID := "card_xxx"
	pm := stripe.PaymentMethod{
		ID: pmID,
	}

	cusID := "cus_xxx"
	cus := stripe.Customer{
		ID: cusID,
	}

	intent := stripe.PaymentIntent{
		ID:            "intent_xxx",
		PaymentMethod: &pm,
		Customer:      &cus,
	}

	appintent := apppaymentintent.FromStripeToAppIntent(intent)

	if appintent.GetNextAction() != "" {
		t.Errorf("incorrect next_action intent domain conversion, got: %v want: ''", appintent.GetNextAction())
	}
}

func TestFromStripeToAppIntentWithoutPaymentMethod(t *testing.T) {
	nat := "stripe_sdk"
	na := stripe.PaymentIntentNextAction{
		Type: stripe.PaymentIntentNextActionType(nat),
	}

	cusID := "cus_xxx"
	cus := stripe.Customer{
		ID: cusID,
	}

	intent := stripe.PaymentIntent{
		ID:         "intent_xxx",
		NextAction: &na,
		Customer:   &cus,
	}

	appintent := apppaymentintent.FromStripeToAppIntent(intent)

	if appintent.GetSource() != nil {
		t.Errorf("incorrect payment_method intent domain conversion, got: %v want: nil", appintent.GetSource().GetGatewayReference())
	}
}

func TestFromStripeToAppIntentWithoutCustomer(t *testing.T) {
	pmID := "card_xxx"
	pm := stripe.PaymentMethod{
		ID: pmID,
	}

	nat := "stripe_sdk"
	na := stripe.PaymentIntentNextAction{
		Type: stripe.PaymentIntentNextActionType(nat),
	}

	intent := stripe.PaymentIntent{
		ID:            "intent_xxx",
		NextAction:    &na,
		PaymentMethod: &pm,
	}

	appintent := apppaymentintent.FromStripeToAppIntent(intent)

	if appintent.GetCustomer() != nil {
		t.Errorf("incorrect customer intent domain conversion, got: %v want: nil", appintent.GetCustomer().GetGatewayReference())
	}
}
