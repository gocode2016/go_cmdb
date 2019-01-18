package main

import (
	"github.com/astaxie/beego"
	_ "cmdb/models"
	_ "cmdb/routers"
	"cmdb/controllers"
)

func main() {
	beego.ErrorController(&controllers.ErrorController{})
	beego.SetLogger("file", `{"filename":"logs/cmdb.log","separate":["emergency", "alert", "critical", "error"],"daily":true,"maxdays":7,"level":7}`)
	beego.SetLogFuncCall(true)
	beego.Run()
}
