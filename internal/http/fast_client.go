package http

import (
	"context"
	"slices"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/valyala/fasthttp"
)

const defaultUserAgentName = "outline-go-client/1.0" // User-Agent header

// Client is a fasthttp-based HTTP client that implements the contracts.Doer interface.
//
// Memory Usage Considerations:
// The client uses fasthttp's pooled Request/Response objects for efficiency.
// However, since pooled responses are released immediately after the Do() call completes,
// the response body MUST be copied out of the pool before returning. This means:
//   - Each response body is fully cloned into new memory allocation
//   - Large response bodies may cause significant memory usage
//   - The response body remains in memory until the caller is done with it
//
// For typical Outline API responses (JSON metadata, access keys, metrics), this is not
// a concern as responses are typically small (<100KB). If you expect large responses,
// consider implementing streaming or chunked processing.
type Client struct {
	client *fasthttp.Client
}

func NewClient() *Client {
	fc := &fasthttp.Client{
		Name: defaultUserAgentName,
	}

	return &Client{
		client: fc,
	}
}

func (c *Client) Do(ctx context.Context, req *contracts.Request) (*contracts.Response, error) {
	fastReq := fasthttp.AcquireRequest()
	fastResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(fastReq)
	defer fasthttp.ReleaseResponse(fastResp)

	// Устанавливаем URI и метод
	fastReq.SetRequestURI(req.URL)
	fastReq.Header.SetMethod(req.Method)

	// Устанавливаем заголовки
	for key, value := range req.Headers {
		fastReq.Header.Set(key, value)
	}

	// Устанавливаем тело запроса, если оно есть
	if req.Body != nil {
		fastReq.SetBody(req.Body)
	}

	// Запускаем фактический HTTP-запрос в отдельной горутине
	errCh := make(chan error, 1)
	go func() {
		errCh <- c.client.Do(fastReq, fastResp)
	}()

	// Ждём либо завершения запроса, либо отмены контекста
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		// При отмене контекста возвращаем её ошибку
		return nil, ctx.Err()
	}

	// Преобразуем fasthttp.Response в наш Response
	headers := make(map[string]string, fastResp.Header.Len())
	fastResp.Header.All()(func(key, value []byte) bool {
		// Копируем key и value, так как они могут быть перезаписаны
		headers[string(key)] = string(value)
		return true // продолжаем итерацию
	})

	// IMPORTANT: Body cloning is required here because fasthttp uses pooled Response objects.
	// The fastResp will be released back to the pool via defer, so we MUST copy the body bytes
	// out of the pooled response before it's released. This ensures the returned Response
	// contains valid data after this function returns.
	//
	// Memory implications:
	// - The entire response body is allocated in new memory via slices.Clone()
	// - For large responses, this creates a full copy in heap memory
	// - The copied body persists until the caller releases the contracts.Response
	//
	// This is the correct behavior for fasthttp pooling, but callers should be aware
	// that holding onto Response objects will retain the full body in memory.
	bodyBytes := fastResp.Body()
	cloneBodyBytes := slices.Clone(bodyBytes)

	resp := &contracts.Response{
		StatusCode: fastResp.StatusCode(),
		Headers:    headers,
		Body:       cloneBodyBytes,
	}
	return resp, nil
}
