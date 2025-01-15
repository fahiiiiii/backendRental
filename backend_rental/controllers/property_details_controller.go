package controllers

import (
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailsController struct {
    beego.Controller
    propertyService services.PropertyService
}

// func (c *PropertyDetailsController) Get() {
//     hotelID := c.GetString("hotel_id")
//     if hotelID == "" {
//         c.Ctx.Output.SetStatus(400)
//         c.Data["json"] = map[string]string{"error": "hotel_id is required"}
//         c.ServeJSON()
//         return
//     }

//     details, err := c.propertyService.GetPropertyDetails(hotelID)
//     if err != nil {
//         c.Ctx.Output.SetStatus(500)
//         c.Data["json"] = map[string]string{"error": err.Error()}
//         c.ServeJSON()
//         return
//     }

//     c.Data["json"] = details
//     c.ServeJSON()
// }

func (s *propertyDetailsService) GetPropertyDetails(hotelID string) (*models.PropertyDetails, error) {
    // Fetch property details directly
    apiPropertyDetails, err := s.apiClient.FetchPropertyDetails(context.Background(), hotelID)
    if err != nil {
        return nil, fmt.Errorf("error fetching property details: %v", err)
    }

    // Create details with initial values from apiPropertyDetails
    details := &models.PropertyDetails{
        HotelID:      apiPropertyDetails.ID,
        PropertyName: apiPropertyDetails.Name,
        Type:          apiPropertyDetails.AccommodationTypeName// You'll need to set this appropriately
        Bedrooms:     apiPropertyDetails.BlockCount,  // You'll need to set this appropriately
        Bathroom:     0,  // Initialize to 0
        Amenities:    make([]models.Facility, 0), // Initialize empty slice
    }

    // Convert original amenities to Facility type
    for _, amenity := range apiPropertyDetails.Amenities {
        details.Amenities = append(details.Amenities, models.Facility{Name: amenity})
    }

    // Prepare to unmarshal additional data
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

    // Convert the property details to JSON bytes
    jsonData, err := json.Marshal(apiPropertyDetails)
    if err != nil {
        return nil, fmt.Errorf("error marshaling property details: %v", err)
    }

    // Unmarshal using the JSON bytes
    if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
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