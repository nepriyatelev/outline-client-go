package outline

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// === CRUD Operations for Access Keys ===

// CreateAccessKey creates a new access key on the server with the provided configuration.
// It returns the created access key or an error if the operation fails.
//
// It returns [*ClientError] for unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) CreateAccessKey(ctx context.Context, createAccessKey *types.CreateAccessKey) (
	*types.AccessKey, error,
) {
	var reqBodyBytes []byte

	if createAccessKey != nil {
		reqBodyBytes, _ = json.Marshal(createAccessKey)
	}

	req := &contracts.Request{
		Method:  http.MethodPost,
		URL:     c.postAccessKeyPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "CreateAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoCreateAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// GetAccessKeys retrieves all access keys from the server.
// It returns a slice of access keys or an error if the operation fails.
//
// It returns [*ClientError] for unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) GetAccessKeys(ctx context.Context) ([]*types.AccessKey, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getAccessKeysPath.String(),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "GetAccessKeys", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetAccessKeys(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalAccessKeysResponse[types.AccessKey](resp.Body)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// GetAccessKey retrieves a specific access key by its ID from the server.
// It returns the access key or an error if not found or if the operation fails.
//
// It returns [*ClientError] with code 404 if the access key is not found,
// [*ClientError] for other unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) GetAccessKey(ctx context.Context, accessKeyID string) (*types.AccessKey, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     setIDInPath(*c.getAccessKeyPath, accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "GetAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	case http.StatusNotFound:
		return nil, errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateAccessKey updates an existing access key with the provided data.
// It returns the updated access key or an error if not found or if the operation fails.
//
// It returns [*ClientError] with code 404 if the access key is not found,
// [*ClientError] for other unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateAccessKey(ctx context.Context, accessKeyID string,
	updateAccessKey *types.AccessKey,
) (*types.AccessKey, error) {
	var reqBodyBytes []byte

	if updateAccessKey != nil {
		reqBodyBytes, _ = json.Marshal(updateAccessKey)
	}

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     setIDInPath(*c.putAccessKeyPath, accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoUpdateAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	case http.StatusNotFound:
		return nil, errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// DeleteAccessKey deletes an access key by its ID from the server.
// It returns an error if the access key is not found or if the operation fails.
//
// It returns [*ClientError] with code 404 if the access key is not found,
// [*ClientError] for other unexpected HTTP status codes,
// or [*DoError] if the HTTP request fails.
func (c *Client) DeleteAccessKey(ctx context.Context, accessKeyID string) error {
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     setIDInPath(*c.deleteAccessKeyPath, accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "DeleteAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoDeleteAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// === Management Operations for Access Keys ===

func (c *Client) UpdateNameAccessKey(ctx context.Context, accessKeyID, newName string) error {
	var reqBody struct {
		Name string `json:"name"`
	}

	reqBody.Name = newName
	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     setIDInPath(*c.putAccessKeyNamePath, accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateNameAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateNameAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateDataLimitAccessKey sets a data transfer limit for an access key.
// It returns an error if the access key is not found, the limit is invalid, or if the operation fails.
//
// It returns [*ClientError] with code 400 if the data limit is invalid,
// [*ClientError] with code 404 if the access key is not found,
// [*ClientError] for other unexpected HTTP status codes,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateDataLimitAccessKey(
	ctx context.Context, accessKeyID string, bytes uint64,
) error {
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	reqBody.Limit.Bytes = bytes

	reqBodyBytes, _ := json.Marshal(reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     setIDInPath(*c.putAccessKeyDataLimitPath, accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateDataLimitAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateDataLimitAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidDataLimit(http.StatusBadRequest, bytes)
	case http.StatusNotFound:
		return errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// DeleteDataLimitAccessKey removes the data transfer limit for an access key.
// It returns an error if the access key is not found or if the operation fails.
//
// It returns [*ClientError] with code 404 if the access key is not found,
// [*ClientError] for other unexpected HTTP status codes,
// or [*DoError] if the HTTP request fails.
func (c *Client) DeleteDataLimitAccessKey(ctx context.Context, accessKeyID string) error {
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     setIDInPath(*c.deleteAccessKeyDataLimitPath, accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "DeleteDataLimitAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoDeleteDataLimitAccessKey(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errAccessKeyNotFound(http.StatusNotFound, accessKeyID)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}
