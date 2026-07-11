//go:build stripe
// +build stripe

package apppaymentintentconfirm_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	appstripetest "github.com/lelledev/upaygo/internal/stripetest"
	apppaymentintentconfirm "github.com/lelledev/upaygo/payment/intent/confirm"

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

func TestConfirm(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentCreateParams{
		Amount:             new(am),
		Currency:           new(cur.GetISO4217()),
		ConfirmationMethod: new("manual"),
		PaymentMethod:      new("pm_card_visa"),
		// payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	intent, e := appstripetest.NewIntent(cur.GetISO4217(), pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentconfirm.Confirm(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to confirm %v payment intent: %v", intent.ID, e)
	}

	if appintent.RequiresConfirmation() {
		t.Error("intent confirmation is incorrect, got an intent that requires confirmation")
	}

	appstripetest.CancelIntent(cur.GetISO4217(), appintent.GetGatewayReference())
}

func TestConfirmWithoutID(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	_, e := apppaymentintentconfirm.Confirm("", cur)
	if e == nil {
		t.Error("expecting an error if confirm an intent without ID")
	}
}

func TestConfirmWithoutCurrency(t *testing.T) {
	_, e := apppaymentintentconfirm.Confirm("in_xxx", nil)
	if e == nil {
		t.Error("expecting an error if confirm an intent without currency")
	}
}
