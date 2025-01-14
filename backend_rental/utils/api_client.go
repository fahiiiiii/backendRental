
// utils/api_client.go
package utils

import (
    "context"
    "io/ioutil"
    "net/http"
    "time"
    "golang.org/x/time/rate"
)

type APIClientInterface interface {
    MakeRequest(ctx context.Context, url string) ([]byte, error)
}

type APIClient struct {
    httpClient  *http.Client
    rateLimiter *rate.Limiter
    apiKey      string
}

func NewAPIClient(apiKey string) *APIClient {
    return &APIClient{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        rateLimiter: rate.NewLimiter(rate.Every(12*time.Second), 1),
        apiKey:      apiKey,
    }
}

func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("X-RapidAPI-Key", c.apiKey)
    req.Header.Add("X-RapidAPI-Host", "booking-com18.p.rapidapi.com")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}