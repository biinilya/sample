// @APIVersion 1.0.0
// @Title sample demo app (Run Keeper)
// @Description This RunKeeper API allows you to manage users run data
// @Contact me@ilyabiin.com
package routers

import (
	"sample/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/user/:uid/record",
			beego.NSInclude(
				&controllers.RecordController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
