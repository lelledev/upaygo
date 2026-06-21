package appconfig

import (
	"encoding/json"
	"fmt"
	"io"
)

var s config

// ImportConfig reads the config json reader and save the content in s global var
func ImportConfig(r io.Reader) error {
	b, e := io.ReadAll(r)
	if e != nil {
		return fmt.Errorf("impossible to read configuration: %v", e)
	}

	s = config{} // Reset existing configs
	e = json.Unmarshal(b, &s)
	if e != nil {
		return fmt.Errorf("impossible to unmarshal configuration: %v", e)
	}

	return nil
}

// public properties needed for json.Unmarshal
type config struct {
	Stripe apiKeys `json:"stripe"`
	Server server  `json:"server"`
}
