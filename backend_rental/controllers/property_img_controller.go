// controllers/property_img_controller.go
package controllers

import (
    "fmt"
    beego "github.com/beego/beego/v2/server/web"
    "backend_rental/services"
)

type PropertyImageController struct {
    beego.Controller
}

// GetPropertyDetails handles GET requests to /v1/property/images
func (c *PropertyImageController) GetPropertyDetails() {
    fmt.Println("=== Starting GetPropertyDetails ===")
    
    destID := c.GetString("dest_id")
    fmt.Printf("Received dest_id: %s\n", destID)
    
    if destID == "" {
        fmt.Println("Error: dest_id is empty")
        c.Data["json"] = map[string]interface{}{
            "error": "dest_id is required",
        }
        c.ServeJSON()
        return
    }

    fmt.Println("Initializing property service...")
    propertyService, err := services.NewPropertyImageService()
    if err != nil {
        fmt.Printf("Service initialization error: %v\n", err)
        c.Data["json"] = map[string]interface{}{
            "error": fmt.Sprintf("failed to initialize service: %v", err),
        }
        c.ServeJSON()
        return
    }

    fmt.Printf("Fetching property details for dest_id: %s\n", destID)
    propertyDetails, err := propertyService.GetPropertyDetails(destID)
    if err != nil {
        fmt.Printf("Error fetching property details: %v\n", err)
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }

    fmt.Println("Successfully retrieved property details")
    c.Data["json"] = propertyDetails
    c.ServeJSON()
    fmt.Println("=== Completed GetPropertyDetails ===")
}
