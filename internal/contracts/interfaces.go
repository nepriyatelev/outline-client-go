package contracts

import (
	"context"
)

// Request — структура запроса
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// Response — структура ответа
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type Doer interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}