package appcurrency

import (
	"fmt"
	"strings"
)

type Currency interface {
	GetISO4217() string

	Equal(c Currency) bool
}

type c struct {
	ISO4217 string `json:"ISO_4217"` // 3 Char symbol - ISO 4217
}

// GetISO4217 exposes c.ISO4217 value
func (c *c) GetISO4217() string {
	return c.ISO4217
}

// Equal checks if a == b
func (a *c) Equal(b Currency) bool {
	return a.GetISO4217() == b.GetISO4217()
}

// New returns a new instance of c
func New(iso4217 string) (Currency, error) {
	if len(iso4217) != 3 {
		return nil, fmt.Errorf("'%v' is not a currency ISO 4217 format", iso4217)
	}

	return &c{ISO4217: strings.ToUpper(iso4217)}, nil
}
