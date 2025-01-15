
package models

type RentalProperty struct {
    CityID    *int     `json:"city_id" orm:"column(city_id);null"`
    DestID    *int     `json:"dest_id" orm:"column(dest_id);null"`
    Name      string   `json:"name" orm:"column(name)"`
    Type      string   `json:"type" orm:"column(type)"`
    Bedrooms  int      `json:"bedrooms" orm:"column(bedrooms)"`
    Bathroom  int      `json:"bathroom" orm:"column(bathroom)"`
    Amenities string   `json:"amenities" orm:"column(amenities)"`
}

// TableName returns the name of the table for ORM mapping.
func (r *RentalProperty) TableName() string {
	return "RentalProperty"
}