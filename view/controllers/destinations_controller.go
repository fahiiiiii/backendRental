package controllers

import (
    beego "github.com/beego/beego/v2/server/web"
    "net/http"
    "encoding/json"
)

type DestinationsController struct {
    beego.Controller
}

type PropertyResponse struct {
    Data []struct {
        CityName string `json:"city_name"`
    } `json:"data"`
}

type Destination struct {
    CityName string `json:"city_name"`
}

// GetDestinations handles the /v1/destinations route
func (c *DestinationsController) GetDestinations() {
    // Fetch data from the /v1/property/list API here
    resp, err := http.Get("http://localhost:8080/v1/property/list")
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to fetch data from property API"}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    var propertyResponse PropertyResponse
    err = json.NewDecoder(resp.Body).Decode(&propertyResponse)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to parse property API response"}
        c.ServeJSON()
        return
    }

    var destinations []Destination
    for _, data := range propertyResponse.Data {
        destinations = append(destinations, Destination{CityName: data.CityName})
    }

    c.Data["json"] = destinations
    c.ServeJSON()
}
func (c *MainController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
