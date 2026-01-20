package outline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// === CRUD Operations for Access Keys ===

// CreateAccessKey Creates a new access key on the server.
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
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	default:
		return nil, errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) GetAccessKeys(ctx context.Context) ([]*types.AccessKey, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getAccessKeysPath.String(),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "GetAccessKeys", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONSliceOfPointersWithError[types.AccessKey](resp.Body)
	default:
		return nil, errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) GetAccessKey(ctx context.Context, accessKeyID string) (*types.AccessKey, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     setIDInPath(c.getAccessKeyPath.String(), accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "GetAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	case http.StatusNotFound:
		return nil, &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return nil, errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) UpdateAccessKey(ctx context.Context, accessKeyID string,
	updateAccessKey *types.AccessKey,
) (*types.AccessKey, error) {
	var reqBodyBytes []byte

	if updateAccessKey != nil {
		reqBodyBytes, _ = json.Marshal(updateAccessKey)
	}

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     setIDInPath(c.putAccessKeyPath.String(), accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return unmarshalJSONWithError[types.AccessKey](resp.Body)
	case http.StatusNotFound:
		return nil, &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return nil, errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) DeleteAccessKey(ctx context.Context, accessKeyID string) error {
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     setIDInPath(c.deleteAccessKeyPath.String(), accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "DeleteAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
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
		URL:     setIDInPath(c.putAccessKeyNamePath.String(), accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateNameAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

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
		URL:     setIDInPath(c.putAccessKeyDataLimitPath.String(), accessKeyID),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateDataLimitAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return &ClientError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf(
				"The provided data limit is invalid: %d.", bytes),
		}
	case http.StatusNotFound:
		return &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

func (c *Client) DeleteDataLimitAccessKey(ctx context.Context, accessKeyID string) error {
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     setIDInPath(c.deleteAccessKeyDataLimitPath.String(), accessKeyID),
		Headers: DefaultHeaders(),
	}

	c.logRequest(ctx, "DeleteDataLimitAccessKey", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return &ClientError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Access key with ID %s not found.", accessKeyID),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}
