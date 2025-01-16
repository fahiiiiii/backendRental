package services

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"

    "backend_rental/models"
    "backend_rental/utils/apiclient"
)

type LocationProcessingService struct {
    apiClient    *apiclient.APIClient
    locationSvc  *LocationService
}

func NewLocationProcessingService(apiClient *apiclient.APIClient) *LocationProcessingService {
    return &LocationProcessingService{
        apiClient:   apiClient,
        locationSvc: &LocationService{},
    }
}

func (s *LocationProcessingService) ProcessLocationsFromQueries(queries []string) error {
    log.Printf("Starting to process %d queries", len(queries))

    var allCities []models.City
    var mu sync.Mutex
    var wg sync.WaitGroup
    errChan := make(chan error, len(queries))

    // Semaphore for concurrent requests
    semaphore := make(chan struct{}, 3)

    for _, query := range queries {
        wg.Add(1)
        semaphore <- struct{}{}

        go func(q string) {
            defer wg.Done()
            defer func() { <-semaphore }()

            cities, err := s.fetchCitiesForQuery(q)
            if err != nil {
                log.Printf("Error fetching cities for query '%s': %v", q, err)
                errChan <- err
                return
            }

            // Add logging for city data validation
            for _, city := range cities {
                log.Printf("Fetched city: ID='%s', Name='%s', Country='%s'",
                    city.CityID, city.CityName, city.Country)
            }

            mu.Lock()
            allCities = append(allCities, cities...)
            mu.Unlock()

            log.Printf("Fetched %d cities for query '%s'", len(cities), q)
        }(query)
    }

    wg.Wait()
    close(errChan)

    // Check for any errors that occurred during processing
    for err := range errChan {
        if err != nil {
            return fmt.Errorf("error during processing: %v", err)
        }
    }

    // Log raw data before processing
    log.Printf("Total cities fetched before processing: %d", len(allCities))

    return s.locationSvc.ProcessAndStoreCities(allCities)
}

func (s *LocationProcessingService) fetchCitiesForQuery(query string) ([]models.City, error) {
    apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
    
    body, err := s.apiClient.MakeRequest(context.Background(), apiURL)
    if err != nil {
        return nil, fmt.Errorf("API request failed: %v", err)
    }

    var citiesResp models.CityResponse
    if err := json.Unmarshal(body, &citiesResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }

    return citiesResp.Data, nil
}
// package services

// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "log"
//     "sync"
    
//     "backend_rental/models"
//     // "backend_rental/utils"
//     "backend_rental/utils/apiclient"
// )

// type LocationProcessingService struct {
//     apiClient      *apiclient.APIClient
//     locationSvc    *LocationService
// }

// func NewLocationProcessingService(apiClient *apiclient.APIClient) *LocationProcessingService {
//     return &LocationProcessingService{
//         apiClient:    apiClient,
//         locationSvc:  &LocationService{},
//     }
// }

// func (s *LocationProcessingService) ProcessLocationsFromQueries(queries []string) error {
//     log.Printf("Starting to process %d queries", len(queries))
    
//     var allCities []models.City
//     var mu sync.Mutex
//     var wg sync.WaitGroup
//     errChan := make(chan error, len(queries))
    
//     // Semaphore for concurrent requests
//     semaphore := make(chan struct{}, 3)
    
//     for _, query := range queries {
//         wg.Add(1)
//         semaphore <- struct{}{}
        
//         go func(q string) {
//             defer wg.Done()
//             defer func() { <-semaphore }()
            
//             cities, err := s.fetchCitiesForQuery(q)
//             if err != nil {
//                 log.Printf("Error fetching cities for query '%s': %v", q, err)
//                 errChan <- err
//                 return
//             }
            
//             // Add logging for city data validation
//             for _, city := range cities {
//                 log.Printf("Fetched city: ID='%s', Name='%s', Country='%s'",
//                     city.CityID, city.CityName, city.Country)
//             }
            
//             mu.Lock()
//             allCities = append(allCities, cities...)
//             mu.Unlock()
            
//             log.Printf("Fetched %d cities for query '%s'", len(cities), q)
//         }(query)
//     }
    
//     wg.Wait()
//     close(errChan)
    
//     // Log raw data before processing
//     log.Printf("Total cities fetched before processing: %d", len(allCities))
    
//     return s.locationSvc.ProcessAndStoreCities(allCities)
// }
// func (s *LocationProcessingService) fetchCitiesForQuery(query string) ([]models.City, error) {
//     rateLimiter := ratelimiter.GetInstance()
    
//     // Wait for rate limit
//     if err := rateLimiter.Wait(context.Background()); err != nil {
//         return nil, fmt.Errorf("rate limiter error: %v", err)
//     }
    
//     apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
//     body, err := s.apiClient.MakeRequest(context.Background(), apiURL)
//     // apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
    
//     // body, err := s.apiClient.MakeRequest(context.Background(), apiURL)
//     if err != nil {
//         return nil, fmt.Errorf("API request failed: %v", err)
//     }
    
//     var citiesResp models.CityResponse
//     if err := json.Unmarshal(body, &citiesResp); err != nil {
//         return nil, fmt.Errorf("failed to parse response: %v", err)
//     }
    
//     return citiesResp.Data, nil
// }
