package outline

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaskSecretPath(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		secret   string
		expected string
	}{
		{
			name:     "Secret present in path",
			raw:      "/api/secret123/data",
			secret:   "secret123",
			expected: "/api/*****/data",
		},
		{
			name:     "Secret not present in path",
			raw:      "/api/data",
			secret:   "secret123",
			expected: "/api/data",
		},
		{
			name:     "Empty secret",
			raw:      "/api/secret123/data",
			secret:   "",
			expected: "/api/secret123/data",
		},
		{
			name:     "Multiple occurrences",
			raw:      "/secret123/foo/secret123/bar/",
			secret:   "secret123",
			expected: "/*****/foo/*****/bar/",
		},
		{
			name:     "Secret at start and end",
			raw:      "/secret123/",
			secret:   "secret123",
			expected: "/*****/",
		},
		{
			name:     "Secret with special characters",
			raw:      "/api/se$cret/data",
			secret:   "se$cret",
			expected: "/api/*****/data",
		},
		{
			name:     "Secret as substring of another segment",
			raw:      "/api/secret1234/data",
			secret:   "secret123",
			expected: "/api/secret1234/data",
		},
		{
			name:     "Secret at the end without trailing slash",
			raw:      "/api/data/secret123",
			secret:   "secret123",
			expected: "/api/data/secret123",
		},
		{
			name:     "Secret at the beginning without trailing slash",
			raw:      "/secret123",
			secret:   "secret123",
			expected: "/secret123",
		},
		{
			name:     "Path with no slashes",
			raw:      "secret123",
			secret:   "secret123",
			expected: "secret123",
		},
		{
			name:     "Secret surrounded by multiple slashes",
			raw:      "//secret123//foo//",
			secret:   "secret123",
			expected: "//*****//foo//",
		},
		{
			name:     "Secret with unicode characters",
			raw:      "/api/секрет/data",
			secret:   "секрет",
			expected: "/api/*****/data",
		},
		{
			name:     "Secret is a single character",
			raw:      "/a/b/c/a/d",
			secret:   "a",
			expected: "/*****/b/c/*****/d",
		},
		{
			name:     "Secret is a slash",
			raw:      "/api//data/",
			secret:   "/",
			expected: "/api//data/",
		},
		{
			name:     "Raw is empty",
			raw:      "",
			secret:   "secret123",
			expected: "",
		},
		{
			name:     "Secret is whitespace",
			raw:      "/api/ /data/",
			secret:   " ",
			expected: "/api/*****/data/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskSecretPath(tt.raw, tt.secret)
			assert.Equal(t, tt.expected, result, "maskSecretPath(%q, %q)", tt.raw, tt.secret)
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{
			name:     "Zero duration",
			input:    0,
			expected: "0s",
		},
		{
			name:     "Less than a minute",
			input:    42 * time.Second,
			expected: "42s",
		},
		{
			name:     "Exactly one minute",
			input:    1 * time.Minute,
			expected: "1m",
		},
		{
			name:     "Minutes but less than an hour",
			input:    17 * time.Minute,
			expected: "17m",
		},
		{
			name:     "Exactly one hour",
			input:    1 * time.Hour,
			expected: "1h",
		},
		{
			name:     "Hours and minutes (should show only hours)",
			input:    2*time.Hour + 30*time.Minute,
			expected: "2h",
		},
		{
			name:     "Negative duration",
			input:    -5 * time.Second,
			expected: "-5s",
		},
		{
			name:     "Large duration (multiple hours)",
			input:    25 * time.Hour,
			expected: "25h",
		},
		{
			name:     "Duration just under an hour",
			input:    59*time.Minute + 59*time.Second,
			expected: "59m",
		},
		{
			name:     "Duration just under a minute",
			input:    59 * time.Second,
			expected: "59s",
		},
		// New test cases
		{
			name:     "Negative duration, more than a minute",
			input:    -2 * time.Minute,
			expected: "-2m",
		},
		{
			name:     "Negative duration, more than an hour",
			input:    -3 * time.Hour,
			expected: "-3h",
		},
		{
			name:     "Duration just over an hour",
			input:    1*time.Hour + 1*time.Second,
			expected: "1h",
		},
		{
			name:     "Duration just over a minute",
			input:    1*time.Minute + 1*time.Second,
			expected: "1m",
		},
		{
			name:     "Duration is exactly 1 second",
			input:    1 * time.Second,
			expected: "1s",
		},
		{
			name:     "Duration is negative and just under a minute",
			input:    -59 * time.Second,
			expected: "-59s",
		},
		{
			name:     "Duration is negative and just under an hour",
			input:    -59*time.Minute - 1*time.Second,
			expected: "-59m",
		},
		{
			name:     "Duration is negative and just over an hour",
			input:    -1*time.Hour - 1*time.Second,
			expected: "-1h",
		},
		{
			name:     "Duration is negative and just over a minute",
			input:    -1*time.Minute - 1*time.Second,
			expected: "-1m",
		},
		{
			name:     "Very large duration (1000 hours)",
			input:    1000 * time.Hour,
			expected: "1000h",
		},
		{
			name:     "Very small negative duration (-1ns)",
			input:    -1 * time.Nanosecond,
			expected: "0s",
		},
		{
			name:     "Very small positive duration (1ns)",
			input:    1 * time.Nanosecond,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.input)
			assert.Equal(t, tt.expected, result, "formatDuration(%v)", tt.input)
		})
	}
}

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

	ue, ok := err.(*UnmarshalError)
	if assert.True(t, ok, "error should be of type *UnmarshalError") {
		assert.Equal(t, []byte{}, ue.Data)
		assert.Equal(t, fmt.Sprintf("%T", new(testPerson)), ue.Type)
		assert.EqualError(t, ue.Err, "empty body")
	}
}

