package routers

import (
	"rental_view/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Existing routes
	beego.Router("/", &controllers.PropertyController{}, "get:Index")
	beego.Router("/v1/property/list", &controllers.PropertyController{}, "get:List")
	
	// Add new route for property details
	beego.Router("/v1/property/details", &controllers.PropertyDetailsController{}, "get:Details")
}