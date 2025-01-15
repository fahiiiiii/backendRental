package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "golang.org/x/time/rate"
    "backend_rental/models"
    "backend_rental/services"
    "backend_rental/utils"
    "backend_rental/utils/apiclient" 
)



type BookingController struct {
    beego.Controller
    uniqueCountries map[string]bool
    uniqueCities    map[string]bool
    countryCities   map[string][]string
    cityProperties  map[string][]string
    mutex           sync.Mutex
    rateLimiter    *rate.Limiter
    rapidAPIKey     string
    fetchedCities  []models.City
    apiClient      *apiclient.APIClient
    locationService *services.LocationService
}

func NewBookingController() *BookingController {
    // Get RapidAPI key from app.conf
    rapidAPIKey, _ := beego.AppConfig.String("rapidapikey")
    if rapidAPIKey == "" {
        log.Fatal("rapidapikey is not set in app.conf")
    }

    return &BookingController{
        uniqueCountries: make(map[string]bool),
        uniqueCities:    make(map[string]bool),
        countryCities:   make(map[string][]string),
        cityProperties:  make(map[string][]string),
        rateLimiter:    rate.NewLimiter(rate.Every(12*time.Second), 1),
        rapidAPIKey:     rapidAPIKey,
        apiClient:      apiclient.NewAPIClient(rapidAPIKey),
        locationService: &services.LocationService{},
    }
}

func (c *BookingController) Prepare() {
    // Optional: Add any pre-request preparation logic
}

func init() {
    // Automated city processing on application startup
    go func() {
        controller := NewBookingController()
        err := controller.processAllCities()
        if err != nil {
            log.Printf("Automated city processing error: %v", err)
        }
    }()
}

func (c *BookingController) Get() {
    action := c.GetString("action")
    
    var result interface{}
    var err error

    switch action {
    case "summary":
        result, err = c.getSummary()
        if err != nil {
            c.handleError(err)
            return
        }

    case "process":
        err = c.processAllCities()
        if err != nil {
            c.handleError(err)
            return
        }
        result = map[string]string{"status": "success"}

    case "city":
        cityID := c.GetString("id")
        if cityID == "" {
            c.handleError(fmt.Errorf("city ID is required"))
            return
        }
        
        result, err = c.getCityDetails(cityID)
        if err != nil {
            c.handleError(err)
            return
        }

    case "list":
        // Get pagination parameters
        page, _ := c.GetInt("page", 1)
        pageSize, _ := c.GetInt("page_size", 20)
        country := c.GetString("country")
        cityName := c.GetString("city_name")

        result, err = c.listLocations(page, pageSize, country, cityName)
        if err != nil {
            c.handleError(err)
            return
        }

    default:
        c.handleError(fmt.Errorf("invalid action: %s", action))
        return
    }

    c.Data["json"] = result
    c.ServeJSON()
}

func (c *BookingController) getSummary() (*models.BookingSummary, error) {
    // Use LocationService to get unique countries and cities
    countryCities, err := c.locationService.GetUniqueCountriesAndCities()
    if err != nil {
        return nil, err
    }
    
    return &models.BookingSummary{
        Countries:      utils.ConvertToCountryMap(countryCities),
        Cities:         utils.ConvertToCityMap(countryCities),
        CountryCities:  countryCities,
        CityProperties: c.cityProperties,
    }, nil
}

func (c *BookingController) getCityDetails(cityID string) (*models.Location, error) {
    // Delegate to LocationService
    return c.locationService.GetLocationByID(cityID)
}

func (c *BookingController) listLocations(
    page, 
    pageSize int, 
    country, 
    cityName string,
) (map[string]interface{}, error) {
    // Delegate to LocationService
    locations, totalCount, err := c.locationService.GetLocations(page, pageSize, country, cityName)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "locations": locations,
        "pagination": map[string]interface{}{
            "total_count": totalCount,
            "page":        page,
            "page_size":   pageSize,
        },
    }, nil
}

// func (c *BookingController) processAllCities() error {
//     // Generate queries using utility function
//     queries := utils.GenerateLocationQueries()
    
//     var allCities []models.City
//     var mu sync.Mutex
//     var wg sync.WaitGroup
    
//     // Semaphore to limit concurrent requests
//     semaphore := make(chan struct{}, 5)
    