func TestUnmarshalJSONWithError_InvalidJSON(t *testing.T) {
	data := []byte("{invalid")
	res, err := unmarshalJSONWithError[testPerson](data)
	assert.Nil(t, res)
	assert.Error(t, err)

	ue, ok := err.(*UnmarshalError)
	if assert.True(t, ok) {
		assert.Equal(t, data, ue.Data)
		// underlying json error message contains "invalid character"
		assert.Contains(t, ue.Err.Error(), "invalid character")
	}
}

func TestUnmarshalJSONWithError_TypeMismatch(t *testing.T) {
	// JSON string cannot be unmarshaled into int
	data := []byte(`"notanint"`)
	res, err := unmarshalJSONWithError[int](data)
	assert.Nil(t, res)
	assert.Error(t, err)

	ue, ok := err.(*UnmarshalError)
	if assert.True(t, ok) {
		assert.Equal(t, data, ue.Data)
		assert.Contains(t, ue.Err.Error(), "cannot unmarshal")
	}
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
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, []byte{}, ue.Data)
	assert.Equal(t, "int", ue.Type)
	assert.EqualError(t, ue.Err, "empty body")
}

func TestUnmarshalWithErrorInternal_InvalidJSON(t *testing.T) {
	var target map[string]interface{}
	data := []byte(`{"foo":123,}`)
	err := unmarshalWithErrorInternal(data, &target, "map[string]interface{}")
	assert.Error(t, err)
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, data, ue.Data)
	assert.Equal(t, "map[string]interface{}", ue.Type)
	assert.Contains(t, ue.Err.Error(), "invalid character")
}

func TestUnmarshalWithErrorInternal_TypeMismatch(t *testing.T) {
	var target int
	data := []byte(`"notanint"`)
	err := unmarshalWithErrorInternal(data, &target, "int")
	assert.Error(t, err)
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, data, ue.Data)
	assert.Equal(t, "int", ue.Type)
	assert.Contains(t, ue.Err.Error(), "cannot unmarshal")
}

func TestUnmarshalJSONSliceOfPointersWithError_Success(t *testing.T) {
	data := []byte(`[{"name":"Alice","age":30},{"name":"Bob","age":25}]`)
	res, err := unmarshalJSONSliceOfPointersWithError[testPerson](data)
	assert.NoError(t, err)
	if assert.Len(t, res, 2) {
		assert.Equal(t, "Alice", res[0].Name)
		assert.Equal(t, 30, res[0].Age)
		assert.Equal(t, "Bob", res[1].Name)
		assert.Equal(t, 25, res[1].Age)
	}
}

func TestUnmarshalJSONSliceOfPointersWithError_EmptySlice(t *testing.T) {
	data := []byte(`[]`)
	res, err := unmarshalJSONSliceOfPointersWithError[testPerson](data)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)
}

func TestUnmarshalJSONSliceOfPointersWithError_EmptyData(t *testing.T) {
	res, err := unmarshalJSONSliceOfPointersWithError[testPerson]([]byte{})
	assert.Nil(t, res)
	assert.Error(t, err)
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, []byte{}, ue.Data)
	assert.Contains(t, ue.Type, "[]*outline.testPerson")
	assert.EqualError(t, ue.Err, "empty body")
}

func TestUnmarshalJSONSliceOfPointersWithError_InvalidJSON(t *testing.T) {
	data := []byte(`[{"name":"Alice","age":30},]`)
	res, err := unmarshalJSONSliceOfPointersWithError[testPerson](data)
	assert.Nil(t, res)
	assert.Error(t, err)
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, data, ue.Data)
	assert.Contains(t, ue.Type, "[]*outline.testPerson")
	assert.Contains(t, ue.Err.Error(), "invalid character")
}

func TestUnmarshalJSONSliceOfPointersWithError_TypeMismatch(t *testing.T) {
	// Not a slice
	data := []byte(`{"name":"Alice","age":30}`)
	res, err := unmarshalJSONSliceOfPointersWithError[testPerson](data)
	assert.Nil(t, res)
	assert.Error(t, err)
	ue, ok := err.(*UnmarshalError)
	assert.True(t, ok)
	assert.Equal(t, data, ue.Data)
	assert.Contains(t, ue.Type, "[]*outline.testPerson")
	assert.Contains(t, ue.Err.Error(), "cannot unmarshal")
}
