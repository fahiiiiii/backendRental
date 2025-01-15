package routers

import (
	"view/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}

package routers

import (
	"rental_view/controllers" // import the correct controllers package
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Create a namespace for /v1
	ns := beego.NewNamespace("/v1",

		// Create a subnamespace for /v1/destinations
		beego.NSNamespace("/destinations",
			// Define the GET route for /v1/destinations
			beego.NSRouter("/", &controllers.DestinationsController{}, "get:GetDestinations"),
		),
	)

	// Add the namespace to Beego
	beego.AddNamespace(ns)
}
