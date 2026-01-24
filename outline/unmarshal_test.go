package outline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test types for unmarshal tests
type testPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestUnmarshalJSONWithError_Success(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30}`)
	res, err := unmarshalJSONWithError[testPerson](data)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Alice", res.Name)
		assert.Equal(t, 30, res.Age)
	}
}

func TestUnmarshalJSONWithError_Empty(t *testing.T) {
	res, err := unmarshalJSONWithError[testPerson]([]byte{})
	assert.Nil(t, res)
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestUnmarshalJSONWithError_InvalidJSON(t *testing.T) {
	data := []byte("{invalid")
	res, err := unmarshalJSONWithError[testPerson](data)
	assert.Nil(t, res)
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestUnmarshalJSONWithError_TypeMismatch(t *testing.T) {
	// JSON string cannot be unmarshaled into int
	data := []byte(`"notanint"`)
	res, err := unmarshalJSONWithError[int](data)
	assert.Nil(t, res)
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestUnmarshalWithErrorInternal_Success(t *testing.T) {
	type testStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}
	data := []byte(`{"foo":"baz","bar":42}`)
	var target testStruct
	err := unmarshalWithErrorInternal(data, &target, "testStruct")
	assert.NoError(t, err)
	assert.Equal(t, "baz", target.Foo)
	assert.Equal(t, 42, target.Bar)
}

func TestUnmarshalWithErrorInternal_EmptyData(t *testing.T) {
	var target int
	err := unmarshalWithErrorInternal([]byte{}, &target, "int")
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestUnmarshalWithErrorInternal_InvalidJSON(t *testing.T) {
	var target map[string]interface{}
	data := []byte(`{"foo":123,}`)
	err := unmarshalWithErrorInternal(data, &target, "map[string]interface{}")
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestUnmarshalWithErrorInternal_TypeMismatch(t *testing.T) {
	var target int
	data := []byte(`"notanint"`)
	err := unmarshalWithErrorInternal(data, &target, "int")
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestUnmarshalAccessKeysResponse_Success(t *testing.T) {
	data := []byte(`{"accessKeys":[{"name":"Alice","age":30},{"name":"Bob","age":25}]}`)
	res, err := unmarshalAccessKeysResponse[testPerson](data)
	assert.NoError(t, err)
	if assert.Len(t, res, 2) {
		assert.Equal(t, "Alice", res[0].Name)
		assert.Equal(t, 30, res[0].Age)
		assert.Equal(t, "Bob", res[1].Name)
		assert.Equal(t, 25, res[1].Age)
	}
}

func TestUnmarshalAccessKeysResponse_EmptyArray(t *testing.T) {
	data := []byte(`{"accessKeys":[]}`)
	res, err := unmarshalAccessKeysResponse[testPerson](data)
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}

func TestUnmarshalAccessKeysResponse_MissingField(t *testing.T) {
	data := []byte(`{}`)
	res, err := unmarshalAccessKeysResponse[testPerson](data)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestUnmarshalAccessKeysResponse_EmptyData(t *testing.T) {
	res, err := unmarshalAccessKeysResponse[testPerson]([]byte{})
	assert.Nil(t, res)
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestUnmarshalAccessKeysResponse_InvalidJSON(t *testing.T) {
	data := []byte(`{"accessKeys":[invalid}`)
	res, err := unmarshalAccessKeysResponse[testPerson](data)
	assert.Nil(t, res)
	assert.Error(t, err)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}
