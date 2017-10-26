package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["toptal/controllers:RecordController"] = append(beego.GlobalControllerRouter["toptal/controllers:RecordController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:RecordController"] = append(beego.GlobalControllerRouter["toptal/controllers:RecordController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:RecordController"] = append(beego.GlobalControllerRouter["toptal/controllers:RecordController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:record_id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:RecordController"] = append(beego.GlobalControllerRouter["toptal/controllers:RecordController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:record_id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:RecordController"] = append(beego.GlobalControllerRouter["toptal/controllers:RecordController"],
		beego.ControllerComments{
			Method: "WeeklyReport",
			Router: `/report/weekly`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid/credentials`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "ListPerm",
			Router: `/:uid/permission`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddPerm",
			Router: `/:uid/permission/:title`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "DelPerm",
			Router: `/:uid/permission/:title`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["toptal/controllers:UserController"] = append(beego.GlobalControllerRouter["toptal/controllers:UserController"],
		beego.ControllerComments{
			Method: "SignIn",
			Router: `/sign_in`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
