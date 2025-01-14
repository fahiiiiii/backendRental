// models/booking_summary.go
package models

type BookingSummary struct {
    Countries      map[string]bool            `json:"countries"`
    Cities         map[string]bool            `json:"cities"`
    CountryCities  map[string][]string        `json:"country_cities"`
    CityProperties map[string][]string        `json:"city_properties"`
}