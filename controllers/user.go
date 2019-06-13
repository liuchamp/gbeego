package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
	"mul/models"
	"mul/services"
)

type UserController struct {
	BaseController
}

// @Title Get
// @Description find user info by id
// @Param	userId		query 	string	true		"the userid you want to get"
// @Success 200 {object} models.UserDTO
// @Failure 403 userId is empty
// @router /info [get]
func (this *UserController) Get() {
	objectId := this.Ctx.Request.Form.Get("userId")
	vpas := this.Ctx.Input.Params()
	logs.Debug(vpas)
	if objectId != "" {
		ob, err := services.FindUserById(&this.Service, objectId)

		if err != nil {
			if err == mgo.ErrNotFound {
				this.Data["content"] = fmt.Sprintf("Can not find User By userId %s", objectId)
				this.Abort("404")
			}
		} else {
			rs := models.User2UserDto(*ob)
			this.Data["json"] = rs
		}
	}
	this.ServeJSON()
}

// @Title Transfer
// @Description  Transfer operation interface
// @Param	body		body 	models.TransferDTO	true		"The transfer DTO content"
// @Success 200
// @Failure 403 body is empty
// @Failure 405 operation failed
// @Failure 500 inner error
// @router /moveMoney [post]
func (this *UserController) Post() {
	var td models.TransferDTO
	e := json.Unmarshal(this.Ctx.Input.RequestBody, &td)
	if e != nil {
		this.Data["content"] = fmt.Sprintf("Body Error %s", this.Ctx.Input.RequestBody)
		this.Abort("406")
	}
	if td.Money < 0 {
		this.Data["content"] = "money Must be greater than 0"
		this.Abort("406")
	}
	err := services.MoveMoney(&this.Service, &td)
	if err != nil {
		this.Data["content"] = err.Error()
		this.Abort("500")
	}

	this.ServeJSON()
}

// @Title registered
// @Description Registered a user
// @Param	body		body 	models.RegisteredDTO	true		"The Registered DTO content"
// @Success 200
// @Failure 403 body is empty
// @Failure 406 operation failed, MS code expired or not exist
// @router /reg [post]
func (uc *UserController) Register() {
	var ob models.RegisteredDTO
	err := json.Unmarshal(uc.Ctx.Input.RequestBody, &ob)
	if err != nil {
		logs.Error("Register parse body error ->", err)
		uc.Data["json"] = err.Error()

	} else {
		if okCode := services.CheckSmsCodeExsit(ob.SmsCode); okCode {
			// code 存在
			nU := models.User{}
			nU.UserPhone = ob.UserPhone
			nU.UserName = ob.UserPhone
			nU.UserId = ob.UserPhone

			_, er := services.CreateUser(&uc.Service, &nU)
			if er != nil {
				uc.Data["json"] = er.Error()
			} else {
				uc.Data["json"] = ""
			}
		} else {
			uc.Data["content"] = "SMS code expired or not exist"
			uc.Abort("406")
		}
	}
	uc.ServeJSON()

}
