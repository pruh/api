package utils

import (
	"net/http"
	"time"
)

// HTTPClient defines list of supported methods.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	c HTTPClient
}

// Do makes request and returns result or error.
func (c *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}

// NewHTTPClient creates new HTTP client.
func NewHTTPClient() HTTPClient {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return &defaultHTTPClient{
		c: client,
	}
}
