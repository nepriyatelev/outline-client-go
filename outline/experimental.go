package outline

import (
	"context"
	"net/http"
	"time"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// === Experimental Metrics ===

func (c *Client) GetExperimentalMetrics(ctx context.Context, since time.Duration) (
	*types.ExperimentalMetricsResponse, error,
) {
	requestURL := *c.getExperimentalMetricsPath
	sinceQueryParamName := "since"
	q := requestURL.Query()
	q.Set(sinceQueryParamName, formatDuration(since))
	requestURL.RawQuery = q.Encode()

	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     requestURL.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "GetExperimentalMetrics", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	return unmarshalJSONWithError[types.ExperimentalMetricsResponse](resp.Body)
}
