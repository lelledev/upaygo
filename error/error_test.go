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
	t.Parallel()

	se := &stripe.Error{
		Msg:            "No such customer: cus_xxx",
		HTTPStatusCode: http.StatusBadRequest,
		Code:           stripe.ErrorCodeResourceMissing,
		Type:           stripe.ErrorTypeInvalidRequest,
	}

	tests := []struct {
		name    string
		err     error
		wantOK  bool
		wantMsg string
	}{
		{"direct", se, true, se.Msg},
		{"wrapped", fmt.Errorf("create customer: %w", se), true, se.Msg},
		{"plain", errors.New("plain"), false, ""},
		{"nil", nil, false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, ok := apperror.AsStripeError(tc.err)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !tc.wantOK {
				return
			}
			if got.Msg != tc.wantMsg {
				t.Errorf("Msg = %q, want %q", got.Msg, tc.wantMsg)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	se := &stripe.Error{Msg: "card declined"}
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"stripe", se, "card declined"},
		{"wrapped stripe", fmt.Errorf("capture intent: %w", se), "card declined"},
		{"plain", errors.New("plain failure"), "plain failure"},
		{"nil", nil, ""},
		{"invalid", apperror.Invalid("missing currency"), "missing currency"},
		{"invalid cause", apperror.InvalidCause("currency parsing", errors.New("bad iso")), "currency parsing: bad iso"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := apperror.Message(tc.err); got != tc.want {
				t.Errorf("Message() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, http.StatusOK},
		{"invalid", apperror.Invalid("missing currency"), http.StatusBadRequest},
		{"wrapped invalid", fmt.Errorf("params: %w", apperror.Invalid("bad")), http.StatusBadRequest},
		{"invalid cause", apperror.InvalidCause("parse", errors.New("x")), http.StatusBadRequest},
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
			t.Parallel()
			if got := apperror.HTTPStatus(tc.err); got != tc.want {
				t.Errorf("HTTPStatus() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestIsInvalid(t *testing.T) {
	t.Parallel()

	cause := errors.New("underlying")
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"invalid", apperror.Invalid("x"), true},
		{"invalid cause", apperror.InvalidCause("x", cause), true},
		{"wrapped invalid", fmt.Errorf("wrap: %w", apperror.Invalid("x")), true},
		{"plain", errors.New("nope"), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := apperror.IsInvalid(tc.err); got != tc.want {
				t.Errorf("IsInvalid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInvalidCauseUnwrap(t *testing.T) {
	t.Parallel()

	cause := errors.New("bad iso")
	err := apperror.InvalidCause("currency parsing", cause)
	if !errors.Is(err, cause) {
		t.Fatal("InvalidCause must unwrap to the original cause")
	}
	if got := err.Error(); got != "currency parsing: bad iso" {
		t.Errorf("Error() = %q", got)
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
		wantBody   bool
	}{
		{
			name: "stripe error",
			err: fmt.Errorf("get intent: %w", &stripe.Error{
				Msg:            "No such payment_intent: pi_xxx",
				HTTPStatusCode: http.StatusNotFound,
			}),
			wantStatus: http.StatusNotFound,
			wantMsg:    "No such payment_intent: pi_xxx",
			wantBody:   true,
		},
		{
			name:       "invalid error",
			err:        apperror.Invalid("missing currency"),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "missing currency",
			wantBody:   true,
		},
		{
			name:       "invalid cause",
			err:        apperror.InvalidCause("currency parsing", errors.New("bad iso")),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "currency parsing: bad iso",
			wantBody:   true,
		},
		{
			name:       "nil error",
			err:        nil,
			wantStatus: http.StatusOK,
			wantBody:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			apperror.WriteJSON(rec, tc.err)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if !tc.wantBody {
				if rec.Body.Len() != 0 {
					t.Errorf("expected empty body, got %q", rec.Body.String())
				}
				return
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			var body apperror.RESTError
			if e := json.NewDecoder(rec.Body).Decode(&body); e != nil {
				t.Fatalf("decode body: %v", e)
			}
			if body.M != tc.wantMsg {
				t.Errorf("body.error = %q, want %q", body.M, tc.wantMsg)
			}
		})
	}
}

func TestRESTErrorError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  *apperror.RESTError
		want string
	}{
		{"value", &apperror.RESTError{M: "fail"}, "fail"},
		{"nil receiver", nil, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}
