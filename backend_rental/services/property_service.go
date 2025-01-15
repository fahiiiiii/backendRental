// services/property_service.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "sync"
    
    "github.com/beego/beego/v2/core/logs"  // Changed to correct import
    "backend_rental/models"
    "backend_rental/utils/ratelimiter"
    
)

type PropertyServiceInterface interface {
    GetAllProperties(ctx context.Context) (map[string][]string, error)
}

type PropertyService struct {
    uniqueCities   map[string][]string
    cityProperties map[string][]string
    mutex          sync.RWMutex
    httpClient     *http.Client
    rapidAPIKey    string
}

func NewPropertyService(rapidAPIKey string) PropertyServiceInterface {
    return &PropertyService{
        uniqueCities:   make(map[string][]string),
        cityProperties: make(map[string][]string),
        httpClient:     &http.Client{},
        rapidAPIKey:    rapidAPIKey,
    }
}

// fetchPropertyData handles the API request and response parsing
// func (s *PropertyService) fetchPropertyData(apiURL string) ([]models.Property, error) {
//     req, err := http.NewRequest("GET", apiURL, nil)
//     if err != nil {
//         return nil, err
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", s.rapidAPIKey)

//     resp, err := s.httpClient.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, err
//     }

//     if resp.StatusCode == 429 || strings.Contains(string(body), "Too many requests") {
//         return nil, fmt.Errorf("rate limit exceeded")
//     }

//     var response struct {
//         Data []struct {
//             DestID   string `json:"dest_id"`
//             Name     string `json:"name"`
//             CityName string `json:"city_name"`
//         } `json:"data"`
//     }

//     err = json.Unmarshal(body, &response)
//     if err != nil {
//         return nil, err
//     }

//     properties := make([]models.Property, 0, len(response.Data))
//     for _, item := range response.Data {
//         properties = append(properties, models.Property{
//             DestID:   item.DestID,
//             Name:     item.Name,
//             CityName: item.CityName,
//         })
//     }

//     return properties, nil
// }
func (s *PropertyService) fetchPropertyData(apiURL string) ([]models.Property, error) {
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", s.rapidAPIKey)

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode == 429 || strings.Contains(string(body), "Too many requests") {
        return nil, fmt.Errorf("rate limit exceeded")
    }

    var response struct {
        Data []struct {
            DestID string `json:"dest_id"`
            Name   string `json:"name"`
            CityID string  `json:"city_id"`

            // Remove CityName from response struct if not needed
        } `json:"data"`
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        return nil, err
    }

    properties := make([]models.Property, 0, len(response.Data))
    for _, item := range response.Data {
        properties = append(properties, models.Property{
            DestID: item.DestID,
            Name:   item.Name,
            CityID: item.CityID,
            // Remove CityName if not part of your model
        })
    }

    return properties, nil
}
func (s *PropertyService) GetAllProperties(ctx context.Context) (map[string][]string, error) {
    propertyResults := make(chan struct {
        City       models.CityKey
        Properties []models.Property
        Err        error
    }, len(s.uniqueCities))
    
    var wg sync.WaitGroup
    rateLimiter := ratelimiter.GetInstance()
    
    for country, cities := range s.uniqueCities {
        for _, cityName := range cities {
            wg.Add(1)
            go func(city, country string) {
                defer wg.Done()
                
                if err := rateLimiter.Wait(ctx); err != nil {
                    propertyResults <- struct {
                        City       models.CityKey
                        Properties []models.Property
                        Err        error
                    }{
                        City: models.CityKey{Name: city, Country: country},
                        Err:  err,
                    }
                    return
                }
                
                properties, err := s.fetchPropertiesWithRetry(ctx, city, country, 3)
                propertyResults <- struct {
                    City       models.CityKey
                    Properties []models.Property
                    Err        error
                }{
                    City:       models.CityKey{Name: city, Country: country},
                    Properties: properties,
                    Err:        err,
                }
            }(cityName, country)
        }
    }
    
    go func() {
        wg.Wait()
        close(propertyResults)
    }()
    
    for result := range propertyResults {
        if result.Err != nil {
            logs.Error("Error fetching properties for %s, %s: %v",
                result.City.Name, result.City.Country, result.Err)
            continue
        }
        s.processPropertyResult(result)
    }
    
    return s.cityProperties, nil
}

func (s *PropertyService) fetchPropertiesWithRetry(ctx context.Context, cityName, country string, maxRetries int) ([]models.Property, error) {
    var lastErr error
    for attempt := 0; attempt < maxRetries; attempt++ {
        properties, err := s.fetchPropertiesForCity(ctx, cityName, country)
        if err == nil {
            return properties, nil
        }
        lastErr = err
    }
    return nil, fmt.Errorf("failed to fetch properties after %d attempts: %v", maxRetries, lastErr)
}

func (s *PropertyService) fetchPropertiesForCity(ctx context.Context, cityName, country string) ([]models.Property, error) {
    uniqueProperties := make(map[string]models.Property)
    searchQueries := []string{
        cityName,
        fmt.Sprintf("%s hotels", cityName),
        fmt.Sprintf("%s accommodation", cityName),
    }
    
    for _, query := range searchQueries {
        encodedQuery := url.QueryEscape(query)
        apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)
        properties, err := s.fetchPropertyData(apiURL)
        if err != nil {
            continue
        }
        for _, prop := range properties {
            if prop.DestID != "" {
                uniqueProperties[prop.DestID] = prop
            }
        }
    }
    
    result := make([]models.Property, 0, len(uniqueProperties))
    for _, prop := range uniqueProperties {
        result = append(result, prop)
    }
    return result, nil
}

