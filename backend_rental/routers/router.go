// routers/router.go
package routers

import (
    "backend_rental/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    ns := beego.NewNamespace("/v1",
        beego.NSNamespace("/booking",
            beego.NSInclude(
                &controllers.BookingController{},
            ),
            beego.NSRouter("/summary", &controllers.BookingController{}, "get:Get"),
            beego.NSRouter("/process", &controllers.BookingController{}, "get:Get"),
        ),
        beego.NSNamespace("/properties",
            beego.NSInclude(
                &controllers.PropertyController{},
            ),
            // Routes for the new PropertyController
            beego.NSRouter("/fetch", &controllers.PropertyController{}, "get:GetProperties"),
            beego.NSRouter("/summary", &controllers.PropertyController{}, "get:GetPropertiesSummary"),
            // Add the new batch processing endpoint
            beego.NSRouter("/process-all", &controllers.PropertyController{}, "get:ProcessAllProperties"),
        ),
        beego.NSNamespace("/property",
            // Existing location routes
            beego.NSRouter("/list", &controllers.LocationController{}, "get:List"),
            beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),
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
// 		beego.NSNamespace("/properties",
//             beego.NSInclude(
//                 &controllers.PropertyController{},
//             ),
//             beego.NSRouter("/fetch", &controllers.PropertyController{}, "get:Get"),
//             beego.NSRouter("/summary", &controllers.PropertyController{}, "get:Get"),
//         ),
//         beego.NSNamespace("/property",
//             // Route to list locations
//             beego.NSRouter("/list", &controllers.LocationController{}, "get:List"),
//             // Route to get unique countries and cities
//             beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),
//         ),
//     )
//     beego.AddNamespace(ns)
// }

// // routers/router.go
// // package routers

// // import (
// //     "backend_rental/controllers"
// //     beego "github.com/beego/beego/v2/server/web"
// // )

// // func init() {
// //     ns := beego.NewNamespace("/v1",
// //         // Booking routes
// //         beego.NSNamespace("/booking",
// //             beego.NSInclude(
// //                 &controllers.BookingController{},
// //             ),
// //             beego.NSRouter("/summary", &controllers.BookingController{}, "get:Get?action=summary"),
// //             beego.NSRouter("/process", &controllers.BookingController{}, "get:Get?action=process"),
// //         ),

// //         // Properties routes
// //         beego.NSNamespace("/properties",
// //             beego.NSInclude(
// //                 &controllers.PropertyController{},
// //             ),
// //             beego.NSRouter("/fetch", &controllers.PropertyController{}, "get:Get"),
// //             beego.NSRouter("/summary", &controllers.PropertyController{}, "get:Get"),
// //         ),

// //         // Location routes - using the new LocationController
// //         beego.NSNamespace("/locations",
// //             beego.NSInclude(
// //                 &controllers.LocationController{},
// //             ),
// //             // Basic CRUD operations
// //             beego.NSRouter("/", &controllers.LocationController{}, "get:Get?action=list"),
// //             beego.NSRouter("/summary", &controllers.LocationController{}, "get:Get?action=summary"),
// //             beego.NSRouter("/process", &controllers.LocationController{}, "get:Get?action=process"),
            
// //             // Additional location-specific endpoints
// //             beego.NSRouter("/countries-cities", &controllers.LocationController{}, "get:GetCountriesAndCities"),
// //         ),
// //     )
    
// //     beego.AddNamespace(ns)
// // }

// // ----------------------------------------------------------

// // // routers/router.go
// // package routers

// // import (
// // 	"backend_rental/controllers"

// // 	beego "github.com/beego/beego/v2/server/web"
// // )

// // func init() {
// // 	// Create a new namespace for API versioning
// // 	ns := beego.NewNamespace("/v1",
// // 		// Home route
// // 		beego.NSNamespace("/",
// // 			beego.NSInclude(
// // 				&controllers.HomeController{},
// // 			),
// // 		),
		
// // 		// Booking-related routes
// // 		beego.NSNamespace("/booking",
// // 			beego.NSInclude(
// // 				&controllers.BookingController{},
// // 			),
// // 			// Add specific routes for search
// // 			beego.NSRouter("/search", &controllers.BookingController{}, "get:SearchCities"),
// // 		),
		
// // 		// Cities routes
// // 		beego.NSNamespace("/cities",
// // 			beego.NSRouter("/search", &controllers.BookingController{}, "get:SearchCities"),
// // 		),
// // 	)

// // 	// Add the namespace to Beego
// // 	beego.AddNamespace(ns)

// // 	// Optional: Add root path handler
// // 	beego.Router("/", &controllers.HomeController{})
// // }