package apppaymentintent

import (
	"time"

	appamount "github.com/lelledev/upaygo/amount"
	appcustomer "github.com/lelledev/upaygo/customer"
	apppaymentsource "github.com/lelledev/upaygo/payment/source"
)

type Intent interface {
	GetGatewayReference() string
	GetConfirmationMethod() string
	GetNextAction() string
	IsOffSession() bool
	GetCreatedAt() time.Time
	GetCustomer() appcustomer.Customer
	GetSource() apppaymentsource.Source
	GetAmount() appamount.Amount

	// From the status
	IsCanceled() bool
	IsSucceeded() bool
	RequiresCapture() bool
	RequiresConfirmation() bool
}

type i struct {
	R    string    `json:"gateway_reference"`   // Gateway reference
	CM   string    `json:"confirmation_method"` // Confirmation method
	NA   string    `json:"next_action"`         // Next Action
	OFFS bool      `json:"off_session"`         // Off session
	CT   time.Time `json:"created_at"`          // Create At

	C  appcustomer.Customer    `json:"customer"`
	PS apppaymentsource.Source `json:"payment_source"`
	A  appamount.Amount        `json:"amount"`
	S  status                  `json:"status"`
}

// GetGatewayReference exposes i.R value
func (i *i) GetGatewayReference() string {
	return i.R
}

// GetGatewayReference exposes i.CM value
func (i *i) GetConfirmationMethod() string {
	return i.CM
}

// GetGatewayReference exposes i.NA value
func (i *i) GetNextAction() string {
	return i.NA
}

// GetGatewayReference exposes i.OFFS value
func (i *i) IsOffSession() bool {
	return i.OFFS
}

// GetGatewayReference exposes i.CT. value
func (i *i) GetCreatedAt() time.Time {
	return i.CT
}

// GetGatewayReference exposes i.C value
func (i *i) GetCustomer() appcustomer.Customer {
	return i.C
}

// GetGatewayReference exposes i.PS value
func (i *i) GetSource() apppaymentsource.Source {
	return i.PS
}

// GetGatewayReference exposes i.A value
func (i *i) GetAmount() appamount.Amount {
	return i.A
}

// IsCanceled if status.s is canceled
func (i *i) IsCanceled() bool {
	return i.S.S == canceled
}

// IsSucceeded if status.s is succeeded
func (i *i) IsSucceeded() bool {
	return i.S.S == succeeded
}

// RequiresCapture if status.s is requirescapture
func (i *i) RequiresCapture() bool {
	return i.S.S == requirescapture
}

// RequiresConfirmation if status.s is requiresconfirmation
func (i *i) RequiresConfirmation() bool {
	return i.S.S == requiresconfirmation
}
