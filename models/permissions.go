package models

import "github.com/astaxie/beego/orm"

type Permission struct {
	Id          int64   `orm:"auto"`
	Title       string  `orm:"size(128);unique"`
	Description string  `orm:"type(longtext)"`
	Users       []*User `orm:"reverse(many)"`
}

func PermissionByTitle(o orm.Ormer, title string) (perm *Permission, err error) {
	perm = new(Permission)
	err = o.QueryTable(perm).Filter("Title", title).One(perm)
	return
}

func init() {
	orm.RegisterModel(new(Permission))
}

const (
	PERM_USER    = "user"
	PERM_MANAGER = "manager"
	PERM_ADMIN   = "admin"
)
