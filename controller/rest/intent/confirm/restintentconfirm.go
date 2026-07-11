package apprestintentconfirm

import (
	"encoding/json"
	"fmt"
	"net/http"

	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintentconfirm "github.com/lelledev/upaygo/payment/intent/confirm"

	"github.com/gorilla/mux"
)

const (
	URL    = "/payment_intents/{id}/confirm"
	Method = http.MethodPost

	responseTye = "application/json"

	errorParsingParam        = "error during the payload parsing"
	errorParamPayloadMissing = "missing payload mandatory parameters to confirm a payment intent"
	errorCurrencyParsing     = "error during the currency parsing"
	errorIntentEncoding      = "error during the intent encoding: %w"
)

// @Summary Confirm an intent
// @Description Confirm an unconfirmed intent
// @Tags Intent
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Intent's ID"
// @Param currency formData string true "Intent's currency"
// @Success 200 {interface} apppaymentintent.Intent
// @Failure 400 {object} apperror.RESTError
// @Failure 405 {object} apperror.RESTError
// @Failure 500 {object} apperror.RESTError
// @Router /payment_intents/{id}/confirm [post]
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", responseTye)

	ID, cur, e := getParams(r)
	if e != nil {
		apperror.WriteJSON(w, e)
		return
	}

	appintent, e := apppaymentintentconfirm.Confirm(ID, cur)
	if e != nil {
		apperror.WriteJSON(w, e)
		return
	}

	if e = json.NewEncoder(w).Encode(appintent); e != nil {
		apperror.WriteJSON(w, fmt.Errorf(errorIntentEncoding, e))
	}
}

// Get and transform the payload params into domain structs
func getParams(r *http.Request) (string, appcurrency.Currency, error) {
	vars := mux.Vars(r)
	ID := vars["id"]

	if e := r.ParseForm(); e != nil {
		return "", nil, apperror.InvalidCause(errorParsingParam, e)
	}

	p := r.Form
	if p.Get("currency") == "" {
		return "", nil, apperror.Invalid(errorParamPayloadMissing)
	}

	cur, e := appcurrency.New(p.Get("currency"))
	if e != nil {
		return "", nil, apperror.InvalidCause(errorCurrencyParsing, e)
	}

	return ID, cur, nil
}
