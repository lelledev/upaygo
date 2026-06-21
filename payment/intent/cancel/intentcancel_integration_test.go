//go:build stripe
// +build stripe

package apppaymentintentcancel_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apppaymentintentcancel "github.com/lelledev/upaygo/payment/intent/cancel"

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

func TestCancel(t *testing.T) {
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
	pip := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(am),
		Currency:           stripe.String(cur.GetISO4217()),
		ConfirmationMethod: stripe.String("manual"),
		Confirm:            stripe.Bool(false),
		CaptureMethod:      stripe.String("manual"),
		PaymentMethod:      stripe.String("pm_card_authenticationRequiredOnSetup"),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
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
