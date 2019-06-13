package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mul/services"
	mongo "mul/utils"
	"runtime"
)

type BaseController struct {
	beego.Controller
	services.Service
}

/**
开获取请求前， 根据userid创建MongoDB连接session
*/

func (baseController *BaseController) Prepare() {
	baseController.UserID = baseController.GetString("userID")
	if baseController.UserID == "" {
		baseController.UserID = baseController.GetString(":userID")
	}
	if baseController.UserID == "" {
		baseController.UserID = "Unknown"
	}

	if err := baseController.Service.Prepare(); err != nil {
		logs.Error(err, baseController.UserID, "BaseController.Prepare", baseController.Ctx.Request.URL.Path)
		baseController.ServeError(err)
		return
	}

	logs.Info(baseController.UserID, "BaseController.Prepare", fmt.Sprintf("UserID[%s] Path[%s]", baseController.UserID, baseController.Ctx.Request.URL.Path))
}
func (baseController *BaseController) Finish() {
	defer func() {
		if baseController.MongoSession != nil {
			mongo.CloseSession(baseController.UserID, baseController.MongoSession)
			baseController.MongoSession = nil
		}
	}()

	logs.Info(baseController.UserID, "Finish", baseController.Ctx.Request.URL.Path)
}

//** EXCEPTIONS

// ServeError prepares and serves an Error exception.
func (baseController *BaseController) ServeError(err error) {
	baseController.Data["json"] = struct {
		Error string `json:"Error"`
	}{err.Error()}
	baseController.Ctx.Output.SetStatus(500)
	baseController.ServeJSON()
}

// ServeValidationErrors prepares and serves a validation exception.
func (baseController *BaseController) ServeValidationErrors(Errors []string) {
	baseController.Data["json"] = struct {
		Errors []string `json:"Errors"`
	}{Errors}
	baseController.Ctx.Output.SetStatus(409)
	baseController.ServeJSON()
}

/**
捕获Controller层异常
*/
func (baseController *BaseController) CatchPanic(functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		logs.Warning(baseController.Service.UserID, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		baseController.ServeError(fmt.Errorf("%v", r))
	}
}
