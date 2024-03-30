//go:build stripe
// +build stripe

package apppaymentintentconfirm_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledaniele/upaygo/config"
	appcurrency "github.com/lelledaniele/upaygo/currency"
	apppaymentintentconfirm "github.com/lelledaniele/upaygo/payment/intent/confirm"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
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

func TestConfirm(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(am),
		Currency:           stripe.String(cur.GetISO4217()),
		ConfirmationMethod: stripe.String("manual"),
		PaymentMethod:      stripe.String("pm_card_visa"),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
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

	_, _ = paymentintent.Cancel(appintent.GetGatewayReference(), nil)
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
