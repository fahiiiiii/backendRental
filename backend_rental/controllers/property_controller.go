package controllers

import (
    "strconv"
    "backend_rental/services"
    "backend_rental/utils"
    beego "github.com/beego/beego/v2/server/web"
)

// PropertyController handles property-related endpoints
type PropertyController struct {
    beego.Controller
    propertyService services.PropertyServiceInterface
}

// NewPropertyController creates a new instance of PropertyController
func NewPropertyController() *PropertyController {
    // Initialize the database connection here
    db := utils.GetDB() // Assuming GetDB() returns a *gorm.DB instance
    rapidAPIKey, _ := beego.AppConfig.String("rapidapikey") // Fetch API key from config
    
    return &PropertyController{
        propertyService: services.NewPropertyService(db, rapidAPIKey),
    }
}

// GetProperties handles fetching properties
// @router /fetch [get]
func (c *PropertyController) GetProperties() {
    page, _ := strconv.Atoi(c.GetString("page", "1"))
    pageSize, _ := strconv.Atoi(c.GetString("pageSize", "10"))

    properties, total, err := c.propertyService.ListProperties(page, pageSize)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]interface{}{
        "data": properties,
        "total": total,
        "page": page,
        "pageSize": pageSize,
    }
    c.ServeJSON()
}

// GetPropertiesSummary handles the summary endpoint
// @router /summary [get]
func (c *PropertyController) GetPropertiesSummary() {
    page, _ := strconv.Atoi(c.GetString("page", "1"))
    pageSize, _ := strconv.Atoi(c.GetString("pageSize", "10"))

    properties, _, err := c.propertyService.ListProperties(page, pageSize)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Create a summary of the properties
    summary := make(map[string]int)
    for _, prop := range properties {
        if prop.CityName != "" { // Ensure CityName is valid
            summary[prop.CityName]++ // Count properties by city
        }
    }

    c.Data["json"] = map[string]interface{}{
        "data": summary,
    }
    c.ServeJSON()
}

// ProcessAllProperties handles batch processing of properties
// @router /process-all [get]
func (c *PropertyController) ProcessAllProperties() {
    page, _ := strconv.Atoi(c.GetString("page", "1"))
    pageSize, _ := strconv.Atoi(c.GetString("pageSize", "10"))

    properties, _, err := c.propertyService.ListProperties(page, pageSize)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }

    processedData := map[string]interface{}{
        "total_cities": len(properties),
        "total_properties": 0,
    }

    // Count total properties across all cities
    processedData["total_properties"] = len(properties)

    c.Data["json"] = map[string]interface{}{
        "status": "success",
        "data": processedData,
    }
    c.ServeJSON()
}



// @router /v1/property/list [get]
func (c *PropertyController) ListProperties() {
    page, _ := strconv.Atoi(c.GetString("page", "1"))
    pageSize, _ := strconv.Atoi(c.GetString("pageSize", "10"))
    
    properties, total, err := c.propertyService.ListProperties(page, pageSize)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }
    
    c.Data["json"] = map[string]interface{}{
        "data": properties,
        "total": total,
        "page": page,
        "pageSize": pageSize,
    }
    c.ServeJSON()
}
