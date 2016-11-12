package dota2api

import "net/http"

// Client contains the Dota 2 Web API calls as methods
type Client struct {
	key string
	hc  *http.Client
}

// NewClient instantiates a new client for the Dota 2 Web API
func NewClient(key string, options ...func(*Client)) *Client {
	c := &Client{
		key: key,
		hc:  http.DefaultClient,
	}
	for _, o := range options {
		o(c)
	}
	return c
}

// HTTPClient is an option to set the underlying HTTP client to use
func HTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.hc = hc
	}
}
