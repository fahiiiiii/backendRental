// utils/api/api_utils.go
package api

import (
    "fmt"
    "net/url"
)

// BuildAPIURL constructs the API URL for city queries
func BuildAPIURL(query string) string {
    baseURL := "https://booking-com18.p.rapidapi.com/api/v1/cities"
    
    params := url.Values{}
    params.Add("query", query)
    params.Add("limit", "100")
    
    return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
