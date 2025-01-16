// backend_rental/utils/location_utils.go
package utils

import (
    "strings"
    "backend_rental/models"
    "regexp"
    "log"
)

// FilterAndCleanCities filters and cleans city data
func FilterAndCleanCities(cities []models.City) []models.City {
    log.Printf("Starting cleaning process with %d cities", len(cities))
    var cleanedCities []models.City
    seen := make(map[string]bool)
    
    for i, city := range cities {
        // Log each city being processed
        log.Printf("Processing city %d: Name='%s', Country='%s'", i, city.CityName, city.Country)
        
        // Skip if name or country is empty
        if city.CityName == "" || city.Country == "" {
            log.Printf("Skipping city %d: empty name or country", i)
            continue
        }
        
        // Clean the data
        cleanCity := models.City{
            CityID:   generateCityID(city.CityName, city.Country),
            CityName: strings.TrimSpace(city.CityName),
            Country:  strings.TrimSpace(city.Country),
        }
        
        // Create a unique key for deduplication
        key := strings.ToLower(cleanCity.CityName + "|" + cleanCity.Country)
        
        if seen[key] {
            log.Printf("Skipping duplicate city: %s, %s", cleanCity.CityName, cleanCity.Country)
            continue
        }
        
        seen[key] = true
        cleanedCities = append(cleanedCities, cleanCity)
        log.Printf("Added cleaned city: %s, %s, ID: %s", cleanCity.CityName, cleanCity.Country, cleanCity.CityID)
    }
    
    log.Printf("Finished cleaning process. %d cities retained", len(cleanedCities))
    return cleanedCities
}
// func FilterAndCleanCities(cities []models.City) []models.City {
//     var cleanedCities []models.City
//     seen := make(map[string]bool)
    
//     for _, city := range cities {
//         // Skip if required fields are empty
//         if city.CityID == "" || city.CityName == "" || city.Country == "" {
//             continue
//         }
        
//         // Clean the data
//         cleanCity := models.City{
//             CityID:   strings.TrimSpace(city.CityID),
//             CityName: strings.TrimSpace(city.CityName),
//             Country:  strings.TrimSpace(city.Country),
//         }
        
//         // Create a unique key for deduplication
//         key := cleanCity.CityID + "|" + cleanCity.CityName + "|" + cleanCity.Country
        
//         if seen[key] {
//             continue
//         }
        
//         seen[key] = true
//         cleanedCities = append(cleanedCities, cleanCity)
//     }
    
//     return cleanedCities
// }
func ConvertCitiesToLocations(cities []models.City) []models.Location {
    log.Printf("Starting conversion of %d cities to locations", len(cities))
    locations := make([]models.Location, 0, len(cities))
    
    for i, city := range cities {
        location := models.Location{
            ID:          city.CityID,
            CityName:    city.CityName,
            Country:     city.Country,
            CountryCode: ExtractCountryCode(city.Country),
            Latitude:    0.0,  // Default values since API doesn't provide coordinates
            Longitude:   0.0,
        }
        locations = append(locations, location)
        log.Printf("Converted city %d to location: %s, %s", i, location.CityName, location.Country)
    }
    
    log.Printf("Finished conversion. Created %d locations", len(locations))
    return locations
}
// func ConvertCitiesToLocations(cities []models.City) []models.Location {
//     locations := make([]models.Location, 0, len(cities))
    
//     for _, city := range cities {
//         location := models.Location{
//             ID:          GenerateUniqueLocationID(city),
//             CityName:    city.CityName,
//             Country:     city.Country,
//             CountryCode: ExtractCountryCode(city.Country),
//             Latitude:    0.0,
//             Longitude:   0.0,
//         }
//         locations = append(locations, location)
//     }
    
//     return locations
// }
func generateCityID(cityName, country string) string {
    // Create a slug-style ID by combining city and country
    combined := strings.ToLower(cityName + "-" + country)
    // Replace spaces and special characters with dashes
    reg := regexp.MustCompile(`[^a-z0-9]+`)
    slug := reg.ReplaceAllString(combined, "-")
    // Remove leading and trailing dashes
    slug = strings.Trim(slug, "-")
    // Remove consecutive dashes
    reg = regexp.MustCompile(`-+`)
    slug = reg.ReplaceAllString(slug, "-")
    return slug
}
// GenerateUniqueLocationID creates a unique ID for a location
func GenerateUniqueLocationID(city models.City) string {
    return strings.ToLower(strings.ReplaceAll(city.CityID, " ", "-"))
}

// ExtractCountryCode extracts a 2-letter country code from country name
func ExtractCountryCode(country string) string {
    // TODO: Consider implementing a proper country code mapping
    country = strings.TrimSpace(country)
    if len(country) >= 2 {
        return strings.ToUpper(country[:2])
    }
    return "XX"
}

// GenerateLocationQueries generates queries for location search
func GenerateLocationQueries() []string {
    // letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
        // "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
    letters := []string{ "C"}
    
    // words := []string{"new", "old", "north", "south", "east", "west", "central"}
    words := []string{"new"}
    
    // cities := []string{"New York", "London", "Tokyo", "Paris", "Dubai"}
    cities := []string{"Dubai"}
    
    queries := make([]string, 0, len(letters)+len(words)+len(cities))
    queries = append(queries, letters...)
    queries = append(queries, words...)
    queries = append(queries, cities...)
    
    return queries
}

// ConvertToCountryMap converts country-cities map to country map
func ConvertToCountryMap(countryCities map[string][]string) map[string]bool {
    countries := make(map[string]bool)
    for country := range countryCities {
        countries[country] = true
    }
    return countries
}

// ConvertToCityMap converts country-cities map to city map
func ConvertToCityMap(countryCities map[string][]string) map[string]bool {
    cities := make(map[string]bool)
    for _, cityList := range countryCities {
        for _, city := range cityList {
            cities[city] = true
        }
    }
    return cities
}
func RemoveDuplicateLocations(locations []models.Location) []models.Location {
    log.Printf("Starting duplicate removal from %d locations", len(locations))
    seen := make(map[string]bool)
    result := make([]models.Location, 0)

    for _, loc := range locations {
        // Create a unique key combining city name and country
        key := strings.ToLower(loc.CityName + "|" + loc.Country)
        
        if !seen[key] {
            seen[key] = true
            result = append(result, loc)
        }
    }
    
    log.Printf("Finished removing duplicates. %d locations remaining", len(result))
    return result
}
// // RemoveDuplicateLocations removes duplicate locations based on ID
// func RemoveDuplicateLocations(locations []models.Location) []models.Location {
//     seen := make(map[string]bool)
//     result := make([]models.Location, 0)

//     for _, location := range locations {
//         // Skip if we've seen this ID before
//         if seen[location.ID] {
//             continue
//         }
        
//         seen[location.ID] = true
//         result = append(result, location)
//     }

//     return result
// }