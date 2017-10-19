package lib

import "github.com/astaxie/beego/orm"

var GetDB = OrmDefaultFactory

var OrmDefaultFactory = func() orm.Ormer {
	return orm.NewOrm()
}
