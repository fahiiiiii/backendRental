// services/location_service.go
package services

import (
    "fmt"
	"log"
    "backend_rental/models"
    "backend_rental/utils"
    "github.com/beego/beego/v2/client/orm"
)

type LocationService struct{}

func (s *LocationService) GetLocations(
    page, 
    pageSize int, 
    country, 
    cityName string,
) ([]models.Location, int64, error) {
    o := orm.NewOrm()
    
    // Create a query table
    qs := o.QueryTable(new(models.Location))
    
    // Apply filters if provided
    if country != "" {
        qs = qs.Filter("country", country)
    }
    
    if cityName != "" {
        qs = qs.Filter("city_name__icontains", cityName)
    }
    
    // Count total matching locations
    totalCount, err := qs.Count()
    if err != nil {
        return nil, 0, err
    }
    
    // Retrieve paginated locations
    var locations []models.Location
    _, err = qs.Limit(pageSize, (page-1)*pageSize).All(&locations)
    if err != nil {
        return nil, 0, err
    }
    
    return locations, totalCount, nil
}

// Bulk create or update locations
func (s *LocationService) BulkCreateLocations(locations []models.Location) error {
    o := orm.NewOrm()
    
    // Start transaction
    tx, err := o.Begin()
    if err != nil {
        return err
    }

    // Use a defer to handle potential rollback
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Remove duplicates using utility function
    uniqueLocations := utils.RemoveDuplicateLocations(locations)

    for _, location := range uniqueLocations {
        // Skip invalid locations
        if location.ID == "" || location.CityName == "" {
            continue
        }

        // Upsert logic
        _, err := o.InsertOrUpdate(&location)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("error upserting location %s: %v", location.ID, err)
        }
    }

    // Commit transaction
    return tx.Commit()
}

// Optional: Get a single location by ID
func (s *LocationService) GetLocationByID(cityID string) (*models.Location, error) {
    o := orm.NewOrm()
    
    location := &models.Location{ID: cityID}
    err := o.Read(location)
    if err != nil {
        if err == orm.ErrNoRows {
            return nil, fmt.Errorf("location not found")
        }
        return nil, err
    }
    
    return location, nil
}

// // Process and store cities
// func (s *LocationService) ProcessAndStoreCities(cities []models.City) error {
//     // Filter and clean cities using utility function
//     cleanedCities := utils.FilterAndCleanCities(cities)
    
//     // Convert cities to locations
//     locations := s.convertCitiesToLocationsInternal(cleanedCities)
    
//     // Bulk create
//     return s.BulkCreateLocations(locations)
// }
func (s *LocationService) ProcessAndStoreCities(cities []models.City) error {
    // Add logging
    log.Printf("Processing %d cities", len(cities))
    
    // Filter and clean cities
    cleanedCities := utils.FilterAndCleanCities(cities)
    log.Printf("After cleaning: %d cities", len(cleanedCities))
    
    // Convert to locations
    locations := s.convertCitiesToLocationsInternal(cleanedCities)
    log.Printf("Converted to %d locations", len(locations))
    
    // Add validation
    if len(locations) == 0 {
        return fmt.Errorf("no valid locations to store")
    }
    
    // Bulk create with error handling
    err := s.BulkCreateLocations(locations)
    if err != nil {
        log.Printf("Error storing locations: %v", err)
        return err
    }
    
    log.Printf("Successfully stored %d locations", len(locations))
    return nil
}
// Internal method to convert cities to locations
func (s *LocationService) convertCitiesToLocationsInternal(cities []models.City) []models.Location {
    locations := make([]models.Location, 0, len(cities))
    
    for _, city := range cities {
        // Skip invalid cities
        if city.CityName == "" || city.CityID == "" {
            continue
        }
        
        location := models.Location{
            ID:          utils.GenerateUniqueLocationID(city),
            CityName:    city.CityName,
            Country:     city.Country,
            CountryCode: utils.ExtractCountryCode(city.Country),
            Latitude:    0.0, // Add logic to fetch latitude if needed
            Longitude:   0.0, // Add logic to fetch longitude if needed
        }
        
        locations = append(locations, location)
    }
    
    return locations
}

// Get unique countries and cities
func (s *LocationService) GetUniqueCountriesAndCities() (map[string][]string, error) {
    o := orm.NewOrm()
    
    // Query to get unique countries and their cities
    var results []orm.Params
    _, err := o.Raw(`
        SELECT DISTINCT country, 
               ARRAY_AGG(DISTINCT city_name) AS cities 
        FROM location 
        GROUP BY country
    `).Values(&results)
    
    if err != nil {
        return nil, err
    }
    
    countryCities := make(map[string][]string)
    for _, result := range results {
        country := result["country"].(string)
        cities := result["cities"].([]string)
        countryCities[country] = cities
    }
    log.Printf("Fetched countries and cities: %v", countryCities)

    return countryCities, nil
}













// // services/location_service.go
// package services

// import (
// 	"fmt"
// 	"backend_rental/models"
// 	"github.com/beego/beego/v2/client/orm"
// )

// type LocationServiceInterface interface {
// 	GetProperties(filter PropertyFilter) ([]models.RentalProperty, error)
// }

// type LocationService struct{}

// type PropertyFilter struct {
// 	CityID      string
// 	Type        string
// 	MinBedrooms string
// }

// func NewLocationService() LocationServiceInterface {
// 	return &LocationService{}
// }

// func (s *LocationService) GetProperties(filter PropertyFilter) ([]models.RentalProperty, error) {
// 	o := orm.NewOrm()
	
// 	// Start building the query
// 	qs := o.QueryTable("rental_property")
	
// 	// Apply filters
// 	if filter.CityID != "" {
// 		qs = qs.Filter("city_id", filter.CityID)
// 	}
// 	if filter.Type != "" {
// 		qs = qs.Filter("type", filter.Type)
// 	}
// 	if filter.MinBedrooms != "" {
// 		qs = qs.Filter("bedrooms__gte", filter.MinBedrooms)
// 	}
	
// 	// Execute query
// 	var properties []models.RentalProperty
// 	_, err := qs.All(&properties)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch properties: %v", err)
// 	}
	
// 	return properties, nil
// }