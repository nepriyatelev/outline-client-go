package outline

import (
	"context"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
)

// logRequest формирует и отправляет два сообщения: Info и Debug.
// methodName — имя вызываемой клиентской функции, например "GetExperimentalMetrics".
// req — конечный HTTP запрос.
func (c *Client) logRequest(ctx context.Context, methodName string, req *contracts.Request) {
	// Скрываем секрет в Info-логе
	maskedURL := maskSecretPath(req.URL, c.secret)
	c.logger.Infof(
		ctx,
		"%s: sending request: method=%s url=%s headers=%v",
		methodName,
		req.Method,
		maskedURL,
		req.Headers,
	)
	// В debug-логе показываем полный URL
	c.logger.Debugf(
		ctx,
		"%s: sending request: method=%s url=%s headers=%v",
		methodName,
		req.Method,
		req.URL,
		req.Headers,
	)
}
