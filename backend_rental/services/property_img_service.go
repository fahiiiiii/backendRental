// services/property_img_service.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    beego "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/client/orm"
    "backend_rental/models"
    "backend_rental/utils/ratelimiter"
)

type PropertyImageService struct {
    rateLimiter *ratelimiter.APIRateLimiter
    rapidAPIKey string
}

// NewPropertyImageService creates and returns a new PropertyImageService instance
func NewPropertyImageService() (*PropertyImageService, error) {
    rapidAPIKey, err := beego.AppConfig.String("rapidapikey")
    if err != nil {
        return nil, fmt.Errorf("failed to get rapidapikey from config: %v", err)
    }
    return &PropertyImageService{
        rateLimiter: ratelimiter.GetInstance(),
        rapidAPIKey: rapidAPIKey,
    }, nil
}

func (s *PropertyImageService) GetPropertyDetails(destID string) (*models.PropertyDescription, error) {
    // First check if we have the data in our database
    o := orm.NewOrm()
    propertyDesc := &models.PropertyDescription{DestID: destID}
    err := o.Read(propertyDesc)
    if err == nil {
        return propertyDesc, nil
    }

    // If not in database, fetch from API
    images, err := s.fetchImagesFromAPI(destID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch images: %v", err)
    }

    // Convert images to JSON string for storage
    imagesJSON, err := json.Marshal(images)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal images: %v", err)
    }

    // Create new property description
    propertyDesc = &models.PropertyDescription{
        DestID:      destID,
        Images:      string(imagesJSON),
        Description: "Default description",
        Rating:      0,
        Review:      "",
    }

    // Save to database
    _, err = o.Insert(propertyDesc)
    if err != nil {
        return nil, fmt.Errorf("failed to insert into database: %v", err)
    }

    return propertyDesc, nil
}

func (s *PropertyImageService) fetchImagesFromAPI(destID string) (*models.CategorizedImages, error) {
    err := s.rateLimiter.Wait(context.Background())
    if err != nil {
        return nil, err
    }

    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/get-photos?hotelId=%s", destID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", s.rapidAPIKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var apiResponse struct {
        Data []struct {
            ID     int      `json:"id"`
            Tag    string   `json:"tag"`
            Images []string `json:"images"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        return nil, err
    }

    images := &models.CategorizedImages{}
    for _, item := range apiResponse.Data {
        switch item.Tag {
        case "Property building":
            images.PropertyBuilding = append(images.PropertyBuilding, item.Images...)
        case "Property":
            images.Property = append(images.Property, item.Images...)
        case "Room":
            images.Room = append(images.Room, item.Images...)
        }
    }

    return images, nil
}