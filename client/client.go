package client

import (
	"io"
	"net/http"
	"time"
)

// Client is a wrapper around http.Client that provides
// simpler methods for making HTTP requests and parsing responses.
type Client struct {
	*http.Client
	Timeout time.Duration
}

// NewClient creates a new HTTP client with an optional timeout.
// If timeout == 0, it uses the default http.Client timeout behavior.
func NewClient(timeout time.Duration) *Client {
	httpClient := &http.Client{}
	if timeout > 0 {
		httpClient.Timeout = timeout
	}
	return &Client{
		Client:  httpClient,
		Timeout: timeout,
	}
}

// executeRequest creates and sends an HTTP request, returning a Response wrapper.
func (c *Client) executeRequest(method, baseURL string, params, headers map[string]string, body io.Reader) (*Response, error) {
	fullURL, err := buildURL(baseURL, params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{Response: resp}, nil
}

// Get sends a GET request with optional query parameters and headers.
func (c *Client) Get(baseURL string, params, headers map[string]string) (*Response, error) {
	return c.executeRequest(http.MethodGet, baseURL, params, headers, nil)
}

// Post sends a POST request with optional query parameters, headers, and body.
func (c *Client) Post(baseURL string, params, headers map[string]string, payload io.Reader, contentType string) (*Response, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if contentType != "" {
		headers["Content-Type"] = contentType
	}
	return c.executeRequest(http.MethodPost, baseURL, params, headers, payload)
}
