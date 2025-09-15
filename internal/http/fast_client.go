package http

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/nepriyatelev/outline-client-go/outline"
	"github.com/valyala/fasthttp"
)

const (
	defaultUserAgentName       = "outline-go-client/1.0" // User-Agent header
	defaultMaxConnsPerHost     = 256                     // 256 connections
	defaultMaxIdleConnDuration = 30 * time.Second        // 30 seconds
	defaultReadTimeout         = 5 * time.Second         // 5 seconds
	defaultWriteTimeout        = 5 * time.Second         // 5 seconds
	defaultMaxConnDuration     = 1 * time.Minute         // 1 minute
	defaultReadBufferSize      = 4096                    // 4KB
	defaultWriteBufferSize     = 4096                    // 4KB
)

type Client struct {
	client *fasthttp.Client
}

func NewClient() *Client {
	fc := &fasthttp.Client{
		Name:                defaultUserAgentName,
		MaxConnsPerHost:     defaultMaxConnsPerHost,
		MaxIdleConnDuration: defaultMaxIdleConnDuration,
		ReadTimeout:         defaultReadTimeout,
		WriteTimeout:        defaultWriteTimeout,
		MaxConnDuration:     defaultMaxConnDuration,
		ReadBufferSize:      defaultReadBufferSize,
		WriteBufferSize:     defaultWriteBufferSize,
	}

	return &Client{
		client: fc,
	}
}

func (c *Client) Do(ctx context.Context, req *outline.Request) (*outline.Response, error) {

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
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		fastReq.SetBody(body)
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
	resp := &outline.Response{
		StatusCode: fastResp.StatusCode(),
		Headers:    headers,
		Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
	}
	return resp, nil
}
