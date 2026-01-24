package outline

import (
	"context"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
)

// logRequest formats and sends two messages: Info and Debug.
// methodName — the name of the calling client function, e.g. "GetExperimentalMetrics".
// req — the final HTTP request.
func (c *Client) logRequest(ctx context.Context, methodName string, req *contracts.Request) {
	// Mask the secret in the Info log
	maskedURL := maskSecretPath(req.URL, c.secret)
	c.logger.Infof(
		ctx,
		"%s: sending request: method=%s url=%s headers=%v",
		methodName,
		req.Method,
		maskedURL,
		req.Headers,
	)
	// In the debug log, show the full URL
	c.logger.Debugf(
		ctx,
		"%s: sending request: method=%s url=%s headers=%v",
		methodName,
		req.Method,
		req.URL,
		req.Headers,
	)
}
