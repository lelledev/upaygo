//go:build unit
// +build unit

package appcustomer_test

import (
	"testing"

	appcustomer "github.com/lelledev/upaygo/customer"
)

func TestNew(t *testing.T) {
	ref, email := "GATEWAY REFERENCE", "email@email.com"
	got := appcustomer.New(ref, email)

	if got.GetGatewayReference() != ref {
		t.Errorf("The new customer.gateway_reference is incorrect, got: %v want %v", got.GetGatewayReference(), ref)
	}

	if got.GetEmail() != email {
		t.Errorf("The new customer.email is incorrect, got: %v want %v", got.GetEmail(), email)
	}
}
