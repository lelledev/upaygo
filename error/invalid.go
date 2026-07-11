package apperror

import "errors"

// InvalidError is a client-facing validation / bad-request error.
// Use Invalid or InvalidCause so HTTP handlers can map it to 400.
// When Err is set, Unwrap exposes it for errors.Is / errors.As.
type InvalidError struct {
	Msg string
	Err error
}

func (e *InvalidError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		if e.Msg == "" {
			return e.Err.Error()
		}
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

// Unwrap returns the underlying cause, if any.
func (e *InvalidError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Invalid returns an *InvalidError with the given message and no cause.
func Invalid(msg string) error {
	return &InvalidError{Msg: msg}
}

// InvalidCause returns an *InvalidError that wraps cause while keeping msg as
// the validation context. Prefer this over fmt.Sprintf so callers can still
// use errors.Is / errors.As on the original error.
func InvalidCause(msg string, cause error) error {
	return &InvalidError{Msg: msg, Err: cause}
}

// IsInvalid reports whether err is or wraps an *InvalidError.
func IsInvalid(err error) bool {
	var inv *InvalidError
	return errors.As(err, &inv)
}
