// backend_rental/models/city.go
package models

type City struct {
    CityID      string  `json:"city_id"`
    CityName    string  `json:"city_name"`
    Country     string  `json:"country"`
}

type CityResponse struct {
    Data []City `json:"data"`
}

// // // models/city.go
// package models

// type City struct {
//     CityID      string  `json:"city_id"`
//     CityName    string  `json:"city_name"`
//     Country     string  `json:"country"`
// }

// type CityResponse struct {
//     Data []City `json:"data"`
// }