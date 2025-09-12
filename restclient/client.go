package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		AuthToken:  token,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) newRequest(method, endpoint string, payload any) ([]byte, error) {
	var body io.Reader

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, err
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	return c.newRequest(http.MethodGet, endpoint, nil)
}

func (c *Client) Post(endpoint string, payload any) ([]byte, error) {
	return c.newRequest(http.MethodPost, endpoint, payload)
}

func (c *Client) Put(endpoint string, payload any) ([]byte, error) {
	return c.newRequest(http.MethodPut, endpoint, payload)
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	return c.newRequest(http.MethodDelete, endpoint, nil)
}
