//go:build unit
// +build unit

package appamount_test

import (
	"testing"

	appamount "github.com/lelledev/upaygo/amount"
	appcurrency "github.com/lelledev/upaygo/currency"
)

func TestNew(t *testing.T) {
	a := 4099
	c, _ := appcurrency.New("EUR")
	got, e := appamount.New(a, c.GetISO4217())
	if e != nil {
		t.Errorf("error during the amount creation: %v", e)
	}

	if got.GetCurrency().GetISO4217() != c.GetISO4217() {
		t.Errorf("error during the amount creation, currency value incorrect, got: %v want: %v", got.GetCurrency(), c)
	}

	if got.GetAmount() != a {
		t.Errorf("error during the amount creation, amount value incorrect, got: %v want: %v", got.GetAmount(), a)
	}
}

func TestNewWrongCurrency(t *testing.T) {
	_, e := appamount.New(10, "I AM A WRONG CURRENCY")
	if e == nil {
		t.Error("new amount with wrong currency created")
	}
}
