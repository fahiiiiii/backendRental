// controllers/destinations.go
package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/beego/beego/v2/server/web"
)

type DestinationsController struct {
    web.Controller
}

// Match the actual API response structure
type PropertyListResponse struct {
    Locations []Location `json:"locations"`
    Pagination struct {
        Page       int `json:"page"`
        PageSize   int `json:"page_size"`
        TotalCount int `json:"total_count"`
    } `json:"pagination"`
}

type Location struct {
    // Add fields that exist in your locations array
    // For now we'll just include city_name
    CityName string `json:"city_name"`
}

func (c *DestinationsController) Get() {
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

    var propertyList PropertyListResponse
    if err := json.NewDecoder(resp.Body).Decode(&propertyList); err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": "Failed to parse response: " + err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    // Convert locations to the format expected by the frontend
    cities := make([]map[string]string, 0)
    for _, location := range propertyList.Locations {
        cities = append(cities, map[string]string{
            "city_name": location.CityName,
        })
    }

    c.Data["json"] = cities
    c.ServeJSON()
}

func (c *PropertyController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "property_details.tpl"
}