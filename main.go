package main

import (
	"github.com/astaxie/beego"
	_ "mul/routers"
	"mul/utils"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.ErrorController(&utils.ErrorController{})
	beego.Run()
}
