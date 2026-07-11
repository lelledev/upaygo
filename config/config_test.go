//go:build unit
// +build unit

package appconfig_test

import (
	"strings"
	"testing"

	appconfig "github.com/lelledev/upaygo/config"
)

const (
	confStripeWithDefault = `
{
  "stripe": {
    "api_keys": {
	  "EUR": {
		"pk_key": "pk_EUR",
		"sk_key": "sk_EUR"
	  },
	  "default": {
        "pk_key": "pk_DEFAULT",
        "sk_key": "sk_DEFAULT"
      }
    }
  }
}
`
	confStripeWithoutDefault = `
{
  "stripe": {
    "api_keys": {
      "EUR": {
        "pk_key": "pk_EUR",
        "sk_key": "sk_EUR"
      }
    }
  }
}
`
	confServer = `
{
  "server": {
    "protocol": "https://",
    "domain": "localhost",
    "port": "8080"
  }
}
`
)

func TestDefaultStripeAPIConfig(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	got, e := appconfig.GetStripeAPIConfigByCurrency("NOT_FOUND_CURRENCY")
	if e != nil {
		t.Errorf("error during the retrieve of Stripe API config by inexistent currency: %v", e)
	}

	if got.GetPK() != "pk_DEFAULT" || got.GetSK() != "sk_DEFAULT" {
		t.Errorf("Stripe API config for an inexistent currency does not return the default value")
	}
}

func TestWithoutDefaultStripeAPIConfig(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithoutDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	_, e = appconfig.GetStripeAPIConfigByCurrency("NOT_FOUND_CURRENCY")
	if e == nil {
		t.Error("configuration without Stripe default API keys must return an error")
	}
}

func TestCurrencyStripeAPIConfig(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	got, e := appconfig.GetStripeAPIConfigByCurrency("EUR")
	if e != nil {
		t.Errorf("error during the retrieve of Stripe API config by EUR currency: %v", e)
	}

	if got.GetPK() != "pk_EUR" || got.GetSK() != "sk_EUR" {
		t.Errorf("incorrect Stripe API config for EUR currency, got %v and %v, want %v %v", got.GetPK(), got.GetSK(), "pk_EUR", "sk_EUR")
	}
}

func TestLowercaseCurrencyStripeAPIConfig(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	got, e := appconfig.GetStripeAPIConfigByCurrency("eur")
	if e != nil {
		t.Errorf("error during the retrieve of Stripe API config by eur currency: %v", e)
	}

	if got.GetPK() != "pk_EUR" || got.GetSK() != "sk_EUR" {
		t.Errorf("incorrect Stripe API config for eur currency, got %v and %v, want %v %v", got.GetPK(), got.GetSK(), "pk_EUR", "sk_EUR")
	}
}

func TestServerConfig(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confServer))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	got := appconfig.GetServerConfig()

	if got.GetPort() != "8080" {
		t.Errorf("incorrect server config PORT, got: %v want: %v", got.GetPort(), "8080")
	}

	if got.GetDomain() != "localhost" {
		t.Errorf("incorrect server config domain, got: %v want: %v", got.GetDomain(), "localhost")
	}

	if got.GetProtocol() != "https://" {
		t.Errorf("incorrect server config protocol, got: %v want: %v", got.GetProtocol(), "https://")
	}

	if got.GetURI() != "https://localhost:8080" {
		t.Errorf("incorrect server config URI, got: %v want: %v", got.GetURI(), "https://localhost:8080")
	}
}

func TestClientForCurrency(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	sc, e := appconfig.ClientForCurrency("EUR")
	if e != nil {
		t.Errorf("error creating Stripe client for EUR: %v", e)
	}
	if sc == nil {
		t.Error("expected a non-nil Stripe client for EUR")
	}

	sc, e = appconfig.ClientForCurrency("NOT_FOUND_CURRENCY")
	if e != nil {
		t.Errorf("error creating Stripe client for default currency: %v", e)
	}
	if sc == nil {
		t.Error("expected a non-nil Stripe client for default currency")
	}
}

func TestClientForCurrencyWithoutDefault(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithoutDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	_, e = appconfig.ClientForCurrency("NOT_FOUND_CURRENCY")
	if e == nil {
		t.Error("configuration without Stripe default API keys must return an error")
	}
}

func TestClientForCurrencyOpWithoutDefault(t *testing.T) {
	e := appconfig.ImportConfig(strings.NewReader(confStripeWithoutDefault))
	if e != nil {
		t.Errorf("error during the config import: %v", e)
	}

	_, e = appconfig.ClientForCurrencyOp("NOT_FOUND_CURRENCY", "payment intent capture")
	if e == nil {
		t.Error("configuration without Stripe default API keys must return an error")
	}
	if e != nil && !strings.Contains(e.Error(), "payment intent capture") {
		t.Errorf("operation context missing from error, got: %v", e)
	}
}
