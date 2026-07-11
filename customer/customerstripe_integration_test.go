//go:build stripe
// +build stripe

package appcustomer_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	appcurrency "github.com/lelledev/upaygo/currency"

	appconfig "github.com/lelledev/upaygo/config"
	appcustomer "github.com/lelledev/upaygo/customer"
	"github.com/lelledev/upaygo/internal/stripetest"
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

	sc := stripetest.Client(t, c.GetISO4217())

	got, e := appcustomer.NewStripe(email, c)
	if e != nil {
		t.Fatalf("error during the appcustomer creation with Stripe: %v", e)
	}

	customerID := got.GetGatewayReference()
	t.Cleanup(func() {
		if _, err := sc.V1Customers.Delete(context.Background(), customerID, nil); err != nil {
			t.Errorf("cleanup delete customer %s: %v", customerID, err)
		}
	})

	if customerID == "" {
		t.Errorf("The new customer.gateway_reference is empty, got: %v", customerID)
	}

	if got.GetEmail() != email {
		t.Errorf("The new customer.email is incorrect, got: %v want %v", got.GetEmail(), email)
	}
}
