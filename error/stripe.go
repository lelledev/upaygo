package apperror

import (
	"encoding/json"

	"github.com/stripe/stripe-go/v76"
)

// Stripe returns a built-in error
// stripeError.Error() returns a JSON string with multiple keys
// New intent error examples
// {"code":"resource_missing","status":400,"message":"No such customer: cus_xxx","param":"customer","request_id":"req_8PAgbbyTbIufgS","type":"invalid_request_error"}
// {"status":401,"message":"Invalid API Key provided: sk_xxx","type":"invalid_request_error"}

// GetStripeErrorMessage extracts the message from the built-in Stripe error
func GetStripeErrorMessage(e error) (string, error) {
	var es stripe.Error

	ej := json.Unmarshal([]byte(e.Error()), &es)
	if ej != nil {
		return "", ej
	}

	return es.Msg, nil
}
