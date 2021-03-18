package api

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("https://%s@api.sncf.com/v1", token),
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
