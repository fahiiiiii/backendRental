package main

import (
    "backend_rental/utils"
    _ "backend_rental/routers"
    "log"
    beego "github.com/beego/beego/v2/server/web"

)


func main() {
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }

	if err := utils.InitDatabaseFromConfig(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    log.Println("Service initialized")
    log.Println("Starting application...")
    beego.Run()
}


