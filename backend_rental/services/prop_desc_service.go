package services

import (
    "context"
    "sync"
    "time"
    "backend_rental/models"
    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/core/logs"
    "golang.org/x/time/rate"
)

type PropDescService struct {
    mutex       sync.Mutex
    rateLimiter *rate.Limiter
}

func NewPropDescService() *PropDescService {
    return &PropDescService{
        rateLimiter: rate.NewLimiter(rate.Every(time.Second), 1), // 1 request per second
    }
}

func (s *PropDescService) GetPropertyDescription(destID string) (*models.PropertyDescription, error) {
    o := orm.NewOrm()
    details := models.PropertyDescription{DestID: destID}
    err := o.Read(&details)
    if err == orm.ErrNoRows {
        return nil, nil
    }
    return &details, err
}

func (s *PropDescService) SavePropertyDescription(details *models.PropertyDescription) error {
    o := orm.NewOrm()
    _, err := o.InsertOrUpdate(details)
    return err
}

// fetchPropertyDescriptionFromAPI fetches property description from external API
func (s *PropDescService) fetchPropertyDescriptionFromAPI(destID string) (*models.PropertyDescription, error) {
    return &models.PropertyDescription{
        DestID:      destID,
        Description: "Sample description",
        Rating:      4.5,
        Review:      "Sample review",
    }, nil
}

func (s *PropDescService) ProcessPropertyDescriptions(destIDs []string) error {
    results := make(chan *models.PropertyDescription, len(destIDs))
    var wg sync.WaitGroup

    for _, destID := range destIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            
            err := s.rateLimiter.Wait(context.Background())
            if err != nil {
                logs.Error("Rate limiter error for property %s: %v", id, err)
                return
            }

            details, err := s.fetchPropertyDescriptionFromAPI(id)
            if err != nil {
                logs.Error("Error fetching details for property %s: %v", id, err)
                return
            }

            results <- details
        }(destID)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for details := range results {
        if err := s.SavePropertyDescription(details); err != nil {
            logs.Error("Error saving property details: %v", err)
        }
    }

    return nil
}