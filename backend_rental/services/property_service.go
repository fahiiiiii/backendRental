// services/property_service.go
package services

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "backend_rental/models"
    "gorm.io/gorm"
)

// PropertyServiceInterface defines the methods that PropertyService must implement
type PropertyServiceInterface interface {
    ListProperties(page, pageSize int) ([]models.Property, int64, error)
    FetchAndStoreProperties() error
}

type PropertyService struct {
    db          *gorm.DB
    httpClient  *http.Client
    rapidAPIKey string
}

// NewPropertyService creates a new PropertyService with db connection and API key
func NewPropertyService(db *gorm.DB, rapidAPIKey string) PropertyServiceInterface {
    return &PropertyService{
        db:          db,
        httpClient:  &http.Client{},
        rapidAPIKey: rapidAPIKey,
    }
}

// ListProperties retrieves properties from the database with pagination
func (s *PropertyService) ListProperties(page, pageSize int) ([]models.Property, int64, error) {
    var properties []models.Property
    var total int64
    
    if err := s.db.Model(&models.Property{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    offset := (page - 1) * pageSize
    err := s.db.Offset(offset).Limit(pageSize).Find(&properties).Error
    if err != nil {
        return nil, 0, err
    }
    
    return properties, total, nil
}

// FetchAndStoreProperties fetches data from an external API and stores it in the database
func (s *PropertyService) FetchAndStoreProperties() error {
    cities := []string{"New York", "London", "Paris", "Tokyo"}
    
    for _, city := range cities {
        properties, err := s.fetchPropertiesForCity(city)
        if err != nil {
            return fmt.Errorf("error fetching properties for %s: %v", city, err)
        }
        
        for _, prop := range properties {
            err := s.db.Where(models.Property{DestID: prop.DestID}).
                Assign(prop).
                FirstOrCreate(&prop).Error
            
            if err != nil {
                return fmt.Errorf("error storing property %s: %v", prop.Name, err)
            }
        }
    }
    
    return nil
}

// fetchPropertiesForCity fetches property data from an external API for a specific city
func (s *PropertyService) fetchPropertiesForCity(city string) ([]models.Property, error) {
    encodedQuery := url.QueryEscape(city)
    apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)
    
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", s.rapidAPIKey)
    
    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var response struct {
        Data []struct {
            DestID   string `json:"dest_id"`
            Name     string `json:"name"`
            CityID   string `json:"city_id"`
            CityName string `json:"city_name"`
        } `json:"data"`
    }
    
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, err
    }
    
    properties := make([]models.Property, len(response.Data))
    for i, item := range response.Data {
        properties[i] = models.Property{
            DestID:   item.DestID,
            Name:     item.Name,
            CityID:   item.CityID,
            CityName: item.CityName,
        }
    }
    
    return properties, nil
}


// // services/property_service.go
// package services

// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "io"
//     "net/http"
//     "net/url"
//     "backend_rental/models"
//     "gorm.io/gorm"
// )

// // type PropertyServiceInterface interface {
// //     ListProperties(page, pageSize int) ([]models.Property, int64, error)
// //     FetchAndStoreProperties() error
// // }

// type PropertyService struct {
//     db          *gorm.DB
//     httpClient  *http.Client
//     rapidAPIKey string
// }

// func NewPropertyService(rapidAPIKey string) PropertyServiceInterface {
//     return &PropertyService{
//         db:          database.GetDB(), // Make sure to implement database.GetDB()
//         httpClient:  &http.Client{},
//         rapidAPIKey: rapidAPIKey,
//     }
// }

// // ListProperties retrieves properties from the database with pagination
// func (s *PropertyService) ListProperties(page, pageSize int) ([]models.Property, int64, error) {
//     var properties []models.Property
//     var total int64
    
//     // Get total count of properties
//     if err := s.db.Model(&models.Property{}).Count(&total).Error; err != nil {
//         return nil, 0, err
//     }
    
//     // Calculate offset
//     offset := (page - 1) * pageSize
    
//     // Get properties with pagination
//     err := s.db.Offset(offset).Limit(pageSize).Find(&properties).Error
//     if err != nil {
//         return nil, 0, err
//     }
    
//     return properties, total, nil
// }

// // FetchAndStoreProperties fetches data from Booking.com API and stores in database
// func (s *PropertyService) FetchAndStoreProperties() error {
//     // Cities to fetch properties for
//     cities := []string{"New York", "London", "Paris", "Tokyo"}
    
//     for _, city := range cities {
//         properties, err := s.fetchPropertiesForCity(city)
//         if err != nil {
//             return fmt.Errorf("error fetching properties for %s: %v", city, err)
//         }
        
//         // Store each property in database
//         for _, prop := range properties {
//             // Upsert - update if exists, insert if doesn't
//             err := s.db.Where(models.Property{DestID: prop.DestID}).
//                 Assign(prop).
//                 FirstOrCreate(&prop).Error
            
//             if err != nil {
//                 return fmt.Errorf("error storing property %s: %v", prop.Name, err)
//             }
//         }
//     }
    
//     return nil
// }

// func (s *PropertyService) fetchPropertiesForCity(city string) ([]models.Property, error) {
//     encodedQuery := url.QueryEscape(city)
//     apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)
    
//     req, err := http.NewRequest("GET", apiURL, nil)
//     if err != nil {
//         return nil, err
//     }
    
//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", s.rapidAPIKey)
    
//     resp, err := s.httpClient.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()
    
//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, err
//     }
    
//     var response struct {
//         Data []struct {
//             DestID   string `json:"dest_id"`
//             Name     string `json:"name"`
//             CityID   string `json:"city_id"`
//             CityName string `json:"city_name"`
//         } `json:"data"`
//     }
    
//     if err := json.Unmarshal(body, &response); err != nil {
//         return nil, err
//     }
    
//     properties := make([]models.Property, len(response.Data))
//     for i, item := range response.Data {
//         properties[i] = models.Property{
//             DestID:   item.DestID,
//             Name:     item.Name,
//             CityID:   item.CityID,
//             CityName: item.CityName,
//         }
//     }
    
//     return properties, nil
// }
