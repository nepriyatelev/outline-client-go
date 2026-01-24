package outline

import (
	"net/url"
	"strings"
)

// maskSecretPath masks the secret part in the raw URL path by replacing it with *****.
// This is used for logging to avoid exposing sensitive information.
func maskSecretPath(raw, secret string) string {
	if secret == "" {
		return raw
	}

	parts := strings.Split(raw, "/")

	for i, part := range parts {
		if part == secret {
			parts[i] = "*****"
		}
	}

	return strings.Join(parts, "/")
}

// setIDInPath replaces the {id} placeholder in the URL path with the actual id.
// It returns the full URL string with the id substituted.
func setIDInPath(u url.URL, id string) string {
	replacedPath := strings.Replace(u.Path, "{id}", id, 1)
	u.Path = replacedPath
	return u.String()
}
