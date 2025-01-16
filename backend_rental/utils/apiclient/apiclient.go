package apiclient

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "log"
    "time"
    "backend_rental/utils/ratelimiter"
)

type APIClient struct {
    client      *http.Client
    rateLimit   *ratelimiter.APIRateLimiter
    rapidAPIKey string
}

func NewAPIClient(rapidAPIKey string) *APIClient {
    return &APIClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        rateLimit:   ratelimiter.GetInstance(),
        rapidAPIKey: rapidAPIKey,
    }
}

func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
    if err := c.rateLimit.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit error: %v", err)
    }

    start := time.Now()
    defer func() {
        log.Printf("Request to %s took %v", url, time.Since(start))
    }()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code: %d, body: %s",
            resp.StatusCode, string(body))
    }

    return body, nil
}

func (c *APIClient) MakePostRequest(ctx context.Context, url string, body io.Reader) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "POST", url, body)
    if err != nil {
        return nil, fmt.Errorf("error creating POST request: %v", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", c.rapidAPIKey)
    req.Header.Add("Content-Type", "application/json")

    if err := c.rateLimit.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit error: %v", err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending POST request: %v", err)
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API POST request failed with status code: %d, body: %s",
            resp.StatusCode, string(responseBody))
    }

    return responseBody, nil
}

func (c *APIClient) MakeRequestWithRetry(ctx context.Context, url string) ([]byte, error) {
    var lastErr error
    for retries := 0; retries < 3; retries++ {
        data, err := c.MakeRequest(ctx, url)
        if err == nil {
            return data, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(retries+1))
    }
    return nil, fmt.Errorf("all retries failed: %v", lastErr)
}
// // utils/apiclient/apiclient.go
// package apiclient

// import (
//     "context"
//     "fmt"
//     "io"
//     "net/http"
//     "log"
//     "time"
//     "backend_rental/utils/ratelimiter"
// )

// type APIClient struct {
//     client      *http.Client
//     rateLimit *ratelimiter.APIRateLimiter
//     rapidAPIKey string
// }

// func NewAPIClient(rapidAPIKey string) *APIClient {
//     return &APIClient{
//         client: &http.Client{
//             Timeout: 10 * time.Second,
//         },
//         rateLimiter: ratelimiter.GetInstance(),
//         rapidAPIKey: rapidAPIKey,
//     }
// }
// // func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
// func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
//     if err := c.rateLimit.Wait(ctx); err != nil {
//         return nil, fmt.Errorf("rate limit error: %v", err)
//     }
// start := time.Now()
//     defer func() {
//         log.Printf("Request to %s took %v", url, time.Since(start))
//     }()
// // func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
//     req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     // Wait for rate limiter
//     if err := c.rateLimiter.Wait(ctx); err != nil {
//         return nil, fmt.Errorf("rate limiter error: %v", err)
//     }

//     resp, err := c.client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error sending request: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response body: %v", err)
//     }

//     if resp.StatusCode != http.StatusOK {
//         return nil, fmt.Errorf("API request failed with status code: %d, body: %s",
//             resp.StatusCode, string(body))
//     }

//     return body, nil
// }

// // Additional methods can be added as needed
// func (c *APIClient) MakePostRequest(ctx context.Context, url string, body io.Reader) ([]byte, error) {
//     req, err := http.NewRequestWithContext(ctx, "POST", url, body)
//     if err != nil {
//         return nil, fmt.Errorf("error creating POST request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)
//     req.Header.Add("Content-Type", "application/json")

//     // Wait for rate limiter
//     if err := c.rateLimiter.Wait(ctx); err != nil {
//         return nil, fmt.Errorf("rate limiter error: %v", err)
//     }

//     resp, err := c.client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error sending POST request: %v", err)
//     }
//     defer resp.Body.Close()

//     responseBody, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response body: %v", err)
//     }

//     if resp.StatusCode != http.StatusOK {
//         return nil, fmt.Errorf("API POST request failed with status code: %d, body: %s",
//             resp.StatusCode, string(responseBody))
//     }

//     return responseBody, nil
// }

// func (c *APIClient) MakeRequestWithRetry(ctx context.Context, url string) ([]byte, error) {
//     var lastErr error
//     for retries := 0; retries < 3; retries++ {
//         data, err := c.MakeRequest(ctx, url)
//         if err == nil {
//             return data, nil
//         }
//         lastErr = err
//         time.Sleep(time.Second * time.Duration(retries+1))
//     }
//     return nil, fmt.Errorf("all retries failed: %v", lastErr)
// }