//     for _, query := range queries {
//         wg.Add(1)
//         semaphore <- struct{}{}
        
//         go func(q string) {
//             defer wg.Done()
//             defer func() { <-semaphore }()
            
//             cities, err := c.fetchCities(q)
//             if err != nil {
//                 log.Printf("Error fetching cities for query '%s': %v", q, err)
//                 return
//             }
            
//             // Filter and clean cities
//             cleanedCities := utils.FilterAndCleanCities(cities)
            
//             // Thread-safe append
//             mu.Lock()
//             allCities = append(allCities, cleanedCities...)
//             mu.Unlock()
//         }(query)
//     }
    
//     // Wait for all goroutines to complete
//     wg.Wait()
    
//     // Use LocationService to process and store cities
//     return c.locationService.ProcessAndStoreCities(allCities)
// }
func (c *BookingController) processAllCities() error {
    // Create service instances
    apiClient := apiclient.NewAPIClient(c.rapidAPIKey)
    locationProcessingService := services.NewLocationProcessingService(apiClient)
    
    // Generate queries
    queries := utils.GenerateLocationQueries()
    
    // Let the service handle the processing
    return locationProcessingService.ProcessLocationsFromQueries(queries)
}
// func (c *BookingController) processAllCities() error {
//     queries := utils.GenerateLocationQueries()
//     log.Printf("Starting to process %d queries", len(queries))
    
//     var allCities []models.City
//     var mu sync.Mutex
//     var wg sync.WaitGroup
//     errChan := make(chan error, len(queries))
    
//     // Semaphore for rate limiting
//     semaphore := make(chan struct{}, 3) // Reduced concurrent requests
    
//     for _, query := range queries {
//         wg.Add(1)
//         semaphore <- struct{}{}
        
//         go func(q string) {
//             defer wg.Done()
//             defer func() { <-semaphore }()
            
//             // Add retry logic
//             var cities []models.City
//             var err error
//             for retries := 0; retries < 3; retries++ {
//                 cities, err = c.fetchCities(q)
//                 if err == nil {
//                     break
//                 }
//                 log.Printf("Retry %d for query '%s': %v", retries+1, q, err)
//                 time.Sleep(time.Second * time.Duration(retries+1))
//             }
            
//             if err != nil {
//                 errChan <- fmt.Errorf("failed to fetch cities for query '%s': %v", q, err)
//                 return
//             }
            
//             cleanedCities := utils.FilterAndCleanCities(cities)
//             log.Printf("Fetched and cleaned %d cities for query '%s'", len(cleanedCities), q)
            
//             mu.Lock()
//             allCities = append(allCities, cleanedCities...)
//             mu.Unlock()
//         }(query)
//     }
    
//     // Wait for all goroutines
//     wg.Wait()
//     close(errChan)
    
//     // Check for errors
//     var errors []error
//     for err := range errChan {
//         errors = append(errors, err)
//     }
    
//     if len(errors) > 0 {
//         return fmt.Errorf("encountered %d errors while fetching cities: %v", len(errors), errors)
//     }
    
//     // Process cities if we have any
//     if len(allCities) == 0 {
//         return fmt.Errorf("no cities were fetched from any query")
//     }
    
//     log.Printf("Successfully fetched %d total cities", len(allCities))
//     return c.locationService.ProcessAndStoreCities(allCities)
// }


func (c *BookingController) fetchCities(query string) ([]models.City, error) {
    // Use rate-limited API client
    apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
    
    // Apply rate limiting
    if err := c.rateLimiter.Wait(context.Background()); err != nil {
        return nil, fmt.Errorf("rate limiter error: %v", err)
    }
    
    // Make API request
    body, err := c.apiClient.MakeRequest(context.Background(), apiURL)
    if err != nil {
        return nil, err
    }

    // Parse response
    var citiesResp models.CityResponse
    if err := json.Unmarshal(body, &citiesResp); err != nil {
        return nil, err
    }

    return citiesResp.Data, nil
}
  
func (c *BookingController) handleError(err error) {
    // Centralized error handling
    log.Printf("Controller error: %v", err)
    
    c.Data["json"] = map[string]string{
        "error":   err.Error(),
        "message": "An error occurred while processing your request",
    }
    c.Ctx.Output.SetStatus(400)
    c.ServeJSON()
}



