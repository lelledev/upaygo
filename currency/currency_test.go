//go:build unit
// +build unit

package appcurrency_test

import (
	"strings"
	"testing"

	appcurrency "github.com/lelledev/upaygo/currency"
)

func TestNew(t *testing.T) {
	cs := "EUR"
	c, e := appcurrency.New(strings.ToLower(cs))
	if e != nil {
		t.Errorf("error during the appcurrency.New func: %v", e)
	}

	if c.GetISO4217() != cs {
		t.Errorf("appcurrency.New func returns a currency ISO4217 incorrect: got: %v want %v", c.GetISO4217(), cs)
	}
}

func TestNewWithWrongCurrency(t *testing.T) {
	_, e := appcurrency.New("I AM A WRONG CURRENCY")

	if e == nil {
		t.Error("No error during appcurrency creation with a wrong ISO4217")
	}
}
