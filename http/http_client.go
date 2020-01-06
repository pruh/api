package http

import (
	"net/http"
	"time"
)

// Client defines list of supported methods.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	c Client
}

// Do makes request and returns result or error.
func (c *defaultHTTPClient) Do(r *http.Request) (*http.Response, error) {
	return c.c.Do(r)
}

// NewHTTPClient creates new HTTP client.
func NewHTTPClient() Client {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return &defaultHTTPClient{
		c: client,
	}
}
