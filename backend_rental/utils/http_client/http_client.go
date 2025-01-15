package http_client

import (
    "context"
    "encoding/json"
    "net/http"
)

type HTTPClient interface {
    Get(ctx context.Context, url string) (*http.Response, error)
}

type Client struct {
    client *http.Client
}

func NewHTTPClient() HTTPClient {
    return &Client{
        client: &http.Client{},
    }
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Add("X-RapidAPI-Key", "your-api-key")
    req.Header.Add("X-RapidAPI-Host", "booking-com18.p.rapidapi.com")
    
    return c.client.Do(req)
}