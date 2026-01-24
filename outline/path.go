package outline

import "strings"

func maskSecretPath(raw, secret string) string {
	if secret == "" {
		return raw
	}

	// Split by slash to get path segments
	parts := strings.Split(raw, "/")

	// Mask each segment that exactly matches the secret
	for i, part := range parts {
		if part == secret {
			parts[i] = "*****"
		}
	}

	return strings.Join(parts, "/")
}

func setIDInPath(path string, id string) string {
	replacedPath := strings.Replace(path, "{id}", id, 1)
	return replacedPath
}
