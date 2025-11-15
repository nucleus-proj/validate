package fluent_validator

import (
	"encoding/json"
	"errors"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// NewJSONError takes a list of strings and generates a json string error output
func NewJSONError(errs []string) error {
	if len(errs) == 0 {
		return nil
	}

	resp := ErrorResponse{Errors: errs}

	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return errors.New(string(b))
}

// NewJSONError takes a list of strings and generates a concatenated readable string error output
func NewErrorFromStrings(errs []string) error {
	if len(errs) == 0 {
		return nil
	}

	eList := make([]error, 0, len(errs))
	for _, msg := range errs {
		eList = append(eList, errors.New(msg))
	}

	return errors.Join(eList...)
}
