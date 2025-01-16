package models

import (
    // "time"
    "gorm.io/gorm"
)

type Property struct {
    gorm.Model
    DestID   string    `json:"dest_id" gorm:"uniqueIndex"`
    Name     string    `json:"name"`
    CityID   string    `json:"city_id"`
    CityName string    `json:"city_name"`
}
    // package models
    
    // type CityKey struct {
    //     Name    string
    //     Country string
    // }
    
    // type Property struct {
    //     DestID      string `json:"dest_id"`
    //     Name        string `json:"name"`
    //     CityID      string  `json:"city_id"`
    //     Country     string `json:"country"`
    //     // Add other fields as needed from the API response
    // }
