package outline

import (
	"context"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// === Transfer Metrics ===

func (c *Client) GetMetricsTransfer(ctx context.Context) (*types.MetricsTransfer, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getMetricsTransferPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "GetMetricsTransfer", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONWithError[types.MetricsTransfer](resp.Body)
	default:
		return nil, errUnexpected(resp.StatusCode, resp.Body)
	}
}
