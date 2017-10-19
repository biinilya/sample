package main

import (
	_ "toptal/routers"

	_ "github.com/lib/pq"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase(
		"default",
		"postgres",
		"postgres://ibiin@localhost:5432/toptal?sslmode=disable",
	)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
