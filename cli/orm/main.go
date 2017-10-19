package main

import (
	_ "toptal/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func main() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase(
		"default",
		"postgres",
		"postgres://ibiin@localhost:5432/toptal?sslmode=disable",
	)

	orm.RunCommand()
}
