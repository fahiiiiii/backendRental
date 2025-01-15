package controllers

import (
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyDescriptionController struct {
    beego.Controller
    propDescService *services.PropDescService
}

func (c *PropertyDescriptionController) Prepare() {
    c.propDescService = services.NewPropDescService()
}

// @router /v1/property/details [get]
func (c *PropertyDescriptionController) GetPropertyDescription() {
    destID := c.GetString("dest_id")
    if destID == "" {
        c.Data["json"] = map[string]interface{}{
            "error": "dest_id is required",
        }
        c.ServeJSON()
        return
    }

    details, err := c.propDescService.GetPropertyDescription(destID)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "error": err.Error(),
        }
        c.ServeJSON()
        return
    }

    if details == nil {
        c.Data["json"] = map[string]interface{}{
            "error": "Property details not found",
        }
        c.ServeJSON()
        return
    }

    c.Data["json"] = details
    c.ServeJSON()
}