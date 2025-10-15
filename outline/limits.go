package outline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
)

var errInvalidDataLimit = func(bytes uint64) *ClientError {
	return &ClientError{
		Code: 400,
		Message: fmt.Sprintf("Invalid data limit: %d.",
			bytes),
	}
}

func (c *Client) SetAllKeyLimitBytes(ctx context.Context, bytes uint64) error {
	requestURL := *c.serverAccessKeyDataLimitURL

	var reqBody struct {
		Limit struct {
			Bytes uint64 `json:"bytes"`
		} `json:"limit"`
	}

	reqBody.Limit.Bytes = bytes

	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     requestURL.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "SetAllKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidDataLimit(bytes)
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) DeleteAllKeyLimitBytes(ctx context.Context) error {
	requestURL := *c.serverAccessKeyDataLimitURL

	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     requestURL.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "DeleteAllKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

var errAccessKeyInexistent = func(id uint64) *ClientError {
	return &ClientError{
		Code:    404,
		Message: fmt.Sprintf("Access key inexistent: %d.", id),
	}
}

func (c *Client) SetKeyLimitBytes(ctx context.Context, id uint64, bytes uint64) error {
	requestURL := *c.serverIndividualAccessKeyDataLimitURL

	setIDInPath(&requestURL, id)

	// TODO: повтор структуры
	var reqBody struct {
		Limit struct {
			Bytes uint64 `json:"bytes"`
		} `json:"limit"`
	}

	reqBody.Limit.Bytes = bytes

	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     requestURL.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "SetKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidDataLimit(bytes)
	case http.StatusNotFound:
		return errAccessKeyInexistent(id)
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) DeleteKeyLimitBytes(ctx context.Context, id uint64) error {
	requestURL := *c.serverIndividualAccessKeyDataLimitURL

	setIDInPath(&requestURL, id)

	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     requestURL.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "DeleteKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errAccessKeyInexistent(id)
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}
