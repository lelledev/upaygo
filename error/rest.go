package apperror

import (
	"encoding/json"
	"net/http"
)

// RESTError is the JSON body returned for failed HTTP requests.
type RESTError struct {
	M string `json:"error"`
}

// Error implements the error interface.
func (e *RESTError) Error() string {
	if e == nil {
		return ""
	}
	return e.M
}

// WriteJSON writes err as a JSON RESTError with Content-Type application/json
// and a status from HTTPStatus(err). Encoding failures are ignored after the
// status has been written (the client already received headers).
func WriteJSON(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatus(err))
	_ = json.NewEncoder(w).Encode(RESTError{M: Message(err)})
}
