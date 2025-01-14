package main

import (
	"backend_rental/utils"
	"log"
	_ "backend_rental/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Ensure InitDatabase is called only once during application startup
	utils.InitDatabase()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	log.Println("Starting application...")
	beego.Run()
}
