package tests

import (
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

var dbUrl string

func init() {
	testPath, _ := filepath.Abs("../tests")
	beego.TestBeegoInit(testPath)
	dbUrl = beego.AppConfig.String("DbConn")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase(
		"default",
		"postgres",
		dbUrl,
	)
	//orm.Debug = true
}

func DisableDB() {
	//lib.GetDB = mocks.BrokenOrmerFactory(errors.New("Broken DB"))
}

func EnableDB() {
	//lib.GetDB = lib.OrmDefaultFactory
}
