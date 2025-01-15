// controllers/property.go
package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
    web.Controller
}

type Property struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Price       float64  `json:"price"`
    Image       string   `json:"image"`
    BadgeText   string   `json:"badgeText"`
    Rating      float64  `json:"rating"`
    Bedrooms    int      `json:"bedrooms"`
    Amenities   []string `json:"amenities"`
    ReviewCount int      `json:"reviewCount"`
    CityID      string   `json:"city_id"`
}

func (c *PropertyController) GetPropertiesByCity() {
    cityID := c.GetString("city_id")
    if cityID == "" {
        c.Data["json"] = map[string]interface{}{
            "error": "City ID is required",
        }
        c.Ctx.Output.SetStatus(400)
        c.ServeJSON()
        return
    }

    // Make request to internal API
    resp, err := http.Get("http://localhost:8080/v1/property/list")
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    var properties []Property
    if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": "Failed to parse properties: " + err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    // Filter properties by city_id
    cityProperties := []Property{}
    for _, prop := range properties {
        if prop.CityID == cityID {
            cityProperties = append(cityProperties, prop)
        }
    }

    c.Data["json"] = cityProperties
    c.ServeJSON()
}