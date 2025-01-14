
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "sync"

    "backend_rental/models"
    "backend_rental/utils"
    "backend_rental/utils/apiclient"
)

type LocationProcessingService struct {
    apiClient *apiclient.APIClient
}

func NewLocationProcessingService(apiClient *apiclient.APIClient) *LocationProcessingService {
    return &LocationProcessingService{
        apiClient: apiClient,
    }
}

func (s *LocationProcessingService) ProcessLocationsFromQueries(queries []string) error {
    var allLocations []models.Location
    var mu sync.Mutex
    var wg sync.WaitGroup

    semaphore := make(chan struct{}, 5)

    for _, query := range queries {
        wg.Add(1)
        semaphore <- struct{}{}

        go func(q string) {
            defer wg.Done()
            defer func() { <-semaphore }()

            cities, err := s.fetchCitiesForQuery(q)
            if err != nil {
                log.Printf("Error fetching cities for query '%s': %v", q, err)
                return
            }

            locationBatch := s.convertCitiesToLocations(cities)

            mu.Lock()
            allLocations = append(allLocations, locationBatch...)
            mu.Unlock()
        }(query)
    }

    wg.Wait()

    // Remove duplicates
    uniqueLocations := utils.RemoveDuplicateLocations(allLocations)

    // Bulk store locations
    locationService := &LocationService{}
    return locationService.BulkCreateLocations(uniqueLocations)
}

func (s *LocationProcessingService) fetchCitiesForQuery(query string) ([]models.City, error) {
    apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)

    body, err := s.apiClient.MakeRequest(context.Background(), apiURL)
    if err != nil {
        return nil, err
    }

    var citiesResp models.CityResponse
    if err := json.Unmarshal(body, &citiesResp); err != nil {
        return nil, err
    }

    return citiesResp.Data, nil
}

func (s *LocationProcessingService) convertCitiesToLocations(cities []models.City) []models.Location {
    locations := make([]models.Location, 0, len(cities))

    for _, city := range cities {
        if city.CityName == "" || city.CityID == "" {
            continue
        }

        location := models.Location{
            ID:          city.CityID,
            CityName:    city.CityName,
            Country:     city.Country,
            CountryCode: generateCountryCode(city.Country),
            Latitude:    0.0, // Add logic to fetch latitude if needed
            Longitude:   0.0, // Add logic to fetch longitude if needed
        }

        locations = append(locations, location)
    }

    return locations
}

// Helper function to generate country code
func generateCountryCode(country string) string {
    // Simple implementation - you might want to use a more comprehensive mapping
    if len(country) >= 2 {
        return strings.ToUpper(country[:2])
    }
    return ""
}