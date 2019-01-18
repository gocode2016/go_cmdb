package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "error/404.html"
}

func (c *ErrorController) Error500() {
	c.Data["content"] = "server error"
	c.TplName = "error/500.html"
}

func (c *ErrorController) Error403() {
	c.Data["content"] = "page forbidden"
	c.TplName = "error/403.html"
}

