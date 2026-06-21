//go:build stripe
// +build stripe

package apprestintentget_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	apprestintentget "github.com/lelledev/upaygo/controller/rest/intent/get"

	appconfig "github.com/lelledev/upaygo/config"
	appcurrency "github.com/lelledev/upaygo/currency"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
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
	pip := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(am),
		Currency:           stripe.String(cur.GetISO4217()),
		PaymentMethod:      stripe.String("pm_card_visa"),
		SetupFutureUsage:   stripe.String("off_session"),
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(true),
		CaptureMethod:      stripe.String("manual"),
	}

	sck, _ := appconfig.GetStripeAPIConfigByCurrency(cur.GetISO4217())
	stripe.Key = sck.GetSK()

	intent, e := paymentintent.New(pip)
	if e != nil {
		return "", fmt.Errorf("impossible to create a new payment intent for testing: %v", e)
	}

	return intent.ID, e
}

// Test a create intent request
func Test(t *testing.T) {
	intentID, e := createTestIntent()
	if e != nil {
		t.Error(e.Error())
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req = mux.SetURLVars(req, map[string]string{"id": intentID})

	q := req.URL.Query()
	q.Add("currency", "EUR")
	req.URL.RawQuery = q.Encode()

	apprestintentget.Handler(w, req)

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

	_, _ = paymentintent.Cancel(intentID, nil)
}
