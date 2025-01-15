// controllers/property_controller.go
package controllers

import (
    "context"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

// PropertyController handles property-related endpoints
type PropertyController struct {
    beego.Controller
    propertyService services.PropertyServiceInterface
}

// NewPropertyController creates a new instance of PropertyController
func NewPropertyController() *PropertyController {
    rapidAPIKey, _ := beego.AppConfig.String("rapidapikey")
    return &PropertyController{
        propertyService: services.NewPropertyService(rapidAPIKey),
    }
}

// GetProperties handles fetching properties
// @router /fetch [get]
func (c *PropertyController) GetProperties() {
    ctx := context.Background()
    properties, err := c.propertyService.GetAllProperties(ctx)
    
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }
    
    c.Data["json"] = map[string]interface{}{
        "data": properties,
    }
    c.ServeJSON()
}

// GetPropertiesSummary handles the summary endpoint
// @router /summary [get]
func (c *PropertyController) GetPropertiesSummary() {
    ctx := context.Background()
    properties, err := c.propertyService.GetAllProperties(ctx)
    
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }
    
    // Create a summary of the properties
    summary := make(map[string]int)
    for city, props := range properties {
        summary[city] = len(props)
    }
    
    c.Data["json"] = map[string]interface{}{
        "data": summary,
    }
    c.ServeJSON()
}

// ProcessAllProperties handles batch processing of properties
// @router /process-all [get]
func (c *PropertyController) ProcessAllProperties() {
    ctx := context.Background()
    properties, err := c.propertyService.GetAllProperties(ctx)
    
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }
    
    // Process the properties
    processedData := map[string]interface{}{
        "total_cities": len(properties),
        "total_properties": 0,
    }
    
    // Count total properties across all cities
    for _, props := range properties {
        processedData["total_properties"] = processedData["total_properties"].(int) + len(props)
    }
    
    c.Data["json"] = map[string]interface{}{
        "status": "success",
        "data": processedData,
    }
    c.ServeJSON()
}
// // controllers/property_controller.go
// package controllers

// import (
//     "context"
//     // "fmt"
//     // "net/url"
//     // "sync"
    
//     // "backend_rental/models"
//     "backend_rental/services"
//     // "backend_rental/utils/ratelimiter"
    
//     beego "github.com/beego/beego/v2/server/web"
//     // "github.com/beego/beego/v2/core/logs" // Correct logs import
// )

// // PropertyController handles property-related endpoints
// type PropertyController struct {
//     beego.Controller
//     propertyService services.PropertyServiceInterface
// }

// // NewPropertyController creates a new instance of PropertyController
// func NewPropertyController() *PropertyController {
//     return &PropertyController{
//         propertyService: services.NewPropertyService(),
//     }
// }

// // GetProperties handles the GET request for fetching properties
// // @router /properties [get]
// func (c *PropertyController) GetProperties() {
//     ctx := context.Background()
//     properties, err := c.propertyService.GetAllProperties(ctx)
    
//     if err != nil {
//         c.Data["json"] = map[string]interface{}{
//             "error": err.Error(),
//         }
//         c.ServeJSON()
//         return
//     }
    
//     c.Data["json"] = map[string]interface{}{
//         "data": properties,
//     }
//     c.ServeJSON()
// }
