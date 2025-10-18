package outline

import (
	"net/url"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/internal/http"
	"github.com/nepriyatelev/outline-client-go/internal/logger"
)

var errParseBaseURL = func(baseURL string, err error) *ParseURLError {
	return &ParseURLError{
		BaseURL: baseURL,
		Err:     err,
	}
}

type Client struct {
	secret string

	// Server endpoints
	getServerInfoPath     *url.URL
	putServerHostnamePath *url.URL
	putServerNamePath     *url.URL
	getMetricsEnabledPath *url.URL
	putMetricsEnabledPath *url.URL
	
	// Access keys endpoints
	putServerPortPath         *url.URL
	putServerDataLimitPath    *url.URL
	deleteServerDataLimitPath *url.URL
	postAccessKeyPath         *url.URL
	getAccessKeysPath         *url.URL
	putAccessKeyPath          *url.URL
	getAccessKeyPath          *url.URL
	deleteAccessKeyPath       *url.URL
	putAccessKeyNamePath      *url.URL
	getMetricsTransferPath    *url.URL

	// Experimental endpoints
	getExperimentalMetricsPath *url.URL

	// Limit endpoints
	putServerAccessKeyDataLimitPath    *url.URL
	deleteServerAccessKeyDataLimitPath *url.URL
	putAccessKeyDataLimitPath          *url.URL
	deleteAccessKeyDataLimitPath       *url.URL

	// Internal
	doer                   contracts.Doer
	logger                 contracts.Logger
}

func NewClient(baseURL, secret string, options ...Option) (*Client, error) {
	return initClient(baseURL, secret, options...)
}

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
	getServerInfoPath     = "/server"
	putServerHostnamePath = "/server/hostname-for-access-keys"
	putServerNamePath     = "/name"
	getMetricsEnabledPath = "/metrics/enabled"
	putMetricsEnabledPath = "/metrics/enabled"

	// Access keys endpoints
	putServerPortPath         = "/server/port-for-new-access-keys"
	putServerDataLimitPath    = "/server/access-key-data-limit"
	deleteServerDataLimitPath = "/server/access-key-data-limit"
	postAccessKeyPath         = "/access-keys"
	getAccessKeysPath         = "/access-keys"
	putAccessKeyPath          = "/access-keys/{id}"
	getAccessKeyPath          = "/access-keys/{id}"
	deleteAccessKeyPath       = "/access-keys/{id}"
	putAccessKeyNamePath      = "/access-keys/{id}/name"
	getMetricsTransferPath    = "/metrics/transfer"

	// Experimental endpoints
	getExperimentalMetricsPath = "/experimental/server/metrics"

	// Limit endpoints
	putServerAccessKeyDataLimitPath    = "/server/access-key-data-limit"
	deleteServerAccessKeyDataLimitPath = "/server/access-key-data-limit"
	putAccessKeyDataLimitPath          = "/access-keys/{id}/data-limit"
	deleteAccessKeyDataLimitPath       = "/access-keys/{id}/data-limit"
)

	c := &Client{
		secret: secret,

		// Server endpoints
		getServerInfoPath:     resolve(getServerInfoPath),            
		putServerHostnamePath: resolve(putServerHostnamePath),
		putServerNamePath:     resolve(putServerNamePath),    
		getMetricsEnabledPath: resolve(getMetricsEnabledPath),
		putMetricsEnabledPath: resolve(putMetricsEnabledPath),

		// Access keys endpoints
		putServerPortPath:         resolve(putServerPortPath),
		putServerDataLimitPath:    resolve(putServerDataLimitPath),
		deleteServerDataLimitPath: resolve(deleteServerDataLimitPath),
		postAccessKeyPath:         resolve(postAccessKeyPath),
		getAccessKeysPath:         resolve(getAccessKeysPath),
		putAccessKeyPath:          resolve(putAccessKeyPath),
		getAccessKeyPath:          resolve(getAccessKeyPath),
		deleteAccessKeyPath:       resolve(deleteAccessKeyPath),
		putAccessKeyNamePath:      resolve(putAccessKeyNamePath),
		getMetricsTransferPath:    resolve(getMetricsTransferPath),

		// Experimental endpoints
		getExperimentalMetricsPath: resolve(getExperimentalMetricsPath),

		// Limit endpoints
		putServerAccessKeyDataLimitPath:    resolve(putServerAccessKeyDataLimitPath),
		deleteServerAccessKeyDataLimitPath: resolve(deleteServerAccessKeyDataLimitPath),
		putAccessKeyDataLimitPath:		  resolve(putAccessKeyDataLimitPath),
		deleteAccessKeyDataLimitPath:       resolve(deleteAccessKeyDataLimitPath),                    


		doer:   http.NewClient(),
		logger: logger.NewNoopLogger(),
	}

	for _, opt := range options {
		opt(c)
	}
	return c, nil
}
