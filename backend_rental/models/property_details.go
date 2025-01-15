package models

type PropertyDetails struct {
    HotelID      string     `json:"hotel_id"`
    PropertyName string     `json:"property_name"`
    Type         string     `json:"type"`
    Bedrooms     int        `json:"bedrooms"`
    Bathroom     int        `json:"bathroom"`
    Amenities    []Facility `json:"amenities"`
}

type Facility struct {
    Name string `json:"name"`
}
