package utils

import (
	"github.com/astaxie/beego"
	"net/http"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error403() {
	dfe := "Forbidden"
	statusOut(c, dfe, http.StatusForbidden)
}
func (c *ErrorController) Error404() {
	dfe := "Resource Not Find"
	statusOut(c, dfe, http.StatusNotFound)
}
func (c *ErrorController) Error406() {
	dfe := "request not can duty"
	statusOut(c, dfe, http.StatusNotAcceptable)
}

func (c *ErrorController) Error501() {
	dfe := "Not Implemented"
	statusOut(c, dfe, http.StatusNotImplemented)
}
func (c *ErrorController) Error500() {
	dfe := "Internal Server Error"
	statusOut(c, dfe, http.StatusInternalServerError)
}

func statusOut(c *ErrorController, defaulErValue string, code int) {
	content := c.Data["content"]
	s, e := checkString(content)
	if e != nil {
		s = defaulErValue
	}
	http.Error(c.Ctx.Output.Context.ResponseWriter, s, code)
}
