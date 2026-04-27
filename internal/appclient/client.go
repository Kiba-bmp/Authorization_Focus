package appclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const internalKeyHeader = "X-Internal-Key"

type Client struct {
	baseURL     string
	internalKey string
	httpClient  *http.Client
}

type UserProvisionRequest struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func New(baseURL, internalKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		internalKey: internalKey,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

func (c *Client) ProvisionUser(ctx context.Context, req UserProvisionRequest) error {
	httpReq, err := c.newJSONRequest(ctx, http.MethodPost, "/internal/users/provision", req)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set(internalKeyHeader, c.internalKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if out == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(out)
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if len(body) == 0 {
			return fmt.Errorf("backend returned %s", resp.Status)
		}
		return fmt.Errorf("backend returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
}
