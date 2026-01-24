package outline

import (
	"testing"

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
			expected: "/api/data/*****",
		},
		{
			name:     "Secret at the beginning without trailing slash",
			raw:      "/secret123",
			secret:   "secret123",
			expected: "/*****",
		},
		{
			name:     "Path with no slashes",
			raw:      "secret123",
			secret:   "secret123",
			expected: "*****",
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

func TestSetIDInPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		id       string
		expected string
	}{
		{
			name:     "Replace {id} in middle",
			path:     "/api/{id}/data",
			id:       "123",
			expected: "/api/123/data",
		},
		{
			name:     "No {id} in path",
			path:     "/api/data",
			id:       "123",
			expected: "/api/data",
		},
		{
			name:     "Multiple {id}, replace only first",
			path:     "/{id}/foo/{id}/bar",
			id:       "123",
			expected: "/123/foo/{id}/bar",
		},
		{
			name:     "{id} at start",
			path:     "{id}/data",
			id:       "123",
			expected: "123/data",
		},
		{
			name:     "{id} at end",
			path:     "/api/{id}",
			id:       "123",
			expected: "/api/123",
		},
		{
			name:     "Empty path",
			path:     "",
			id:       "123",
			expected: "",
		},
		{
			name:     "Empty id",
			path:     "/api/{id}",
			id:       "",
			expected: "/api/",
		},
		{
			name:     "ID with special characters",
			path:     "/api/{id}",
			id:       "a/b?c=d&e=f",
			expected: "/api/a/b?c=d&e=f",
		},
		{
			name:     "No replacement needed",
			path:     "no id here",
			id:       "123",
			expected: "no id here",
		},
		{
			name:     "{id} in word",
			path:     "prefix{id}suffix",
			id:       "123",
			expected: "prefix123suffix",
		},
		{
			name:     "ID with numbers",
			path:     "/user/{id}/profile",
			id:       "456",
			expected: "/user/456/profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setIDInPath(tt.path, tt.id)
			assert.Equal(t, tt.expected, result, "setIDInPath(%q, %q)", tt.path, tt.id)
		})
	}
}
