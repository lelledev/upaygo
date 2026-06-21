package apprestintentget

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appcurrency "github.com/lelledev/upaygo/currency"
	apperror "github.com/lelledev/upaygo/error"
	apppaymentintentget "github.com/lelledev/upaygo/payment/intent/get"

	"github.com/gorilla/mux"
)

const (
	URL    = "/payment_intents/{id}"
	Method = http.MethodGet

	responseTye = "application/json"

	errorParamPathMissing  = "missing URL in-path mandatory parameters to get the payment intent"
	errorParamQueryMissing = "error during the query parsing: '%v'"
	errorAmountCreation    = "error during the intent amount creation: '%v'"
	errorIntentGet         = "error during the intent getter: '%v'"
	errorIntentEncoding    = "error during the intent encoding: '%v'"
)

// @Summary Get an intent
// @Description Get an existing intent
// @Tags Intent
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Intent's ID"
// @Param currency formData string true "Intent's currency"
// @Success 200 {interface} apppaymentintent.Intent
// @Failure 400 {object} apperror.RESTError
// @Failure 405 {object} apperror.RESTError
// @Failure 500 {object} apperror.RESTError
// @Router /payment_intents/{id} [get]
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", responseTye)

	ID, cur, e := getParams(r)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)

		e := apperror.RESTError{
			M: e.Error(),
		}
		_ = json.NewEncoder(w).Encode(e)

		return
	}

	appintent, e := apppaymentintentget.Get(ID, cur)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)

		e := apperror.RESTError{
			M: fmt.Sprintf(errorIntentGet, e),
		}
		_ = json.NewEncoder(w).Encode(e)

		return
	}

	e = json.NewEncoder(w).Encode(appintent)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)

		e := apperror.RESTError{
			M: fmt.Sprintf(errorIntentEncoding, e),
		}
		_ = json.NewEncoder(w).Encode(e)

		return
	}
}

// Get and transform the payload params into domain structs
func getParams(r *http.Request) (string, appcurrency.Currency, error) {
	vars := mux.Vars(r)
	ID, _ := vars["id"]

	cursym := r.URL.Query().Get("currency")
	if cursym == "" {
		return "", nil, errors.New(errorParamQueryMissing)
	}

	cur, e := appcurrency.New(cursym)
	if e != nil {
		return "", nil, fmt.Errorf(errorAmountCreation, e.Error())
	}

	return ID, cur, nil
}
