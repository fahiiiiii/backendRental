package services

import (
    "backend_rental/common"
    "log"
)

// Remove PropertyService declaration here if already in property_service.go
func (ps *PropertyService) InitDatabase() error {
    log.Println("Initializing database in services...")
    return nil
}

func (ps *PropertyService) FetchData() ([]common.Data, error) {
    return []common.Data{}, nil
}