func (s *PropertyService) processPropertyResult(result struct {
    City       models.CityKey
    Properties []models.Property
    Err        error
}) {
    if len(result.Properties) == 0 {
        return
    }
    
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.cityProperties[result.City.Name] = []string{}
    maxProperties := 20
    if len(result.Properties) < maxProperties {
        maxProperties = len(result.Properties)
    }
    
    for _, prop := range result.Properties[:maxProperties] {
        s.cityProperties[result.City.Name] = append(
            s.cityProperties[result.City.Name],
            prop.Name,
        )
    }
}





// package services

// import (
//     "context"
//     "fmt"
//     "net/url"
//     "sync"
    
//     "backend_rental/models"
//     "backend_rental/utils/ratelimiter"
//     "github.com/beego/beego/v2/core/logs"
// )

// type PropertyServiceInterface interface {
//     GetAllProperties(ctx context.Context) (map[string][]string, error)
// }

// type PropertyService struct {
//     uniqueCities   map[string][]string
//     cityProperties map[string][]string
//     mutex          sync.RWMutex
// }

// func NewPropertyService() PropertyServiceInterface {
//     return &PropertyService{
//         uniqueCities:   make(map[string][]string),
//         cityProperties: make(map[string][]string),
//     }
// }

// func (s *PropertyService) GetAllProperties(ctx context.Context) (map[string][]string, error) {
//     propertyResults := make(chan struct {
//         City       models.CityKey
//         Properties []models.Property
//         Err        error
//     }, len(s.uniqueCities))
    
//     var wg sync.WaitGroup
//     rateLimiter := ratelimiter.GetInstance()
    
//     for country, cities := range s.uniqueCities {
//         for _, cityName := range cities {
//             wg.Add(1)
//             go func(city, country string) {
//                 defer wg.Done()
                
//                 if err := rateLimiter.Wait(ctx); err != nil {
//                     propertyResults <- struct {
//                         City       models.CityKey
//                         Properties []models.Property
//                         Err        error
//                     }{
//                         City: models.CityKey{Name: city, Country: country},
//                         Err:  err,
//                     }
//                     return
//                 }
                
//                 properties, err := s.fetchPropertiesWithRetry(ctx, city, country, 3)
//                 propertyResults <- struct {
//                     City       models.CityKey
//                     Properties []models.Property
//                     Err        error
//                 }{
//                     City:       models.CityKey{Name: city, Country: country},
//                     Properties: properties,
//                     Err:        err,
//                 }
//             }(cityName, country)
//         }
//     }
    
//     go func() {
//         wg.Wait()
//         close(propertyResults)
//     }()
    
//     for result := range propertyResults {
//         if result.Err != nil {
//             logs.Error("Error fetching properties for %s, %s: %v",
//                 result.City.Name, result.City.Country, result.Err)
//             continue
//         }
//         s.processPropertyResult(result)
//     }
    
//     return s.cityProperties, nil
// }

// func (s *PropertyService) fetchPropertiesWithRetry(ctx context.Context, cityName, country string, maxRetries int) ([]models.Property, error) {
//     var lastErr error
//     for attempt := 0; attempt < maxRetries; attempt++ {
//         properties, err := s.fetchPropertiesForCity(ctx, cityName, country)
//         if err == nil {
//             return properties, nil
//         }
//         lastErr = err
//     }
//     return nil, fmt.Errorf("failed to fetch properties after %d attempts: %v", maxRetries, lastErr)
// }

// func (s *PropertyService) fetchPropertiesForCity(ctx context.Context, cityName, country string) ([]models.Property, error) {
//     uniqueProperties := make(map[string]models.Property)
//     searchQueries := []string{
//         cityName,
//         fmt.Sprintf("%s hotels", cityName),
//         fmt.Sprintf("%s accommodation", cityName),
//     }
    
//     for _, query := range searchQueries {
//         encodedQuery := url.QueryEscape(query)
//         apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)
//         properties, err := s.fetchPropertyData(ctx, apiURL)
//         if err != nil {
//             continue
//         }
//         for _, prop := range properties {
//             if prop.DestID != "" {
//                 uniqueProperties[prop.DestID] = prop
//             }
//         }
//     }
    
//     result := make([]models.Property, 0, len(uniqueProperties))
//     for _, prop := range uniqueProperties {
//         result = append(result, prop)
//     }
//     return result, nil
// }

// func (s *PropertyService) processPropertyResult(result struct {
//     City       models.CityKey
//     Properties []models.Property
//     Err        error
// }) {
//     if len(result.Properties) == 0 {
//         return
//     }
    
//     s.mutex.Lock()
//     defer s.mutex.Unlock()
    
//     s.cityProperties[result.City.Name] = []string{}
//     maxProperties := 20
//     if len(result.Properties) < maxProperties {
//         maxProperties = len(result.Properties)
//     }
    
//     for _, prop := range result.Properties[:maxProperties] {
//         s.cityProperties[result.City.Name] = append(
//             s.cityProperties[result.City.Name],
//             prop.Name,
//         )
//     }
// }