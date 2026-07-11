//go:build stripe
// +build stripe

package apprestintentcancel_test

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

	appconfig "github.com/lelledev/upaygo/config"
	apprestintentcancel "github.com/lelledev/upaygo/controller/rest/intent/cancel"
	appcurrency "github.com/lelledev/upaygo/currency"
	appstripetest "github.com/lelledev/upaygo/internal/stripetest"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v82"
)

const (
	errorRestCreateIntent = "cancel intent controller failed: %v"
)

type responseIntent struct {
	IntentGatewayReference string         `json:"gateway_reference"`
	Status                 responseStatus `json:"status"`
}

type responseStatus struct {
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

func createTestIntent() (string, error) {
	cur, _ := appcurrency.New("EUR")
	am := int64(4567)
	pip := &stripe.PaymentIntentCreateParams{
		Amount:             new(am),
		Currency:           new(cur.GetISO4217()),
		PaymentMethod:      new("pm_card_visa"),
		SetupFutureUsage:   new("off_session"),
		ConfirmationMethod: new("automatic"),
		Confirm:            new(true),
		CaptureMethod:      new("manual"),
		// payment_method_types is compatible with confirmation_method;
		// automatic_payment_methods is not (Stripe rejects both together).
		PaymentMethodTypes: []*string{new("card")},
	}

	intent, e := appstripetest.NewIntent(cur.GetISO4217(), pip)
	if e != nil {
		return "", fmt.Errorf("impossible to create a new payment intent for testing: %w", e)
	}

	return intent.ID, e
}

// Test a create intent request
func Test(t *testing.T) {
	intentID, e := createTestIntent()
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		appstripetest.CancelIntent("EUR", intentID)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("currency=EUR"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": intentID})

	apprestintentcancel.Handler(w, req)

	res := w.Result()
	resBody, e := io.ReadAll(res.Body)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}
	defer res.Body.Close()

	var resI responseIntent
	e = json.Unmarshal(resBody, &resI)
	if e != nil {
		t.Errorf(errorRestCreateIntent, e)
	}

	if resI.IntentGatewayReference == "" {
		t.Errorf(errorRestCreateIntent, "the body response does not have the gateway reference")
	}

	if resI.Status.R != "canceled" {
		t.Errorf(errorRestCreateIntent, "the body response does not have the status 'canceled'")
	}
}
