// Package outline provides a client for interacting with the Outline server API,
// including server configuration, access key management, metrics, and experimental endpoints.
package outline

import (
	"net/url"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/internal/http"
	"github.com/nepriyatelev/outline-client-go/internal/logger"
)

// Client manages authenticated calls to the Outline server API.
// The zero value is not usable; use [NewClient] or [MustNewClient] to create an instance.
// Client is safe for concurrent use after construction.
type Client struct {
	secret string

	// Server endpoints
	//
	// Get Server Information
	getServerInfoPath *url.URL

	// Server Configuration
	putServerHostnamePath *url.URL
	putServerPortPath     *url.URL
	putServerNamePath     *url.URL
	getMetricsEnabledPath *url.URL
	putMetricsEnabledPath *url.URL

	// Data Limits (Server-wide)
	putServerAccessKeyDataLimitPath    *url.URL
	deleteServerAccessKeyDataLimitPath *url.URL

	// Access keys endpoints
	//
	// CRUD Operations
	postAccessKeyPath   *url.URL
	getAccessKeysPath   *url.URL
	getAccessKeyPath    *url.URL
	putAccessKeyPath    *url.URL
	deleteAccessKeyPath *url.URL

	// Access Key Management
	putAccessKeyNamePath         *url.URL
	putAccessKeyDataLimitPath    *url.URL
	deleteAccessKeyDataLimitPath *url.URL

	// Metrics Endpoints
	//
	// Transfer Metrics
	getMetricsTransferPath *url.URL

	// Experimental Endpoints
	//
	// Experimental Metrics
	getExperimentalMetricsPath *url.URL

	// Internal
	doer   contracts.Doer
	logger contracts.Logger
}

// NewClient creates a [Client] that targets baseURL with the provided secret
// and applies the supplied options.
//
// It returns [*ParseURLError] if the baseURL cannot be parsed or joined with the secret.
func NewClient(baseURL, secret string, options ...Option) (*Client, error) {
	return initClient(baseURL, secret, options...)
}

// MustNewClient behaves like [NewClient] but panics on configuration errors.
func MustNewClient(baseURL, secret string, options ...Option) *Client {
	c, err := initClient(baseURL, secret, options...)
	if err != nil {
		panic(err)
	}

	return c
}

func initClient(baseURL, secret string, options ...Option) (*Client, error) {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return nil, errParseBaseURL(baseURL, err)
	}
	parsedBase.Path, err = url.JoinPath(parsedBase.Path, secret)
	if err != nil {
		return nil, errParseBaseURL(baseURL, err)
	}

	resolve := func(p string) *url.URL {
		return parsedBase.ResolveReference(&url.URL{Path: p})
	}

	var (
		// Server endpoints
		//
		// Get Server Information
		getServerInfoPath = "/server"

		// Server Configuration
		putServerHostnamePath = "/server/hostname-for-access-keys"
		putServerPortPath     = "/server/port-for-new-access-keys"
		putServerNamePath     = "/name"
		getMetricsEnabledPath = "/metrics/enabled"
		putMetricsEnabledPath = "/metrics/enabled"

		// Data Limits (Server-wide)
		putServerAccessKeyDataLimitPath    = "/server/access-key-data-limit"
		deleteServerAccessKeyDataLimitPath = "/server/access-key-data-limit"

		// Access keys endpoints
		//
		// CRUD Operations
		postAccessKeyPath   = "/access-keys"
		getAccessKeysPath   = "/access-keys"
		getAccessKeyPath    = "/access-keys/{id}"
		putAccessKeyPath    = "/access-keys/{id}"
		deleteAccessKeyPath = "/access-keys/{id}"

		// Access Key Management
		putAccessKeyNamePath         = "/access-keys/{id}/name"
		putAccessKeyDataLimitPath    = "/access-keys/{id}/data-limit"
		deleteAccessKeyDataLimitPath = "/access-keys/{id}/data-limit"

		// Metrics Endpoints
		//
		// Transfer Metrics
		getMetricsTransferPath = "/metrics/transfer"

		// Experimental Endpoints
		//
		// Experimental Metrics
		getExperimentalMetricsPath = "/experimental/server/metrics"
	)

	c := &Client{
		secret: secret,

		// Server endpoints
		getServerInfoPath:                  resolve(getServerInfoPath),
		putServerHostnamePath:              resolve(putServerHostnamePath),
		putServerPortPath:                  resolve(putServerPortPath),
		putServerNamePath:                  resolve(putServerNamePath),
		getMetricsEnabledPath:              resolve(getMetricsEnabledPath),
		putMetricsEnabledPath:              resolve(putMetricsEnabledPath),
		putServerAccessKeyDataLimitPath:    resolve(putServerAccessKeyDataLimitPath),
		deleteServerAccessKeyDataLimitPath: resolve(deleteServerAccessKeyDataLimitPath),

		// Access keys endpoints
		postAccessKeyPath:            resolve(postAccessKeyPath),
		getAccessKeysPath:            resolve(getAccessKeysPath),
		getAccessKeyPath:             resolve(getAccessKeyPath),
		putAccessKeyPath:             resolve(putAccessKeyPath),
		deleteAccessKeyPath:          resolve(deleteAccessKeyPath),
		putAccessKeyNamePath:         resolve(putAccessKeyNamePath),
		putAccessKeyDataLimitPath:    resolve(putAccessKeyDataLimitPath),
		deleteAccessKeyDataLimitPath: resolve(deleteAccessKeyDataLimitPath),

		// Metrics Endpoints
		getMetricsTransferPath: resolve(getMetricsTransferPath),

		// Experimental Endpoints
		getExperimentalMetricsPath: resolve(getExperimentalMetricsPath),

		doer:   http.NewClient(),
		logger: logger.NewNoopLogger(),
	}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}
