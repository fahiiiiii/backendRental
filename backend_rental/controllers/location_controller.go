
package controllers

import (
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type LocationController struct {
    beego.Controller
}

func (c *LocationController) List() {
    // Pagination parameters
    page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 20)


    // Optional filtering parameters
    country := c.GetString("country")
    cityName := c.GetString("city_name")

    locationService := &services.LocationService{}

    // Fetch locations with optional filtering
    locations, totalCount, err := locationService.GetLocations(page, pageSize, country, cityName)
    if err != nil {
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    // Prepare response
    response := map[string]interface{}{
        "locations": locations,
        "pagination": map[string]interface{}{
            "total_count": totalCount,
            "page":        page,
            "page_size":   pageSize,
        },
    }

    c.Data["json"] = response
    c.ServeJSON()
}

// New endpoint to get unique countries and cities
func (c *LocationController) GetCountriesAndCities() {
    locationService := &services.LocationService{}

    countryCities, err := locationService.GetUniqueCountriesAndCities()
    if err != nil {
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    c.Data["json"] = countryCities
    c.ServeJSON()
}



// package controllers

// import (
//     "backend_rental/services"
//     beego "github.com/beego/beego/v2/server/web"
// )

// type LocationController struct {
//     beego.Controller
// }

// func (c *LocationController) List() {
//     // Pagination parameters
//     page,  := c.GetInt("page", 1)
//     pageSize,  := c.GetInt("page_size", 20)

//     // Optional filtering parameters
//     country := c.GetString("country")
//     cityName := c.GetString("city_name")

//     locationService := &services.LocationService{}

//     // Fetch locations with optional filtering
//     locations, totalCount, err := locationService.GetLocations(page, pageSize, country, cityName)
//     if err != nil {
//         c.Data["json"] = map[string]string{"error": err.Error()}
//         c.Ctx.Output.SetStatus(500)
//         c.ServeJSON()
//         return
//     }

//     // Prepare response
//     response := map[string]interface{}{
//         "locations": locations,
//         "pagination": map[string]interface{}{
//             "total_count": totalCount,
//             "page":        page,
//             "page_size":   pageSize,
//         },
//     }

//     c.Data["json"] = response
//     c.ServeJSON()
// }

// // New endpoint to get unique countries and cities
// func (c *LocationController) GetCountriesAndCities() {
//     locationService := &services.LocationService{}

//     countryCities, err := locationService.GetUniqueCountriesAndCities()
//     if err != nil {
//         c.Data["json"] = map[string]string{"error": err.Error()}
//         c.Ctx.Output.SetStatus(500)
//         c.ServeJSON()
//         return
//     }

//     c.Data["json"] = countryCities
//     c.ServeJSON()
// }