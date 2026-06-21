//go:build stripe
// +build stripe

package appcustomer_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stripe/stripe-go/customer"

	appcurrency "github.com/lelledev/upaygo/currency"

	appconfig "github.com/lelledev/upaygo/config"
	appcustomer "github.com/lelledev/upaygo/customer"
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

func TestNewStripe(t *testing.T) {
	email := "email@email.com"
	c, _ := appcurrency.New("EUR")

	got, e := appcustomer.NewStripe(email, c)
	if e != nil {
		t.Errorf("error during the appcustomer creation with Stripe: %v", e)
	}

	if got.GetGatewayReference() == "" {
		t.Errorf("The new customer.gateway_reference is empty, got: %v", got.GetGatewayReference())
	}

	if got.GetEmail() != email {
		t.Errorf("The new customer.email is incorrect, got: %v want %v", got.GetEmail(), email)
	}

	_, _ = customer.Del(got.GetGatewayReference(), nil)
}
