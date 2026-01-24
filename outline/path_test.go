package outline

import (
	"net/url"
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
		urlStr   string
		id       string
		expected string
	}{
		{
			name:     "Replace {id} in middle",
			urlStr:   "/api/{id}/data",
			id:       "123",
			expected: "/api/123/data",
		},
		{
			name:     "No {id} in path",
			urlStr:   "/api/data",
			id:       "123",
			expected: "/api/data",
		},
		{
			name:     "Multiple {id}, replace only first",
			urlStr:   "/{id}/foo/{id}/bar",
			id:       "123",
			expected: "/123/foo/%7Bid%7D/bar",
		},
		{
			name:     "{id} at start",
			urlStr:   "{id}/data",
			id:       "123",
			expected: "123/data",
		},
		{
			name:     "{id} at end",
			urlStr:   "/api/{id}",
			id:       "123",
			expected: "/api/123",
		},
		{
			name:     "Empty path",
			urlStr:   "",
			id:       "123",
			expected: "",
		},
		{
			name:     "Empty id",
			urlStr:   "/api/{id}",
			id:       "",
			expected: "/api/",
		},
		{
			name:     "ID with special characters",
			urlStr:   "/api/{id}",
			id:       "a/b?c=d&e=f",
			expected: "/api/a/b%3Fc=d&e=f",
		},
		{
			name:     "No replacement needed",
			urlStr:   "no id here",
			id:       "123",
			expected: "no%20id%20here",
		},
		{
			name:     "{id} in word",
			urlStr:   "prefix{id}suffix",
			id:       "123",
			expected: "prefix123suffix",
		},
		{
			name:     "ID with numbers",
			urlStr:   "/user/{id}/profile",
			id:       "456",
			expected: "/user/456/profile",
		},
		{
			name:     "Full URL with {id}",
			urlStr:   "http://localhost:8081/api/access-keys/{id}/data-limit",
			id:       "key-123",
			expected: "http://localhost:8081/api/access-keys/key-123/data-limit",
		},
		{
			name:     "HTTPS URL with {id}",
			urlStr:   "https://api.example.com/v1/resources/{id}",
			id:       "resource-456",
			expected: "https://api.example.com/v1/resources/resource-456",
		},
		{
			name:     "URL with query parameters",
			urlStr:   "http://example.com/api/{id}?param=value",
			id:       "789",
			expected: "http://example.com/api/789?param=value",
		},
		{
			name:     "URL with port and {id}",
			urlStr:   "http://localhost:3000/users/{id}/profile",
			id:       "user-abc",
			expected: "http://localhost:3000/users/user-abc/profile",
		},
		{
			name:     "URL with multiple path segments and {id}",
			urlStr:   "https://service.domain.com/api/v2/projects/{id}/issues",
			id:       "proj-001",
			expected: "https://service.domain.com/api/v2/projects/proj-001/issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.urlStr)
			assert.NoError(t, err, "url.Parse should not fail for %q", tt.urlStr)
			result := setIDInPath(*u, tt.id)
			assert.Equal(t, tt.expected, result, "setIDInPath(%q, %q)", tt.urlStr, tt.id)
		})
	}
}
