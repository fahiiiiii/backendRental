package services

import (
    "backend_rental/client"
    "backend_rental/models"
    "context"
    "encoding/json"
    "fmt"
)

type PropertyDetailsService interface {
    GetPropertyDetails(hotelID string) (*models.PropertyDetails, error)
}

type propertyDetailsService struct {
    apiClient client.APIClient  // Changed from util.APIClient to client.APIClient
}

func NewPropertyDetailsService(apiClient client.APIClient) PropertyDetailsService {
    return &propertyDetailsService{
        apiClient: apiClient,
    }
}

func (s *propertyDetailsService) GetPropertyDetails(hotelID string) (*models.PropertyDetails, error) {
    // Fetch property details directly
    apiPropertyDetails, err := s.apiClient.FetchPropertyDetails(context.Background(), hotelID)
    if err != nil {
        return nil, fmt.Errorf("error fetching property details: %v", err)
    }

    // Create details with initial values from propertyDetails
    details := &models.PropertyDetails{
        HotelID:      apiPropertyDetails.ID,
        PropertyName: apiPropertyDetails.Name,
        Type:         apiPropertyDetails.AccommodationTypeName, // You'll need to set this appropriately
        Bedrooms:     apiPropertyDetails.BlockCount,  // You'll need to set this appropriately
        Bathroom:     0,  // Initialize to 0
        Amenities:    make([]models.Facility, 0), // Initialize empty slice
    }

    // Convert original amenities to Facility type
    for _, amenity := range apiPropertyDetails.Amenities {
        details.Amenities = append(details.Amenities, models.Facility{Name: amenity})
    }

    // If you need to unmarshal additional data, use the original propertyDetails
    var apiResponse struct {
        Data struct {
            HotelID              string `json:"hotel_id"`
            HotelName            string `json:"hotel_name"`
            AccommodationTypeName string `json:"accommodation_type_name"`
            Rooms                map[string]struct {
                PrivateBathroomCount int `json:"private_bathroom_count"`
            } `json:"rooms"`
            FacilitiesBlock struct {
                Facilities []struct {
                    Name string `json:"name"`
                } `json:"facilities"`
            } `json:"facilities_block"`
            BlockCount int `json:"block_count"`
        } `json:"data"`
    }

    // Unmarshal using the byte data
    if err := json.Unmarshal(apiPropertyDetails, &apiResponse); err != nil {
        return nil, fmt.Errorf("error parsing JSON: %v", err)
    }

    // Update details with additional information from apiResponse
    details.Type = apiResponse.Data.AccommodationTypeName
    details.Bedrooms = apiResponse.Data.BlockCount

    // Extract bathroom count from first room
    for _, room := range apiResponse.Data.Rooms {
        details.Bathroom = room.PrivateBathroomCount
        break
    }

    // Add additional facilities from apiResponse
    for _, facility := range apiResponse.Data.FacilitiesBlock.Facilities {
        if facility.Name != "" {
            details.Amenities = append(details.Amenities, models.Facility{Name: facility.Name})
        }
    }

    return details, nil
}
// func (s *propertyDetailsService) GetPropertyDetails(hotelID string) (*models.PropertyDetails, error) {
//     // Fetch property details directly
//     propertyDetails, err := s.apiClient.FetchPropertyDetails(context.Background(), hotelID)
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
// 	details := &models.PropertyDetails{
//         HotelID:      propertyDetails.ID,
//         PropertyName: propertyDetails.Name,
//         Type:         "", // Set this based on your requirements
//         Bedrooms:     0,  // Set this based on your requirements
//         Bathroom:     propertyDetails.BlockCount,  // Set this based on your requirements
//         Amenities:    propertyDetails.Amenities,
//     }

//     // details := &models.PropertyDetails{
//     //     HotelID:      apiResponse.Data.HotelID,
//     //     PropertyName: apiResponse.Data.HotelName,
//     //     Type:         apiResponse.Data.AccommodationTypeName,
//     //     Bedrooms:     apiResponse.Data.BlockCount,
//     //     Bathroom:     0,
//     //     Amenities:    make([]models.Facility, 0),
//     // }

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
