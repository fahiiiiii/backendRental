package controllers

import (
	"encoding/json"
	"net/http"
	"backend_rental/services"
	// "backend_rental/models"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailsController struct {
	beego.Controller
}
func (c *PropertyDetailsController) GetPropertyDetails() {
	var request struct {
		PropertyIDs []string `json:"property_ids"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.PropertyIDs) == 0 {
		c.CustomAbort(http.StatusBadRequest, "No property IDs provided")
		return
	}

	// Fetch details from the database
	service := services.NewPropertyDetailsService()
	propertyDetails, err := service.FetchStoredPropertyDetails(request.PropertyIDs)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = propertyDetails
	c.ServeJSON()
}
