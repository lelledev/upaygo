//go:build stripe
// +build stripe

package apprestintentcreate_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stripe/stripe-go/customer"

	appconfig "github.com/lelledev/upaygo/config"
	apprestintentcreate "github.com/lelledev/upaygo/controller/rest/intent/create"
	appcurrency "github.com/lelledev/upaygo/currency"
	appcustomer "github.com/lelledev/upaygo/customer"
)

const (
	errorRestCreateIntent = "create intent controller failed: %v"
)

type responseIntent struct {
	IntentGatewayReference string           `json:"gateway_reference"`
	Customer               responseCustomer `json:"customer"`
}

type responseCustomer struct {
	R string `json:"gateway_reference"`
}

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

// Test a create intent request
func Test(t *testing.T) {
	var resI responseIntent

	c, _ := appcurrency.New("EUR")
	cus, e := appcustomer.NewStripe("email@email.com", c)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}

	a, ps, w := 7777, "pm_card_visa", httptest.NewRecorder()
	p := fmt.Sprintf("currency=%v&amount=%v&payment_source=%v&customer_reference=%v", c.GetISO4217(), a, ps, cus.GetGatewayReference())
	req := httptest.NewRequest("POST", "http://example.com", strings.NewReader(p))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	apprestintentcreate.Handler(w, req)

	res := w.Result()
	resBody, e := io.ReadAll(res.Body)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}
	defer res.Body.Close()

	e = json.Unmarshal(resBody, &resI)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}

	if resI.IntentGatewayReference == "" {
		t.Errorf(errorRestCreateIntent, "the body response does not have the gateway reference")
	}

	if resI.Customer.R == "" {
		t.Errorf(errorRestCreateIntent, "the body response does not have the customer reference")
	}

	_, _ = customer.Del(cus.GetGatewayReference(), nil)
}

// Test a create intent request without customer
func TestWithoutCustomer(t *testing.T) {
	var resI responseIntent

	c, a, ps, w := "EUR", 9999, "pm_card_visa", httptest.NewRecorder()
	p := fmt.Sprintf("currency=%v&amount=%v&payment_source=%v", c, a, ps)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(p))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	apprestintentcreate.Handler(w, req)

	res := w.Result()
	resBody, e := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}

	e = json.Unmarshal(resBody, &resI)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}

	if resI.IntentGatewayReference == "" {
		t.Errorf(errorRestCreateIntent, "the body response does not have the gateway reference")
	}
}
