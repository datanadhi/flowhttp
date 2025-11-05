package client

import (
	"net/url"
)

// buildURL appends query parameters to the base URL.
func buildURL(baseURL string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return baseURL, nil
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
