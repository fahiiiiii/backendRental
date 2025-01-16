package handlers

import (
    "backend_rental/models"
    "backend_rental/utils"
    "log"
    "net/http"
    "encoding/json"
    "time"
    "io/ioutil"
    "github.com/go-resty/resty/v2"
)

// HandleCityQuery handles a city query and returns locations
func HandleCityQuery(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    if query == "" {
        http.Error(w, "Query parameter is missing", http.StatusBadRequest)
        return
    }
    
    cities, err := fetchCitiesFromAPI(query)
    if err != nil {
        log.Printf("Error fetching cities: %v", err)
        http.Error(w, "Error fetching cities", http.StatusInternalServerError)
        return
    }

    cleanedCities := utils.CleanCities(cities)
    locations, err := utils.ConvertCitiesToLocations(cleanedCities)
    if err != nil {
        log.Printf("Error converting cities: %v", err)
        http.Error(w, "Error converting cities", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "locations": locations,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// fetchCitiesFromAPI fetches city data from the booking API
func fetchCitiesFromAPI(query string) ([]models.City, error) {
    client := resty.New()
    url := "https://booking-com18.p.rapidapi.com/stays/auto-complete?query=" + query

    resp, err := client.R().
        SetQueryParams(map[string]string{
            "query": query,
        }).
        SetHeaders(map[string]string{
            "X-RapidAPI-Host": "booking-com18.p.rapidapi.com",
            "X-RapidAPI-Key": "YOUR-RAPIDAPI-KEY",
        }).
        Get(url)

    if err != nil {
        return nil, err
    }

    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("received non-200 status code: %v", resp.StatusCode())
    }

    var apiResponse struct {
        Result []struct {
            Name    string `json:"name"`
            Country string `json:"country"`
        } `json:"result"`
    }

    err = json.Unmarshal(resp.Body(), &apiResponse)
    if err != nil {
        return nil, err
    }

    var cities []models.City
    for _, item := range apiResponse.Result {
        cities = append(cities, models.City{
            Name:    item.Name,
            Country: item.Country,
        })
    }

    return cities, nil
}
