package outline

import (
	"encoding/json"
	"fmt"
)

func unmarshalJSONWithError[T any](data []byte) (*T, error) {
	target := new(T)
	if err := unmarshalWithErrorInternal(data, target, fmt.Sprintf("%T", target)); err != nil {
		return nil, err
	}
	return target, nil
}

func unmarshalJSONSliceOfPointersWithError[T any](data []byte) ([]*T, error) {
	var target []*T
	if err := unmarshalWithErrorInternal(data, &target, fmt.Sprintf("%T", target)); err != nil {
		return nil, err
	}
	return target, nil
}

func unmarshalWithErrorInternal(data []byte, target any, typeStr string) error {
	if len(data) == 0 {
		return errUnmarshalEmptyBody(typeStr)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return errUnmarshal(data, typeStr, err)
	}
	return nil
}
