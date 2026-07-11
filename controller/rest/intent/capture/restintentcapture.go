package apprestintentcapture

import (
	"encoding/json"
	"fmt"
	"net/http"

	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintentcapture "github.com/lelledev/upaygo/payment/intent/capture"

	"github.com/gorilla/mux"
)

const (
	URL    = "/payment_intents/{id}/capture"
	Method = http.MethodPost

	responseTye = "application/json"
)

// @Summary Capture an intent
// @Description Capture an confirmed intent
// @Tags Intent
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Intent's ID"
// @Param currency formData string true "Intent's currency"
// @Success 200 {interface} apppaymentintent.Intent
// @Failure 400 {object} apperror.RESTError
// @Failure 405 {object} apperror.RESTError
// @Failure 500 {object} apperror.RESTError
// @Router /payment_intents/{id}/capture [post]
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", responseTye)

	ID, cur, e := getParams(r)
	if e != nil {
		apperror.WriteJSON(w, e)
		return
	}

	appintent, e := apppaymentintentcapture.Capture(ID, cur)
	if e != nil {
		apperror.WriteJSON(w, e)
		return
	}

	if e = json.NewEncoder(w).Encode(appintent); e != nil {
		apperror.WriteJSON(w, fmt.Errorf("encode payment intent response: %w", e))
	}
}

// Get and transform the payload params into domain structs
func getParams(r *http.Request) (string, appcurrency.Currency, error) {
	vars := mux.Vars(r)
	ID := vars["id"]

	if e := r.ParseForm(); e != nil {
		return "", nil, apperror.Invalid(fmt.Sprintf("error during the payload parsing: %v", e))
	}

	p := r.Form
	if p.Get("currency") == "" {
		return "", nil, apperror.Invalid("missing payload mandatory parameters to capture a payment intent")
	}

	cur, e := appcurrency.New(p.Get("currency"))
	if e != nil {
		return "", nil, apperror.Invalid(fmt.Sprintf("error during the intent amount creation: %v", e))
	}

	return ID, cur, nil
}
