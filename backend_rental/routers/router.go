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
            beego.NSRouter("/list", &controllers.LocationController{}, "get:List"),
            beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),

            // Add property details routes here
            beego.NSRouter("/details", &controllers.PropertyDetailsController{}, "post:GetPropertyDetails"),
            beego.NSRouter("/description", &controllers.PropertyDescriptionController{}, "get:GetPropertyDescription"),
            beego.NSRouter("/images", &controllers.PropertyImageController{}, "get:GetPropertyDetails"),
            
        ),
    )
    beego.AddNamespace(ns)
}

// // routers/router.go
// package routers

// import (
//     "backend_rental/controllers"
//     beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
//     ns := beego.NewNamespace("/v1",
//         beego.NSNamespace("/booking",
//             beego.NSInclude(
//                 &controllers.BookingController{},
//             ),
//             beego.NSRouter("/summary", &controllers.BookingController{}, "get:Get"),
//             beego.NSRouter("/process", &controllers.BookingController{}, "get:Get"),
//         ),
        
//         beego.NSNamespace("/property",
//             // Existing location routes
//             beego.NSRouter("/list", &controllers.LocationController{}, "get:List"),
//             beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),
//         ),
//     )
//     beego.AddNamespace(ns)
// }
