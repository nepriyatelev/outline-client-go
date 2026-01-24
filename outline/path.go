package outline

import (
	"net/url"
	"strings"
)

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

func setIDInPath(u url.URL, id string) string {
	replacedPath := strings.Replace(u.Path, "{id}", id, 1)
	u.Path = replacedPath
	return u.String()
}
