// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mul/controllers"
	"mul/utils"
	"os"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/captcha",
			beego.NSInclude(
				&controllers.CaptchaController{},
			),
		),
	)
	beego.AddNamespace(ns)

	err := utils.Startup(utils.MainGoRoutine)
	if err != nil {
		logs.Error(err, utils.MainGoRoutine, "initApp")
		os.Exit(1)
	}

}
