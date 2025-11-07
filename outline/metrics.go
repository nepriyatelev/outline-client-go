package outline

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

func (c *Client) GetMetricsTransfer(ctx context.Context) (*types.MetricsTransfer, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getMetricsTransferPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var metricsTransfer *types.MetricsTransfer
	if err = json.Unmarshal(resp.Body, metricsTransfer); err != nil {
		return nil, err
	}

	return metricsTransfer, nil
}
