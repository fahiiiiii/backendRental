package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"backend_rental/models"
	beego "github.com/beego/beego/v2/server/web"
)

// Define BookingResult within the same file
type BookingResult struct {
	Countries      map[string]bool                `json:"countries"`
	Cities         map[string]*models.CityDetails `json:"cities"`
	CountryCities  map[string][]string            `json:"country_cities"`
	CityProperties map[string][]string            `json:"city_properties"`
}

type BookingController struct {
	beego.Controller
}
func (c *BookingController) SearchCities() {
	// Get search parameters
	query := c.GetString("query")
	country := c.GetString("country")
	
	// Validate input
	if query == "" && country == "" {
		c.Data["json"] = map[string]string{"error": "Please provide a search query or country"}
		c.ServeJSON()
		return
	}

	// Perform search
	cities, err := searchCitiesByCriteria(query, country)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return results
	c.Data["json"] = cities
	c.ServeJSON()
}

// Helper function for advanced city search
func searchCitiesByCriteria(query, country string) ([]models.CityDetails, error) {
	var results []models.CityDetails

	// If no specific query, use a default search
	if query == "" {
		query = "a"
	}

	// Fetch cities based on query
	cities, err := fetchCities(query)
	if err != nil {
		return nil, err
	}

	// Filter cities based on criteria
	for _, city := range cities {
		// Check country filter if provided
		if country != "" && !strings.EqualFold(city.Country, country) {
			continue
		}

		// Convert to CityDetails
		cityDetails := models.CityDetails{
			ID:      city.ID,
			Name:    city.CityName,
			Country: city.Country,
		}

		results = append(results, cityDetails)
	}

	return results, nil
}

// Get handles the main retrieval of booking data
func (c *BookingController) Get() {
	// Initialize data storage
	uniqueCountries := make(map[string]bool)
	uniqueCities := make(map[string]*models.CityDetails)
	countryCities := make(map[string][]string)
	cityProperties := make(map[string][]string)

	// Get query parameter
	query := c.GetString("query")

	// Generate queries
	var queries []string
	if query != "" {
		queries = []string{query}
	} else {
		queries = generateQueries()
	}

	// Concurrency control
	var wg sync.WaitGroup
	var mutex sync.Mutex
	results := make(chan []models.City, len(queries))

	// Process queries concurrently
	for _, q := range queries {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			cities, err := fetchCities(query)
			if err != nil {
				log.Printf("Error fetching cities for query '%s': %v", query, err)
				return
			}
			results <- cities
		}(q)
	}

	// Close results channel when done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for cities := range results {
		for _, city := range cities {
			mutex.Lock()
			
			country := strings.TrimSpace(strings.ToUpper(city.Country))
			cityName := strings.TrimSpace(strings.ToUpper(city.CityName))

			if country != "" {
				uniqueCountries[country] = true
			}

			if cityName != "" {
				// Create or update city details
				cityDetails := &models.CityDetails{
					Name:    cityName,
					Country: country,
					ID:      city.ID,  // Include city ID
				}

				uniqueCities[cityName] = cityDetails

				// Add to country-cities mapping
				if _, exists := countryCities[country]; !exists {
					countryCities[country] = []string{}
				}
				
				// Avoid duplicates
				cityExists := false
				for _, existingCity := range countryCities[country] {
					if existingCity == cityName {
						cityExists = true
						break
					}
				}
				
				if !cityExists {
					countryCities[country] = append(countryCities[country], cityName)
				}
			}
			
			mutex.Unlock()
		}
	}

	// Prepare response
	result := BookingResult{
		Countries:      uniqueCountries,
		Cities:         uniqueCities,
		CountryCities:  countryCities,
		CityProperties: cityProperties,
	}

	// Set response
	c.Data["json"] = result
	c.ServeJSON()
}

// Helper function to generate queries
func generateQueries() []string {
	queries := []string{}
	
	// Alphabet queries
	for char := 'A'; char <= 'Z'; char++ {
		queries = append(queries, string(char))
	}

	// Common prefixes
	prefixes := []string{
		"a", "the", "new", "old", "big", "small", 
		"north", "south", "east", "west", "central",
	}

	queries = append(queries, prefixes...)
	return queries
}

// Helper function to fetch cities from API
func fetchCities(query string) ([]models.City, error) {
	// Get API key with proper error handling
	apiKey, err := beego.AppConfig.String("booking_api_key")
	if err != nil {
		return nil, fmt.Errorf("error getting API key: %v", err)
	}

	apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
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

	var citiesResp struct {
		Data []models.City `json:"data"`
	}
	err = json.Unmarshal(body, &citiesResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return citiesResp.Data, nil
}