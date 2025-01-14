package main

import (
	_ "rental_view/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

