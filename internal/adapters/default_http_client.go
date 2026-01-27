package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DefaultHttpClient struct {
	client  *http.Client
	baseURL string
	token   string
}

func NewDefaultHttpClient(baseURL string, token string) *DefaultHttpClient {
	return &DefaultHttpClient{
		client:  &http.Client{},
		baseURL: baseURL,
		token:   token,
	}
}

func (c *DefaultHttpClient) Post(url string, req any, resp any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", httpResp.StatusCode)
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}
