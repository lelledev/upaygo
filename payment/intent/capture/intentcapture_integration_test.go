//go:build stripe
// +build stripe

package apppaymentintentcapture_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apppaymentintentcapture "github.com/lelledev/upaygo/payment/intent/capture"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
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
	defer fc.Close()

	e = appconfig.ImportConfig(fc)
	if e != nil {
		fmt.Printf("Error durring file config import: %v", e)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCapture(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(am),
		Currency:           stripe.String(cur.GetISO4217()),
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(true),
		CaptureMethod:      stripe.String("manual"),
		PaymentMethod:      stripe.String("pm_card_visa"),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentcapture.Capture(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to capture %v payment intent: %v", intent.ID, e)
	}

	if !appintent.IsSucceeded() {
		t.Error("intent capture is incorrect, got an intent that is not succeeded")
	}

	_, _ = paymentintent.Cancel(appintent.GetGatewayReference(), nil)
}

func TestCaptureWithSCACard(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(am),
		Currency:           stripe.String(cur.GetISO4217()),
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(true),
		CaptureMethod:      stripe.String("manual"),
		PaymentMethod:      stripe.String("pm_card_authenticationRequiredOnSetup"),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	_, e = apppaymentintentcapture.Capture(intent.ID, cur)
	if e == nil {
		t.Errorf("intent %v should not be captured as it should have status requires_action", intent.ID)
	}

	_, _ = paymentintent.Cancel(intent.ID, nil)
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
