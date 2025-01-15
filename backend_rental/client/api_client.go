// backend_rental/client/api_client.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

type APIClient struct {
    httpClient *http.Client
    rapidAPIKey string
}

type PropertyDetails struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Address     string   `json:"address"`
    City        string   `json:"city"`
    Country     string   `json:"country"`
    Rating      float64  `json:"rating"`
    Price       float64  `json:"price"`
    Amenities   []string `json:"amenities"`
}

func NewAPIClient(rapidAPIKey string) *APIClient {
    return &APIClient{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        rapidAPIKey: rapidAPIKey,
    }
}

func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("X-RapidAPI-Key", c.rapidAPIKey)
    req.Header.Add("X-RapidAPI-Host", "booking-com18.p.rapidapi.com")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}

func (c *APIClient) FetchPropertyDetails(ctx context.Context, hotelID string) (*PropertyDetails, error) {
    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/properties/detail?hotel_id=%s", hotelID)
    
    data, err := c.MakeRequest(ctx, url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch property details: %v", err)
    }

    var details PropertyDetails
    if err := json.Unmarshal(data, &details); err != nil {
        return nil, fmt.Errorf("failed to parse property details: %v", err)
    }

    return &details, nil
}

// Interface for mocking in tests
type APIClientInterface interface {
    MakeRequest(ctx context.Context, url string) ([]byte, error)
    FetchPropertyDetails(ctx context.Context, hotelID string) (*PropertyDetails, error)
}
// package services

// import (
//     "backend_rental/models"
//     "backend_rental/client"  // Import the client package
//     "context"
//     "encoding/json"
//     "fmt"
// )

// type PropertyDetailsService interface {
//     GetPropertyDetails(hotelID string) (*models.PropertyDetails, error)
// }

// type propertyDetailsService struct {
//     apiClient client.APIClient  // Use client.APIClient
// }

// func NewPropertyDetailsService(apiClient client.APIClient) PropertyDetailsService {  // Use client.APIClient
//     return &propertyDetailsService{
//         apiClient: apiClient,
//     }
// }

// func (s *propertyDetailsService) GetPropertyDetails(hotelID string) (*models.PropertyDetails, error) {
//     resp, err := s.apiClient.FetchPropertyDetails(context.Background(), hotelID)
//     if err != nil {
//         return nil, fmt.Errorf("error fetching property details: %v", err)
//     }

//     var apiResponse struct {
//         Data struct {
//             HotelID              string `json:"hotel_id"`
//             HotelName            string `json:"hotel_name"`
//             AccommodationTypeName string `json:"accommodation_type_name"`
//             Rooms                map[string]struct {
//                 PrivateBathroomCount int `json:"private_bathroom_count"`
//             } `json:"rooms"`
//             FacilitiesBlock struct {
//                 Facilities []struct {
//                     Name string `json:"name"`
//                 } `json:"facilities"`
//             } `json:"facilities_block"`
//             BlockCount int `json:"block_count"`
//         } `json:"data"`
//     }

//     if err := json.Unmarshal(resp, &apiResponse); err != nil {
//         return nil, fmt.Errorf("error parsing JSON: %v", err)
//     }

//     details := &models.PropertyDetails{
//         HotelID:      apiResponse.Data.HotelID,
//         PropertyName: apiResponse.Data.HotelName,
//         Type:         apiResponse.Data.AccommodationTypeName,
//         Bedrooms:     apiResponse.Data.BlockCount,
//         Bathroom:     0,
//         Amenities:    make([]models.Facility, 0),
//     }

//     // Extract bathroom count from first room
//     for _, room := range apiResponse.Data.Rooms {
//         details.Bathroom = room.PrivateBathroomCount
//         break
//     }

//     // Extract facilities
//     for _, facility := range apiResponse.Data.FacilitiesBlock.Facilities {
//         if facility.Name != "" {
//             details.Amenities = append(details.Amenities, models.Facility{Name: facility.Name})
//         }
//     }

//     return details, nil
// }
