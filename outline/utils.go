package outline

import (
    "strings"
)

func maskSecretPath(raw, secret string) string {
    if secret == "" {
		return raw
	}

	return strings.ReplaceAll(raw, "/"+secret+"/", "/*****/")
}
