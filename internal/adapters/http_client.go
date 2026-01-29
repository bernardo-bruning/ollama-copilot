package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpClient interface {
	Post(url string, req any, resp any) error
}

type DefaultHttpClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func NewHttpClient(baseURL string, headers map[string]string) HttpClient {
	return &DefaultHttpClient{
		client:  &http.Client{},
		baseURL: baseURL,
		headers: headers,
	}
}

func (c *DefaultHttpClient) Post(path string, req any, resp any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for key, value := range c.headers {
		httpReq.Header.Set(key, value)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	if resp == nil {
		return nil
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}