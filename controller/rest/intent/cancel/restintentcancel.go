package apprestintentcancel

import (
	"encoding/json"
	"fmt"
	"net/http"

	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintentcancel "github.com/lelledev/upaygo/payment/intent/cancel"

	"github.com/gorilla/mux"
)

const (
	URL    = "/payment_intents/{id}/cancel"
	Method = http.MethodPost

	responseTye = "application/json"

	errorParsingParam        = "error during the payload parsing: '%v'"
	errorParamPayloadMissing = "missing payload mandatory parameters to cancel a payment intent"
	errorAmountCreation      = "error during the intent amount creation: '%v'"
	errorIntentEncoding      = "error during the intent encoding: '%v'"
)

// @Summary Cancel an intent
// @Description Cancel an confirmed intent
// @Tags Intent
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Intent's ID"
// @Param currency formData string true "Intent's currency"
// @Success 200 {interface} apppaymentintent.Intent
// @Failure 400 {object} apperror.RESTError
// @Failure 405 {object} apperror.RESTError
// @Failure 500 {object} apperror.RESTError
// @Router /payment_intents/{id}/cancel [post]
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", responseTye)

	ID, cur, e := getParams(r)
	if e != nil {
		apperror.WriteJSON(w, e)
		return
	}

	appintent, e := apppaymentintentcancel.Cancel(ID, cur)
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
		return "", nil, apperror.Invalid(fmt.Sprintf(errorParsingParam, e))
	}

	p := r.Form
	if p.Get("currency") == "" {
		return "", nil, apperror.Invalid(errorParamPayloadMissing)
	}

	cur, e := appcurrency.New(p.Get("currency"))
	if e != nil {
		return "", nil, apperror.Invalid(fmt.Sprintf(errorAmountCreation, e))
	}

	return ID, cur, nil
}
