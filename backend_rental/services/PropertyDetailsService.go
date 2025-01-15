package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"backend_rental/models"
	"backend_rental/utils/ratelimiter"
	"github.com/beego/beego/v2/client/orm"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailsService struct {
	rateLimiter *ratelimiter.APIRateLimiter
	apiKey      string
}

func NewPropertyDetailsService() *PropertyDetailsService {
	apiKey, _ := beego.AppConfig.String("rapidapi.key") // Load API key from config
	return &PropertyDetailsService{
		rateLimiter: ratelimiter.GetInstance(),
		apiKey:      apiKey,
	}
}

func (s *PropertyDetailsService) FetchAndProcessPropertyDetails(propertyIDs []string) (map[string]*models.PropertyDetails, error) {
	results := make(chan *models.PropertyDetails, len(propertyIDs))
	var wg sync.WaitGroup

	// Semaphore for rate-limiting
	semaphore := make(chan struct{}, 1)

	for _, id := range propertyIDs {
		wg.Add(1)
		go func(propertyID string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Wait for rate limiter
			if err := s.rateLimiter.Wait(context.Background()); err != nil {
				fmt.Printf("Rate limiter error for property %s: %v\n", propertyID, err)
				return
			}

			// Fetch property details from API
			details, err := s.fetchPropertyDetailsFromAPI(propertyID)
			if err != nil {
				fmt.Printf("Error fetching property %s: %v\n", propertyID, err)
				return
			}

			results <- details
		}(id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	propertyDetails := make(map[string]*models.PropertyDetails)
	for detail := range results {
		if detail != nil {
			propertyDetails[detail.PropertyID] = detail

			// Save to database
			if err := s.savePropertyDetailsToDatabase(detail); err != nil {
				fmt.Printf("Error saving property %s to database: %v\n", detail.PropertyID, err)
			}
		}
	}

	return propertyDetails, nil
}

func (s *PropertyDetailsService) fetchPropertyDetailsFromAPI(propertyID string) (*models.PropertyDetails, error) {
    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%s&checkinDate=2025-01-09&checkoutDate=2025-01-23&units=metric", propertyID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", s.apiKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var apiResponse struct {
        Data struct {
            PropertyID            string `json:"hotel_id"`
            PropertyName          string `json:"hotel_name"`
            AccommodationTypeName string `json:"accommodation_type_name"`
            BlockCount            int    `json:"block_count"`
            Rooms                 map[string]struct {
                PrivateBathroomCount int `json:"private_bathroom_count"`
            } `json:"rooms"`
            FacilitiesBlock struct {
                Facilities []struct {
                    Name string `json:"name"`
                } `json:"facilities"`
            } `json:"facilities_block"`
            Location struct {
                CityName string `json:"city_name"`
            } `json:"location"`
        } `json:"data"`
    }

    if err := json.Unmarshal(body, &apiResponse); err != nil {
        return nil, err
    }

    // Look up CityID from the location
    cityID, err := getCityIDFromCityName(apiResponse.Data.Location.CityName)
    if err != nil {
        return nil, fmt.Errorf("could not find CityID for city: %s", apiResponse.Data.Location.CityName)
    }

    details := &models.PropertyDetails{
        PropertyID:   apiResponse.Data.PropertyID,
        PropertyName: apiResponse.Data.PropertyName,
        Type:         apiResponse.Data.AccommodationTypeName,
        Bedrooms:     apiResponse.Data.BlockCount,
        Amenities:    []models.Facility{},
        CityID:       &cityID, // Use pointer to int here
    }

    // Extract bathroom count
    for _, room := range apiResponse.Data.Rooms {
        details.Bathroom = room.PrivateBathroomCount
        break
    }

    // Extract amenities
    for _, facility := range apiResponse.Data.FacilitiesBlock.Facilities {
        details.Amenities = append(details.Amenities, models.Facility{Name: facility.Name})
    }

    return details, nil
}

// Define getCityIDFromCityName (Dummy Implementation)
func getCityIDFromCityName(cityName string) (int, error) {
    // Dummy implementation for now. You can replace this with actual logic to look up CityID from your data source
    cityMap := map[string]int{
        "New York": 1,
        "Los Angeles": 2,
        "Chicago": 3,
    }
    
    cityID, exists := cityMap[cityName]
    if !exists {
        return 0, fmt.Errorf("City not found: %s", cityName)
    }

    return cityID, nil
}

func (s *PropertyDetailsService) savePropertyDetailsToDatabase(details *models.PropertyDetails) error {
	// Initialize ORM and start a new transaction
	o := orm.NewOrm()

	// Prepare the RentalProperty model
	property := models.RentalProperty{
		CityID:      details.CityID, // Use CityID directly (now a *int)
		DestID:      func() *int { val := 1; return &val }(), // Set DestID as a pointer
		Name:        details.PropertyName,
		Type:        details.Type,
		Bedrooms:    details.Bedrooms,
		Bathroom:    details.Bathroom,
		Amenities:   fmt.Sprintf("%v", details.Amenities), // Store amenities as a JSON string or CSV
	}

	// Insert the data into the database
	_, err := o.Insert(&property)
	if err != nil {
		return fmt.Errorf("error inserting property details into the database: %v", err)
	}

	return nil
}


// package services

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"sync"
//     "time"
// 	"backend_rental/models"
// 	"backend_rental/utils/ratelimiter"
// 	beego "github.com/beego/beego/v2/server/web"

// )

// type PropertyDetailsService struct {
// 	rateLimiter *ratelimiter.APIRateLimiter
// 	apiKey      string
// }

// func NewPropertyDetailsService() *PropertyDetailsService {
// 	apiKey, _ := beego.AppConfig.String("rapidapi.key") // Load API key from config
// 	return &PropertyDetailsService{
// 		rateLimiter: ratelimiter.GetInstance(),
// 		apiKey:      apiKey,
// 	}
// }

// // FetchAndProcessPropertyDetails fetches and processes property details from an API.
// func (s *PropertyDetailsService) FetchAndProcessPropertyDetails(propertyIDs []string) (map[string]*models.PropertyDetails, error) {
// 	results := make(chan *models.PropertyDetails, len(propertyIDs))
// 	var wg sync.WaitGroup

// 	// Semaphore for rate-limiting
// 	semaphore := make(chan struct{}, 1)

// 	for _, id := range propertyIDs {
// 		wg.Add(1)
// 		go func(propertyID string) {
// 			defer wg.Done()
// 			semaphore <- struct{}{}
// 			defer func() { <-semaphore }()

// 			// Wait for rate limiter
// 			if err := s.rateLimiter.Wait(context.Background()); err != nil {
// 				fmt.Printf("Rate limiter error for property %s: %v\n", propertyID, err)
// 				return
// 			}

// 			// Fetch property details from API
// 			details, err := s.fetchPropertyDetailsFromAPI(propertyID)
// 			if err != nil {
// 				fmt.Printf("Error fetching property %s: %v\n", propertyID, err)
// 				return
// 			}

// 			results <- details
// 		}(id)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	propertyDetails := make(map[string]*models.PropertyDetails)
// 	for detail := range results {
// 		if detail != nil {
// 			propertyDetails[detail.PropertyID] = detail
// 		}
// 	}

// 	// Save results to file
// 	if err := s.savePropertyDetailsToFile(propertyDetails); err != nil {
// 		return nil, err
// 	}

// 	return propertyDetails, nil
// }

// func (s *PropertyDetailsService) fetchPropertyDetailsFromAPI(propertyID string) (*models.PropertyDetails, error) {
// 	url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%s&checkinDate=2025-01-09&checkoutDate=2025-01-23&units=metric", propertyID)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
// 	req.Header.Add("x-rapidapi-key", s.apiKey)

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var apiResponse struct {
// 		Data struct {
// 			PropertyID            string `json:"hotel_id"`
// 			PropertyName          string `json:"hotel_name"`
// 			AccommodationTypeName string `json:"accommodation_type_name"`
// 			BlockCount            int    `json:"block_count"`
// 			Rooms                 map[string]struct {
// 				PrivateBathroomCount int `json:"private_bathroom_count"`
// 			} `json:"rooms"`
// 			FacilitiesBlock struct {
// 				Facilities []struct {
// 					Name string `json:"name"`
// 				} `json:"facilities"`
// 			} `json:"facilities_block"`
// 		} `json:"data"`
// 	}

// 	if err := json.Unmarshal(body, &apiResponse); err != nil {
// 		return nil, err
// 	}

// 	details := &models.PropertyDetails{
// 		PropertyID:   apiResponse.Data.PropertyID,
// 		PropertyName: apiResponse.Data.PropertyName,
// 		Type:         apiResponse.Data.AccommodationTypeName,
// 		Bedrooms:     apiResponse.Data.BlockCount,
// 		Amenities:    []models.Facility{},
// 	}

// 	// Extract bathroom count
// 	for _, room := range apiResponse.Data.Rooms {
// 		details.Bathroom = room.PrivateBathroomCount
// 		break
// 	}

// 	// Extract amenities
// 	for _, facility := range apiResponse.Data.FacilitiesBlock.Facilities {
// 		details.Amenities = append(details.Amenities, models.Facility{Name: facility.Name})
// 	}

// 	return details, nil
// }

// func (s *PropertyDetailsService) savePropertyDetailsToFile(propertyDetails map[string]*models.PropertyDetails) error {
// 	outputData, err := json.MarshalIndent(propertyDetails, "", "    ")
// 	if err != nil {
// 		return fmt.Errorf("error marshaling property details: %v", err)
// 	}

// 	return os.WriteFile("property_details.json", outputData, 0644)
// }
