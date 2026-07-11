package apperror

import (
	"errors"
	"net/http"

	"github.com/stripe/stripe-go/v82"
)

// AsStripeError reports whether err is or wraps a *stripe.Error and returns it.
// Prefer this over re-parsing Error() JSON strings.
func AsStripeError(err error) (*stripe.Error, bool) {
	var se *stripe.Error
	if errors.As(err, &se) {
		return se, true
	}
	return nil, false
}

// Message returns a human-readable error string suitable for API responses.
// Stripe errors use their Msg field; otherwise err.Error() is returned.
func Message(err error) string {
	if err == nil {
		return ""
	}
	if se, ok := AsStripeError(err); ok && se.Msg != "" {
		return se.Msg
	}
	return err.Error()
}

// HTTPStatus maps err to an HTTP status code for REST responses.
//
// Mapping rules:
//   - *InvalidError → 400
//   - *stripe.Error with 401 (bad API key) → 500 (server misconfiguration)
//   - *stripe.Error with other 4xx → that status
//   - *stripe.Error with 5xx → 502 Bad Gateway
//   - everything else → 500
func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if IsInvalid(err) {
		return http.StatusBadRequest
	}
	if se, ok := AsStripeError(err); ok {
		switch {
		case se.HTTPStatusCode == http.StatusUnauthorized:
			// Invalid Stripe secret key is a server config problem.
			return http.StatusInternalServerError
		case se.HTTPStatusCode >= 400 && se.HTTPStatusCode < 500:
			return se.HTTPStatusCode
		case se.HTTPStatusCode >= 500:
			return http.StatusBadGateway
		}
	}
	return http.StatusInternalServerError
}
