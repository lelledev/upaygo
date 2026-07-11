package apperror

import "errors"

// InvalidError is a client-facing validation / bad-request error.
// Use Invalid to construct one so HTTP handlers can map it to 400.
type InvalidError struct {
	Msg string
}

func (e *InvalidError) Error() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

// Invalid returns an *InvalidError with the given message.
func Invalid(msg string) error {
	return &InvalidError{Msg: msg}
}

// IsInvalid reports whether err is or wraps an *InvalidError.
func IsInvalid(err error) bool {
	var inv *InvalidError
	return errors.As(err, &inv)
}
