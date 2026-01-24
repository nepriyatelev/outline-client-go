package outline

import (
	"reflect"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
)

// Exported types from internal for users
type (
	Request  = contracts.Request
	Response = contracts.Response
	Doer     = contracts.Doer
	Logger   = contracts.Logger
)

// Option is a function that configures a Client.
type Option func(*Client)

// WithClient sets the HTTP client for the Client.
func WithClient(client Doer) Option {
	return func(c *Client) {
		if isNilInterface(client) {
			return
		}
		c.doer = client
	}
}

// WithLogger sets the logger for the Client.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		if isNilInterface(logger) {
			return
		}
		c.logger = logger
	}
}

// isNilInterface returns true if iface is nil
// or contains a dynamic nil pointer.
func isNilInterface(iface any) bool {
	if iface == nil {
		return true
	}
	v := reflect.ValueOf(iface)
	// If this is a pointer, the interface contains a nil pointer
	return v.Kind() == reflect.Ptr && v.IsNil()
}
