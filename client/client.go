package client

import (
	"net/url"
	"time"

	"github.com/jakub-gawlas/go-retryablehttp"
)

// Client for Fogger service REST API
type Client struct {
	URL  *url.URL
	http *retryablehttp.Client
}

// New returns default Fogger client
func New(url *url.URL) *Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = time.Second * 1
	httpClient.RetryWaitMax = time.Second * 2
	httpClient.RetryMax = 1

	return &Client{
		URL:  url,
		http: httpClient,
	}
}
