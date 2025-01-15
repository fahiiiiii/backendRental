package models

type PropertyDetails struct {
    PropertyID   string     `json:"property_id"`
    PropertyName string     `json:"property_name"`
    Type         string     `json:"type"`
    Bedrooms     int        `json:"bedrooms"`
    Bathroom     int        `json:"bathroom"`
    Amenities    []Facility `json:"amenities"`
    CityID       *int       `json:"city_id"` // Add the CityID field as a pointer
}

type Facility struct {
    Name string `json:"name"`
}
