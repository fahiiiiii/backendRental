// controllers/property.go
package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    // "io"
    "log"
    // "net/http"
    "net/url"
    "strings"
    "sync"
    "time"
    "math"

    beego "github.com/beego/beego/v2/server/web"
    // "github.com/beego/beego/v2/server/web/context"
    "golang.org/x/time/rate"
    "backend_rental/models"
    "backend_rental/utils/apiclient" 
    
)

type PropertyController struct {
    beego.Controller
    cityProperties  map[string][]string
    mutex           sync.Mutex
    rateLimiter    *rate.Limiter
    rapidAPIKey     string
	apiClient      *apiclient.APIClient
}

func NewPropertyController() *PropertyController {
    rapidAPIKey, _ := beego.AppConfig.String("rapidapikey")
    if rapidAPIKey == "" {
        log.Fatal("rapidapikey is not set in app.conf")
    }

    return &PropertyController{
        cityProperties: make(map[string][]string),
        apiClient:     apiclient.NewAPIClient(rapidAPIKey),
    }
}
// func NewPropertyController() *PropertyController {
//     rapidAPIKey, _ := beego.AppConfig.String("rapidapikey")
//     if rapidAPIKey == "" {
//         log.Fatal("rapidapikey is not set in app.conf")
//     }

//     return &PropertyController{
//         cityProperties: make(map[string][]string),
//         rateLimiter:   rate.NewLimiter(rate.Every(12*time.Second), 1),
//         rapidAPIKey:   rapidAPIKey,
//     }
// }

// controllers/property.go

func init() {
    // Create a new controller instance
    controller := NewPropertyController()

    // Start background processing in a goroutine
    go func() {
        // Create a ticker for periodic cleanup and maintenance
        cleanupTicker := time.NewTicker(24 * time.Hour)
        defer cleanupTicker.Stop()

        // Create error channel for handling background processing errors
        errorCh := make(chan error, 1)
        defer close(errorCh)

        // Create context with cancellation for cleanup
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        for {
            select {
            case <-ctx.Done():
                return
            case <-cleanupTicker.C:
                // Perform periodic cleanup of stale data
                controller.mutex.Lock()
                for city, properties := range controller.cityProperties {
                    // Remove cities with no properties or expired data
                    if len(properties) == 0 {
                        delete(controller.cityProperties, city)
                    }
                }
                controller.mutex.Unlock()

                // Reset rate limiter if needed
                controller.rateLimiter = rate.NewLimiter(rate.Every(12*time.Second), 1)

            case err := <-errorCh:
                // Log any errors from background processing
                log.Printf("PropertyController background processing error: %v", err)
            }
        }
    }()
}


// Get handles all GET requests
func (c *PropertyController) Get() {
    action := c.GetString("action")
    
    var result interface{}
    var err error

    switch action {
    case "fetch":
        cityName := c.GetString("city")
        country := c.GetString("country")
        if cityName == "" || country == "" {
            c.handleError("city and country parameters are required")
            return
        }
        err = c.fetchPropertiesForCity(cityName, country)
    case "summary":
        result = c.getSummary()
    default:
        c.handleError("invalid action")
        return
    }

    if err != nil {
        c.handleError(err.Error())
        return
    }

    c.Data["json"] = result
    c.ServeJSON()
}

func (c *PropertyController) getSummary() *models.PropertySummary {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    return &models.PropertySummary{
        TotalProperties: c.getTotalProperties(),
        CitiesCount:     len(c.cityProperties),
        CityProperties:  c.cityProperties,
    }
}

func (c *PropertyController) getTotalProperties() int {
    total := 0
    for _, props := range c.cityProperties {
        total += len(props)
    }
    return total
}

func (c *PropertyController) fetchPropertiesForCity(cityName, country string) error {
    properties, err := c.fetchPropertiesWithRetry(cityName, country, 3)
    if err != nil {
        return err
    }

    result := struct {
        City       models.CityKey
        Properties []models.Property
        Err        error
    }{
        City:       models.CityKey{Name: cityName, Country: country},
        Properties: properties,
    }

    c.processPropertyResult(result)
    return nil
}

func (c *PropertyController) fetchPropertiesWithRetry(cityName, country string, maxRetries int) ([]models.Property, error) {
    for attempt := 0; attempt < maxRetries; attempt++ {
        properties, err := c.doFetchProperties(cityName, country)
        
        if err == nil {
            return properties, nil
        }
        
        if strings.Contains(err.Error(), "Too many requests") || 
           strings.Contains(err.Error(), "You are not subscribed") {
            waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
            time.Sleep(waitTime)
            continue
        }
        
        return nil, err
    }
    
    return nil, fmt.Errorf("failed to fetch properties after %d attempts", maxRetries)
}

func (c *PropertyController) doFetchProperties(cityName, country string) ([]models.Property, error) {
    uniqueProperties := make(map[string]models.Property)
    searchQueries := []string{
        cityName,
        fmt.Sprintf("%s hotels", cityName),
        fmt.Sprintf("%s accommodation", cityName),
    }

    for _, query := range searchQueries {
        encodedQuery := url.QueryEscape(query)
        apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)
        
        properties, err := c.makeAPIRequest(apiURL)
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

// Update makeAPIRequest method
func (c *PropertyController) makeAPIRequest(apiURL string) ([]models.Property, error) {
    body, err := c.apiClient.MakeRequest(context.Background(), apiURL)
    if err != nil {
        return nil, err
    }

    var response models.PropertyResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, err
    }

    return response.Data, nil
}

// func (c *PropertyController) makeAPIRequest(apiURL string) ([]models.Property, error) {
//     if err := c.rateLimiter.Wait(context.Background()); err != nil {
//         return nil, fmt.Errorf("rate limiter error: %v", err)
//     }

//     req, err := http.NewRequest("GET", apiURL, nil)
//     if err != nil {
//         return nil, err
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     client := &http.Client{Timeout: 10 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, err
//     }

//     if resp.StatusCode != http.StatusOK {
//         return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
//     }

//     var response models.PropertyResponse
//     if err := json.Unmarshal(body, &response); err != nil {
//         return nil, err
//     }

//     return response.Data, nil
// }

func (c *PropertyController) processPropertyResult(result struct {
    City       models.CityKey
    Properties []models.Property
    Err        error
}) {
    if len(result.Properties) == 0 {
        return
    }

    c.mutex.Lock()
    defer c.mutex.Unlock()

    c.cityProperties[result.City.Name] = []string{}
    
    maxProperties := 20
    if len(result.Properties) < maxProperties {
        maxProperties = len(result.Properties)
    }
    
    for _, prop := range result.Properties[:maxProperties] {
        c.cityProperties[result.City.Name] = append(
            c.cityProperties[result.City.Name], 
            prop.Name,
        )
    }
}

func (c *PropertyController) handleError(message string) {
    c.Data["json"] = map[string]string{"error": message}
    c.Ctx.Output.SetStatus(400)
    c.ServeJSON()
}