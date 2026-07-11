//go:build stripe
// +build stripe

package apppaymentintentcancel_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apppaymentintentcancel "github.com/lelledev/upaygo/payment/intent/cancel"

	"github.com/stripe/stripe-go/v82"
)

func TestMain(m *testing.M) {
	var fcp string

	flag.StringVar(&fcp, "config", "", "Provide config file as an absolute path")
	flag.Parse()

	if fcp == "" {
		fmt.Print("Integration Stripe test needs the config file absolute path as flag -config")
		os.Exit(1)
	}

	fc, e := os.Open(fcp)
	if e != nil {
		fmt.Printf("Impossible to get configuration file: %v\n", e)
		os.Exit(1)
	}
	e = appconfig.ImportConfig(fc)
	_ = fc.Close()
	if e != nil {
		fmt.Printf("Error during file config import: %v", e)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCancel(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentCreateParams{
		Amount:             new(am),
		Currency:           new(cur.GetISO4217()),
		ConfirmationMethod: new("automatic"),
		Confirm:            new(true),
		CaptureMethod:      new("manual"),
		PaymentMethod:      new("pm_card_visa"),
		// payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	sc, e := appconfig.ClientForCurrency(cur.GetISO4217())
	if e != nil {
		t.Errorf("impossible to create Stripe client: %v", e)
		return
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentcancel.Cancel(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to cancel %v payment intent: %v", intent.ID, e)
	}

	if !appintent.IsCanceled() {
		t.Error("intent cancel is incorrect, got an intent that is not canceled")
	}
}

func TestCancelWithSCACard(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentCreateParams{
		Amount:             new(am),
		Currency:           new(cur.GetISO4217()),
		ConfirmationMethod: new("automatic"),
		Confirm:            new(true),
		CaptureMethod:      new("manual"),
		PaymentMethod:      new("pm_card_authenticationRequiredOnSetup"),
		// payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	sc, e := appconfig.ClientForCurrency(cur.GetISO4217())
	if e != nil {
		t.Errorf("impossible to create Stripe client: %v", e)
		return
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentcancel.Cancel(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to cancel %v payment intent: %v", intent.ID, e)
	}

	if !appintent.IsCanceled() {
		t.Error("intent cancel is incorrect, got an intent that is not canceled")
	}
}

func TestCancelNonConfirmedIntent(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentCreateParams{
		Amount:             new(am),
		Currency:           new(cur.GetISO4217()),
		ConfirmationMethod: new("manual"),
		Confirm:            new(false),
		CaptureMethod:      new("manual"),
		PaymentMethod:      new("pm_card_authenticationRequiredOnSetup"),
		// payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	sc, e := appconfig.ClientForCurrency(cur.GetISO4217())
	if e != nil {
		t.Errorf("impossible to create Stripe client: %v", e)
		return
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentcancel.Cancel(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to cancel %v payment intent: %v", intent.ID, e)
	}

	if !appintent.IsCanceled() {
		t.Error("intent cancel is incorrect, got an intent that is not canceled")
	}
}

func TestCancelWithoutID(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	_, e := apppaymentintentcancel.Cancel("", cur)
	if e == nil {
		t.Error("expecting an error if cancel an intent without ID")
	}
}

func TestCancelWithoutCurrency(t *testing.T) {
	_, e := apppaymentintentcancel.Cancel("in_xxx", nil)
	if e == nil {
		t.Error("expecting an error if cancel an intent without currency")
	}
}
