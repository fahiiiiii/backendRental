package models

import (
    "time"
    // "github.com/beego/beego/v2/client/orm"
)

type Location struct {
    ID          string    `orm:"column(id);pk;unique"`
    CityName    string    `orm:"column(city_name)"`
    Country     string    `orm:"column(country)"`
    CountryCode string    `orm:"column(country_code)"`
    Latitude    float64   `orm:"column(latitude);null"`
    Longitude   float64   `orm:"column(longitude);null"`
    CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
    UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}




// // models/location.go
// package models

// import (
//     "time"
//     // "github.com/beego/beego/v2/client/orm"
// )

// type Location struct {
//     ID          string    `orm:"column(city_id);pk"`
//     CityName    string    `orm:"column(city_name);index"`
//     Country     string    `orm:"column(country);index"`
//     CountryCode string    `orm:"column(country_code)"`
//     Latitude    float64   `orm:"column(latitude)"`
//     Longitude   float64   `orm:"column(longitude)"`
//     CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
//     UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
// }
