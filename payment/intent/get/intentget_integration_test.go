//go:build stripe
// +build stripe

package apppaymentintentget_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"
	apppaymentintentget "github.com/lelledev/upaygo/payment/intent/get"

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
	e = appconfig.ImportConfig(fc)
	_ = fc.Close()
	if e != nil {
		fmt.Printf("Error during file config import: %v", e)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	am := int64(2088)
	pip := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(am),
		Currency: stripe.String(cur.GetISO4217()),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
	if e != nil {
		t.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	appintent, e := apppaymentintentget.Get(intent.ID, cur)
	if e != nil {
		t.Errorf("impossible to get %v payment intent: %v", intent.ID, e)
	}

	if appintent.GetGatewayReference() != intent.ID {
		t.Errorf("intent get is incorrect, got an intent with different ID. Got: %v want: %v", appintent.GetGatewayReference(), intent.ID)
	}

	_, _ = paymentintent.Cancel(appintent.GetGatewayReference(), nil)
}

func TestGetWithoutID(t *testing.T) {
	cur, _ := appcurrency.New("EUR")
	_, e := apppaymentintentget.Get("", cur)
	if e == nil {
		t.Error("expecting an error if get an intent without ID")
	}
}

func TestGetWithoutCurrency(t *testing.T) {
	_, e := apppaymentintentget.Get("in_xxx", nil)
	if e == nil {
		t.Error("expecting an error if get an intent without currency")
	}
}
