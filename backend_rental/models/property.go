// models/property.go
package models

type Property struct {
    DestID      string `json:"dest_id"`
    Name        string `json:"name"`
    Country     string `json:"country"`
    City        string `json:"city"`
    Type        string `json:"type"`
    Rating      float64 `json:"rating,omitempty"`
}

type CityKey struct {
    Name    string `json:"name"`
    Country string `json:"country"`
}

type PropertyResponse struct {
    Status  string     `json:"status"`
    Message string     `json:"message,omitempty"`
    Data    []Property `json:"data"`
}

type PropertySummary struct {
    TotalProperties int                    `json:"total_properties"`
    CitiesCount     int                    `json:"cities_count"`
    CityProperties  map[string][]string    `json:"city_properties"`
}