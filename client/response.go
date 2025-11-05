package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response wraps http.Response and caches the body
// for multiple reads and easier JSON/string parsing.
type Response struct {
	*http.Response
	cachedBody []byte
}

// getDataCopy safely reads the body once, closes it, and rebuilds it
// so that it can be read multiple times.
func (resp *Response) getDataCopy() ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("nil response or body")
	}
	if resp.cachedBody != nil {
		return resp.cachedBody, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	resp.cachedBody = bodyBytes

	return bodyBytes, nil
}

// Json parses the response body into a map[string]any.
// Returns an error if the body is not valid JSON.
func (r *Response) Json() (map[string]any, error) {
	body, err := r.getDataCopy()
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return data, nil
}

// String returns the response body as a string.
func (r *Response) String() (string, error) {
	body, err := r.getDataCopy()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Bytes returns the response body as a byte slice.
func (r *Response) Bytes() ([]byte, error) {
	return r.getDataCopy()
}

// IsSuccess reports whether the HTTP status code is in the 2xx range.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// StatusText returns the textual representation of the status code.
func (r *Response) StatusText() string {
	return http.StatusText(r.StatusCode)
}
