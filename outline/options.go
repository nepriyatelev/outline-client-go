package outline

import (
	"reflect"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
)

// Экспортируем типы из internal для пользователей
type (
	Request  = contracts.Request
	Response = contracts.Response
	Doer     = contracts.Doer
	Logger   = contracts.Logger
)

// Option задаёт опцию настройки клиента.
type Option func(*Client)

func WithClient(client Doer) Option {
	return func(c *Client) {
		if isNilInterface(client) {
			return
		}
		c.doer = client
	}
}

// WithLogger задаёт логгер для клиента.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		if isNilInterface(logger) {
			return
		}
		c.logger = logger
	}
}

// isNilInterface возвращает true, если iface равен nil
// или содержит динамический нулевой указатель.
func isNilInterface(iface any) bool {
	if iface == nil {
		return true
	}
	v := reflect.ValueOf(iface)
	// Если это указатель, интерфейс содержит nil-указатель
	return v.Kind() == reflect.Ptr && v.IsNil()
}
