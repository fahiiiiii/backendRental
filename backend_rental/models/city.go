// // models/city.go
package models

type City struct {
    CityID      string  `json:"city_id"`
    CityName    string  `json:"city_name"`
    Country     string  `json:"country"`
}

type CityResponse struct {
    Data []City `json:"data"`
}
// // models/city.go
// package models

// type City struct {
//     CityID      string  `json:"city_id"`
//     CityName    string  `json:"city_name"`
//     Country     string  `json:"country"`
// }

// type CityResponse struct {
//     Data []City `json:"data"`
// }
// type CityDetails struct {
//     ID          string   `json:"id"`
//     CityName    string   `json:"city_name"`
//     Country     string   `json:"country"`
//     Properties  []string `json:"properties"`
// }

// type BookingSummary struct {
//     Countries      map[string]bool            `json:"countries"`
//     Cities         map[string]bool            `json:"cities"`
//     CountryCities  map[string][]string        `json:"country_cities"`
//     CityProperties map[string][]string        `json:"city_properties"`
// }