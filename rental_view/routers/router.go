package routers

import (
	"rental_view/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	// In routers/router.go
	beego.Router("/destinations", &controllers.DestinationsController{})
	beego.Router("/all-properties", &controllers.PropertyController{}, "get:GetPropertiesByCity")
}
