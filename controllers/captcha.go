package controllers

import (
	"fmt"
	"mul/services"
)

type CaptchaController struct {
	BaseController
}

// @Title Get
// @Description Get New Verification code
// @Success 200 {object} models.VerCode
// @Failure 503 : service error
// @router / [get]
func (o *CaptchaController) Get() {
	fmt.Println("start")
	smsCode := services.GetNewSmsVerCode()
	o.Data["json"] = smsCode.VerCode
	o.ServeJSON()
}
