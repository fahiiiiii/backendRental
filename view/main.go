package main

import (
	_ "view/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

