package utils

import (
    "backend_rental/models"
    "fmt"
)

// GenerateUniqueLocationID generates a unique ID for a location
func GenerateUniqueLocationID(city models.City) string {
    return fmt.Sprintf("%s-%s", city.CityID, city.Country)
}

// ExtractCountryCode extracts a country code from the country name
func ExtractCountryCode(country string) string {
    countryCodes := map[string]string{
        "United States": "US",
        "Canada":        "CA",
        "United Kingdom": "GB",
    }

    return countryCodes[country]
}

// RemoveDuplicateLocations removes duplicate locations from a slice
func RemoveDuplicateLocations(locations []models.Location) []models.Location {
    seen := make(map[string]bool)
    uniqueLocations := []models.Location{}

    for _, location := range locations {
        if _, found := seen[location.ID]; !found {
            seen[location.ID] = true
            uniqueLocations = append(uniqueLocations, location)
        }
    }

    return uniqueLocations
}

// FilterAndCleanCities filters out cities that are missing critical information
func FilterAndCleanCities(cities []models.City) []models.City {
    var cleanedCities []models.City

    for _, city := range cities {
        if city.CityID == "" || city.CityName == "" || city.Country == "" {
            continue
        }

        cleanedCities = append(cleanedCities, city)
    }

    return cleanedCities
}

// ConvertCitiesToLocations converts City models to Location models
func ConvertCitiesToLocations(cities []models.City) []models.Location {
    locations := make([]models.Location, 0, len(cities))
    
    for _, city := range cities {
        if city.CityName == "" || city.CityID == "" {
            continue
        }
        
        location := models.Location{
            ID:          GenerateUniqueLocationID(city),
            CityName:    city.CityName,
            Country:     city.Country,
            CountryCode: ExtractCountryCode(city.Country),
            Latitude:    0.0,
            Longitude:   0.0,
        }
        
        locations = append(locations, location)
    }
    
    return locations
}

func GenerateLocationQueries() []string {
    queries := []string{}
    
    // Alphabet queries
    for char := 'A'; char <= 'Z'; char++ {
        queries = append(queries, string(char))
    }

    // Common prefixes and patterns
    prefixes := []string{
        "a", "the", "new", "old", "big", "small", 
        "north", "south", "east", "west", "central",
    }

    for _, prefix := range prefixes {
        queries = append(queries, prefix)
    }

    // Static example cities (you can modify or extend this list)
    exampleCities := []string{
        "New York", "Los Angeles", "London", "Tokyo",
    }

    queries = append(queries, exampleCities...)
    
    return queries
}


// ConvertToCountryMap converts country and city data into a map format.
func ConvertToCountryMap(countryCities map[string][]string) map[string]bool {
    countryMap := make(map[string]bool)
    for country := range countryCities {
        countryMap[country] = true
    }
    return countryMap
}

// ConvertToCityMap converts country and city data into a map with unique cities as keys.
func ConvertToCityMap(countryCities map[string][]string) map[string]bool {
    cityMap := make(map[string]bool)
    for _, cities := range countryCities {
        for _, city := range cities {
            cityMap[city] = true
        }
    }
    return cityMap
}

// package utils

// import (
//     "backend_rental/models"
//     "fmt"
// )

// // GenerateUniqueLocationID generates a unique ID for a location
// func GenerateUniqueLocationID(city models.City) string {
//     return fmt.Sprintf("%s-%s", city.CityID, city.Country)
// }

// // ExtractCountryCode extracts a country code from the country name
// func ExtractCountryCode(country string) string {
//     countryCodes := map[string]string{
//         "United States": "US",
//         "Canada":        "CA",
//         "United Kingdom": "GB",
//     }

//     return countryCodes[country]
// }

// // RemoveDuplicateLocations removes duplicate locations from a slice
// func RemoveDuplicateLocations(locations []models.Location) []models.Location {
//     seen := make(map[string]bool)
//     uniqueLocations := []models.Location{}

//     for _, location := range locations {
//         if _, found := seen[location.ID]; !found {
//             seen[location.ID] = true
//             uniqueLocations = append(uniqueLocations, location)
//         }
//     }

//     return uniqueLocations
// }

// // FilterAndCleanCities filters out cities that are missing critical information
// func FilterAndCleanCities(cities []models.City) []models.City {
//     var cleanedCities []models.City

//     for _, city := range cities {
//         if city.CityID == "" || city.CityName == "" || city.Country == "" {
//             continue
//         }

//         cleanedCities = append(cleanedCities, city)
//     }

//     return cleanedCities
// }

// // ConvertCitiesToLocations converts City models to Location models
// func ConvertCitiesToLocations(cities []models.City) []models.Location {
//     locations := make([]models.Location, 0, len(cities))
    
//     for _, city := range cities {
//         if city.CityName == "" || city.CityID == "" {
//             continue
//         }
        
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

// package utils

// import (
//     "backend_rental/models"
//     "fmt"
// )

// // GenerateUniqueLocationID generates a unique ID for a location
// func GenerateUniqueLocationID(city models.City) string {
//     // Example: You can create a unique ID based on the city name and country
//     return fmt.Sprintf("%s-%s", city.CityID, city.Country)
// }

// // ExtractCountryCode extracts a country code from the country name
// func ExtractCountryCode(country string) string {
//     // Example: You can use a map or an API to get the country code
//     countryCodes := map[string]string{
//         "United States": "US",
//         "Canada":        "CA",
//         "United Kingdom": "GB",
//         // Add other country mappings here
//     }

//     return countryCodes[country] // Default behavior if not found can be adjusted
// }

// // ConvertCitiesToLocations converts City models to Location models
// func ConvertCitiesToLocations(cities []models.City) []models.Location {
//     locations := make([]models.Location, 0, len(cities))
    
//     for _, city := range cities {
//         if city.CityName == "" || city.CityID == "" {
//             continue
//         }
        
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

// -----------------------------------------------------------------
// // utils/location_utils.go
// package utils

// import (
//     "backend_rental/models"
// )

// // ConvertCitiesToLocations converts City models to Location models
// func ConvertCitiesToLocations(cities []models.City) []models.Location {
//     locations := make([]models.Location, 0, len(cities))
    
//     for _, city := range cities {
//         if city.CityName == "" || city.CityID == "" {
//             continue
//         }
        
//         location := models.Location{
//             ID:          (city),
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
