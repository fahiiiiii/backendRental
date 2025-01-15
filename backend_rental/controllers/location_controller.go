
package controllers

import (
    "log"
    "backend_rental/models"
    "github.com/beego/beego/v2/client/orm"

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

func (c *LocationController) Debug() {
    o := orm.NewOrm()
    var locations []models.Location
    
    _, err := o.QueryTable(new(models.Location)).All(&locations)
    if err != nil {
        log.Printf("Error fetching locations: %v", err)
        c.Data["json"] = map[string]string{"error": err.Error()}
    } else {
        log.Printf("Found %d locations in database", len(locations))
        for i, loc := range locations {
            log.Printf("Location %d: %s, %s, %s", i+1, loc.CityName, loc.Country, loc.ID)
        }
        c.Data["json"] = locations
    }
    c.ServeJSON()
}
// // controllers/location_controller.go
// package controllers

// import (
// 	"net/http"
// 	beego "github.com/beego/beego/v2/server/web"
// 	"backend_rental/services"
// )

// type LocationController struct {
// 	beego.Controller
// 	locationService services.LocationServiceInterface
// }

// func NewLocationController() *LocationController {
// 	return &LocationController{
// 		locationService: services.NewLocationService(),
// 	}
// }

// // List handles GET requests for property listings
// // @router /list [get]
// func (c *LocationController) List() {
// 	// Get query parameters for filtering
// 	cityID := c.GetString("city_id")
// 	propertyType := c.GetString("type")
// 	minBedrooms := c.GetString("min_bedrooms")
	
// 	// Call service to get properties
// 	properties, err := c.locationService.GetProperties(services.PropertyFilter{
// 		CityID:      cityID,
// 		Type:        propertyType,
// 		MinBedrooms: minBedrooms,
// 	})
	
// 	if err != nil {
// 		c.Data["json"] = map[string]interface{}{
// 			"error": err.Error(),
// 			"code":  http.StatusInternalServerError,
// 		}
// 		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
// 		c.ServeJSON()
// 		return
// 	}
	
// 	c.Data["json"] = map[string]interface{}{
// 		"data": properties,
// 		"code": http.StatusOK,
// 	}
// 	c.ServeJSON()
// }