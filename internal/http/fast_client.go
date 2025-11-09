package http

import (
	"context"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/valyala/fasthttp"
)

const defaultUserAgentName = "outline-go-client/1.0" // User-Agent header

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

	bodyBytes := fastResp.Body()
	resp := &contracts.Response{
		StatusCode: fastResp.StatusCode(),
		Headers:    headers,
		Body:       bodyBytes,
	}
	return resp, nil
}
