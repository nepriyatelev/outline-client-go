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
