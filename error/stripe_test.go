//go:build unit
// +build unit

package apperror_test

import (
	"errors"
	"testing"

	apperror "github.com/lelledev/upaygo/error"
)

func TestGetStripeErrorMessage(t *testing.T) {
	es := errors.New("{\"status\":401,\"message\":\"Invalid API Key provided: sk_xxx\",\"type\":\"invalid_request_error\"}")
	got, e := apperror.GetStripeErrorMessage(es)
	if e != nil {
		t.Errorf("impossible to get the error message from Stripe built-in error: %v", e)
	}

	if w := "Invalid API Key provided: sk_xxx"; got != w {
		t.Errorf("Get the error message from Stripe built-in error is incorrect, got: %v want: %v", got, w)
	}
}
