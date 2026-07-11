//go:build unit
// +build unit

package apperror_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apperror "github.com/lelledev/upaygo/error"

	"github.com/stripe/stripe-go/v82"
)

func TestAsStripeError(t *testing.T) {
	se := &stripe.Error{
		Msg:            "No such customer: cus_xxx",
		HTTPStatusCode: http.StatusBadRequest,
		Code:           stripe.ErrorCodeResourceMissing,
		Type:           stripe.ErrorTypeInvalidRequest,
	}

	got, ok := apperror.AsStripeError(se)
	if !ok {
		t.Fatal("expected AsStripeError to match a *stripe.Error")
	}
	if got.Msg != se.Msg {
		t.Errorf("Msg = %q, want %q", got.Msg, se.Msg)
	}

	wrapped := fmt.Errorf("create customer: %w", se)
	got, ok = apperror.AsStripeError(wrapped)
	if !ok {
		t.Fatal("expected AsStripeError to unwrap a wrapped *stripe.Error")
	}
	if got.Msg != se.Msg {
		t.Errorf("unwrapped Msg = %q, want %q", got.Msg, se.Msg)
	}

	if _, ok := apperror.AsStripeError(errors.New("plain")); ok {
		t.Error("plain errors must not match AsStripeError")
	}
	if _, ok := apperror.AsStripeError(nil); ok {
		t.Error("nil must not match AsStripeError")
	}
}

func TestMessage(t *testing.T) {
	se := &stripe.Error{Msg: "card declined"}
	if got := apperror.Message(se); got != "card declined" {
		t.Errorf("Message(stripe) = %q, want card declined", got)
	}

	wrapped := fmt.Errorf("capture intent: %w", se)
	if got := apperror.Message(wrapped); got != "card declined" {
		t.Errorf("Message(wrapped stripe) = %q, want card declined", got)
	}

	if got := apperror.Message(errors.New("plain failure")); got != "plain failure" {
		t.Errorf("Message(plain) = %q, want plain failure", got)
	}
	if got := apperror.Message(nil); got != "" {
		t.Errorf("Message(nil) = %q, want empty", got)
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, http.StatusOK},
		{"invalid", apperror.Invalid("missing currency"), http.StatusBadRequest},
		{"wrapped invalid", fmt.Errorf("params: %w", apperror.Invalid("bad")), http.StatusBadRequest},
		{"stripe 400", &stripe.Error{HTTPStatusCode: 400, Msg: "bad request"}, http.StatusBadRequest},
		{"stripe 402", &stripe.Error{HTTPStatusCode: 402, Msg: "card"}, http.StatusPaymentRequired},
		{"stripe 404", &stripe.Error{HTTPStatusCode: 404, Msg: "missing"}, http.StatusNotFound},
		{"stripe 401 maps to 500", &stripe.Error{HTTPStatusCode: 401, Msg: "bad key"}, http.StatusInternalServerError},
		{"stripe 500 maps to 502", &stripe.Error{HTTPStatusCode: 500, Msg: "stripe down"}, http.StatusBadGateway},
		{"wrapped stripe 404", fmt.Errorf("get intent: %w", &stripe.Error{HTTPStatusCode: 404, Msg: "gone"}), http.StatusNotFound},
		{"plain", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := apperror.HTTPStatus(tc.err); got != tc.want {
				t.Errorf("HTTPStatus() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestIsInvalid(t *testing.T) {
	if !apperror.IsInvalid(apperror.Invalid("x")) {
		t.Error("IsInvalid should match Invalid()")
	}
	if !apperror.IsInvalid(fmt.Errorf("wrap: %w", apperror.Invalid("x"))) {
		t.Error("IsInvalid should match wrapped Invalid()")
	}
	if apperror.IsInvalid(errors.New("nope")) {
		t.Error("IsInvalid must not match plain errors")
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	err := fmt.Errorf("get intent: %w", &stripe.Error{
		Msg:            "No such payment_intent: pi_xxx",
		HTTPStatusCode: http.StatusNotFound,
	})

	apperror.WriteJSON(rec, err)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body apperror.RESTError
	if e := json.NewDecoder(rec.Body).Decode(&body); e != nil {
		t.Fatalf("decode body: %v", e)
	}
	if body.M != "No such payment_intent: pi_xxx" {
		t.Errorf("body.error = %q, want Stripe message", body.M)
	}
}

func TestWriteJSONInvalid(t *testing.T) {
	rec := httptest.NewRecorder()
	apperror.WriteJSON(rec, apperror.Invalid("missing currency"))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var body apperror.RESTError
	if e := json.NewDecoder(rec.Body).Decode(&body); e != nil {
		t.Fatalf("decode body: %v", e)
	}
	if body.M != "missing currency" {
		t.Errorf("body.error = %q", body.M)
	}
}

func TestWriteJSONNil(t *testing.T) {
	rec := httptest.NewRecorder()
	apperror.WriteJSON(rec, nil)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Errorf("nil err should not write a response, code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestRESTErrorError(t *testing.T) {
	e := &apperror.RESTError{M: "fail"}
	if e.Error() != "fail" {
		t.Errorf("RESTError.Error() = %q", e.Error())
	}
	var nilE *apperror.RESTError
	if nilE.Error() != "" {
		t.Errorf("nil RESTError.Error() = %q", nilE.Error())
	}
}
