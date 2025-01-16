package routers

import (
    "backend_rental/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    ns := beego.NewNamespace("/v1",

        // Booking routes
        beego.NSNamespace("/booking",
            beego.NSInclude(
                &controllers.BookingController{},
            ),
            beego.NSRouter("/summary", &controllers.BookingController{}, "get:Get"),
            beego.NSRouter("/process", &controllers.BookingController{}, "get:Get"),
        ),

        // Property routes
        beego.NSNamespace("/property",
            // Fix `/list` route to use PropertyController
            beego.NSRouter("/list", &controllers.PropertyController{}, "get:ListProperties"),
            beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),

            // Property details routes
            beego.NSRouter("/details", &controllers.PropertyDetailsController{}, "post:GetPropertyDetails"),
            beego.NSRouter("/description", &controllers.PropertyDescriptionController{}, "get:GetPropertyDescription"),
            beego.NSRouter("/images", &controllers.PropertyImageController{}, "get:GetPropertyDetails"),
        ),
    )
    beego.AddNamespace(ns)
}

