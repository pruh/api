package utils

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	c HTTPClient
}

func (c *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}

func NewHttpClient() HTTPClient {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return &defaultHTTPClient{
		c: client,
	}
}
