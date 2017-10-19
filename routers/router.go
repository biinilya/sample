// @APIVersion 1.0.0
// @Title Toptal sample project
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"toptal/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/permissions",
			beego.NSInclude(
				&controllers.PermissionsController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
