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

type ResponseError struct {
	StatusCode int
	Message    string
}

func (e *ResponseError) Error() string {
	return e.Message
}

type UserProvisionRequest struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func New(baseURL, internalKey string) *Client {
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		internalKey: internalKey,
		httpClient:  &http.Client{Timeout: time.Second * 10},
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
		return fmt.Errorf("ошибка запроса к backend: %w", err)
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
			return &ResponseError{
				StatusCode: resp.StatusCode,
				Message:    "backend вернул ошибку без описания.",
			}
		}

		var payload struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &payload); err == nil && strings.TrimSpace(payload.Error) != "" {
			return &ResponseError{
				StatusCode: resp.StatusCode,
				Message:    strings.TrimSpace(payload.Error),
			}
		}

		return &ResponseError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}
}
