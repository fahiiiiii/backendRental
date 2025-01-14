// controllers/home_controller.go
package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type HomeController struct {
	beego.Controller
}

// Get handles GET requests to the home route
func (c *HomeController) Get() {
	c.Data["json"] = map[string]interface{}{
		"status":  "ok",
		"message": "Welcome to Backend Rental API",
		"version": "1.0.0",
	}
	c.ServeJSON()
}

// Post handles POST requests to the home route
func (c *HomeController) Post() {
	c.Data["json"] = map[string]string{
		"error": "Method not allowed",
	}
	c.Ctx.Output.SetStatus(405)
	c.ServeJSON()
}