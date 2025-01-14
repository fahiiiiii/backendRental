// package controllers

// import (
// 	"log"
// 	"net/http"

// 	"github.com/beego/beego/v2/server/web"
// )

// type PropertyDetailsController struct {
// 	web.Controller
// }

// func (c *PropertyDetailsController) Details() {
// 	// Enable CORS
// 	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
// 	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type")

// 	// Log incoming request details
// 	log.Printf("Incoming Request - Method: %s, Query: %v", 
// 		c.Ctx.Input.Method(), 
// 		c.Ctx.Input.Query("id"),
// 	)

// 	// Get property ID from multiple sources
// 	propertyID := c.GetString("id")
// 	if propertyID == "" {
// 		propertyID = c.Ctx.Input.Param(":id")
// 	}

// 	// Detailed logging
// 	log.Printf("Property ID Extracted: %s", propertyID)

// 	// Validate property ID
// 	if propertyID == "" {
// 		c.Ctx.Output.SetStatus(400)
// 		c.Data["json"] = map[string]interface{}{
// 			"error":   true,
// 			"message": "Property ID is required",
// 		}
// 		c.ServeJSON()
// 		return
// 	}

// 	// Mock property details with comprehensive logging
// 	log.Printf("Preparing Property Details for ID: %s", propertyID)

// 	propertyDetails := map[string]interface{}{
// 		"id":          propertyID,
// 		"title":       "Luxury Dubai Hostel with Pool",
// 		"rating":      10.0,
// 		"reviewCount": 520,
// 		"bedrooms":    1,
// 		"bathrooms":   1,
// 		"guests":      8,
// 		"description": `Stunning property in the heart of Dubai`,
// 		"images": []string{
// 			"https://images.unsplash.com/photo-1560185893-a55cbc8c5ac8",
// 			"https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd",
// 			"https://images.unsplash.com/photo-1582719478250-c89cae4dc85b",
// 		},
// 	}

// 	log.Printf("Serving Property Details: %+v", propertyDetails)

// 	// Serve JSON response
// 	c.Data["json"] = propertyDetails
// 	c.ServeJSON()
// }

// // Handle OPTIONS for CORS preflight
// func (c *PropertyDetailsController) Options() {
// 	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
// 	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type")
// 	c.Ctx.Output.SetStatus(200)
// }

// func (c *PropertyDetailsController) Index() {
// 	c.TplName = "property_details.tpl"
// }



package controllers

import (
	// "encoding/json"
	"log"

	"github.com/beego/beego/v2/server/web"
)

type PropertyDetailsController struct {
	web.Controller
}

func (c *PropertyDetailsController) Details() {
	// Enable CORS
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type")

	// Get property ID from query parameter
	propertyID := c.GetString("id")

	log.Printf("Requested Property ID: %s", propertyID)

	// Validate property ID
	if propertyID == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"error":   true,
			"message": "Property ID is required",
		}
		c.ServeJSON()
		return
	}

	// Mock property details based on ID
	propertyDetails := map[string]interface{}{
		"id":          propertyID,
		"title":       "Luxury Dubai Hostel with Pool",
		"rating":      10.0,
		"reviewCount": 520,
		"bedrooms":    1,
		"bathrooms":   1,
		"guests":      8,
		"description": `Stunning property in the heart of Dubai`,
		"amenities": []string{
			"Pool", 
			"Gym", 
			"Free WiFi",
		},
		"location": "Dubai, UAE",
		"price": 250,
		"images": []string{
			"https://images.unsplash.com/photo-1560185893-a55cbc8c5ac8",
			"https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd",
			"https://images.unsplash.com/photo-1582719478250-c89cae4dc85b",
		},
	}

	log.Printf("Serving Property Details: %+v", propertyDetails)

	// Serve JSON response
	c.Data["json"] = propertyDetails
	c.ServeJSON()
}


// package controllers

// import (
// 	"github.com/beego/beego/v2/server/web"
// )

// type PropertyDetailsController struct {
// 	web.Controller
// }

// func (c *PropertyDetailsController) Index() {
// 	c.TplName = "property_details.tpl"
// }

// func (c *PropertyDetailsController) Details() {
// 	propertyID := c.GetString("id")
	
// 	// Mock property details - replace with actual API call
// 	propertyDetails := map[string]interface{}{
// 		"id":          propertyID,
// 		"title":       "Hostel · 8 guests | Hostel in Dubai with Pool",
// 		"rating":      10.0,
// 		"reviewCount": 520,
// 		"bedrooms":    1,
// 		"bathrooms":   1,
// 		"guests":      8,
// 		"description": `
// 			Welcome to Peppers Backpackers, the modern gem of Dubai where you can experience a true backpacker's paradise. 
// 			As the only hostel in the city with a backpacker vibe, we offer an enchanting blend of comfort and affordability.

// 			Prepare to be captivated by the magnificent views that greet you every morning. 
// 			With our perfectly situated location, you are just a stone's throw away from the pristine beach. 
// 			From more than 300 meters high on our rooftop, immerse yourself in the breathtaking sight of the iconic Palm Island off on your horizon.
// 		`,
// 		"images": []string{
// 			"https://images.unsplash.com/photo-1560185893-a55cbc8c5ac8",
// 			"https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd",
// 			"https://images.unsplash.com/photo-1582719478250-c89cae4dc85b",
// 			"https://images.unsplash.com/photo-1594560913095-8cf34baf3a39",
// 			"https://images.unsplash.com/photo-1505873242700-f289a29c1d4f",
// 		},
// 	}

// 	c.Data["json"] = propertyDetails
// 	c.ServeJSON()
// }