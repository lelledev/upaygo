//go:build stripe
// +build stripe

package apppaymentintentcapture_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apppaymentintentcapture "github.com/lelledev/upaygo/payment/intent/capture"

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

func TestCapture(t *testing.T) {
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
		t.Fatalf("impossible to create Stripe client: %v", e)
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), pip)
	if e != nil {
		t.Fatalf("impossible to create a new payment intent for testing: %v", e)
	}

	// Succeeded PaymentIntents cannot be cancelled; no teardown cancel after capture.

	appintent, e := apppaymentintentcapture.Capture(intent.ID, cur)
	if e != nil {
		// Capture failed — intent may still be cancelable (e.g. requires_capture).
		t.Cleanup(func() {
			if _, err := sc.V1PaymentIntents.Cancel(context.Background(), intent.ID, nil); err != nil {
				t.Errorf("cleanup cancel payment intent %s: %v", intent.ID, err)
			}
		})
		t.Fatalf("impossible to capture %v payment intent: %v", intent.ID, e)
	}

	if !appintent.IsSucceeded() {
		t.Error("intent capture is incorrect, got an intent that is not succeeded")
	}
}

func TestCaptureWithSCACard(t *testing.T) {
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
		t.Fatalf("impossible to create Stripe client: %v", e)
	}

	intent, e := sc.V1PaymentIntents.Create(context.Background(), pip)
	if e != nil {
		t.Fatalf("impossible to create a new payment intent for testing: %v", e)
	}

	t.Cleanup(func() {
		if _, err := sc.V1PaymentIntents.Cancel(context.Background(), intent.ID, nil); err != nil {
			t.Errorf("cleanup cancel payment intent %s: %v", intent.ID, err)
		}
	})

	_, e = apppaymentintentcapture.Capture(intent.ID, cur)
	if e == nil {
		t.Errorf("intent %v should not be captured as it should have status requires_action", intent.ID)
	}
}

func TestCaptureWithoutID(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	_, e := apppaymentintentcapture.Capture("", cur)
	if e == nil {
		t.Error("expecting an error if capture an intent without ID")
	}
}

func TestCaptureWithoutCurrency(t *testing.T) {
	_, e := apppaymentintentcapture.Capture("in_xxx", nil)
	if e == nil {
		t.Error("expecting an error if capture an intent without currency")
	}
}
