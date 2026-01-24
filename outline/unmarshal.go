package outline

import (
	"encoding/json"
	"fmt"
)

// unmarshalJSONWithError unmarshals JSON data into a new instance of type T.
// It returns a pointer to the unmarshaled value or an error if unmarshaling fails.
func unmarshalJSONWithError[T any](data []byte) (*T, error) {
	target := new(T)
	if err := unmarshalWithErrorInternal(data, target, fmt.Sprintf("%T", target)); err != nil {
		return nil, err
	}
	return target, nil
}

// unmarshalAccessKeysResponse unmarshals the access keys response from JSON.
// It extracts the accessKeys array from the response wrapper.
func unmarshalAccessKeysResponse[T any](data []byte) ([]*T, error) {
	var wrapper struct {
		AccessKeys []*T `json:"accessKeys"`
	}
	if err := unmarshalWithErrorInternal(data, &wrapper, fmt.Sprintf("[]*%T", *new(T))); err != nil {
		return nil, err
	}
	return wrapper.AccessKeys, nil
}

// unmarshalWithErrorInternal performs the actual JSON unmarshaling with error handling.
// It checks for empty data and wraps unmarshaling errors with additional context.
func unmarshalWithErrorInternal(data []byte, target any, typeStr string) error {
	if len(data) == 0 {
		return errUnmarshalEmptyBody(typeStr)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return errUnmarshal(data, typeStr, err)
	}
	return nil
}
