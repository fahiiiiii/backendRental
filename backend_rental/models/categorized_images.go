package models

type CategorizedImages struct {
    PropertyBuilding []string `json:"property_building"`
    Property        []string `json:"property"`
    Room            []string `json:"room"`
}
