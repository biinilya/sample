package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"fmt"
	"toptal/tests"

	_ "toptal/tests/controllers/routers"
)

func SetUp(path string, ctrl beego.ControllerInterface) {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace(path,
			beego.NSInclude(
				ctrl,
			),
		),
	)
	beego.AddNamespace(ns)
	tests.EnableDB()
}

func TearDown() {
	for _, tName := range []string{"user"} {
		if _, dbErr := orm.NewOrm().Raw(
			fmt.Sprintf(`TRUNCATE TABLE "%s" CASCADE`, tName),
		).Exec(); dbErr != nil {
			panic(dbErr)
		}
	}
}
