package models

type CityKey struct {
    Name    string
    Country string
}

type Property struct {
    DestID      string `json:"dest_id"`
    Name        string `json:"name"`
    CityID      string  `json:"city_id"`
    Country     string `json:"country"`
    // Add other fields as needed from the API response
}