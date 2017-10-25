package controllers

import (
	"toptal/lib"
	"toptal/models"

	"toptal/database/filter"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// ApiController defines basic operations for API resources
type ApiController struct {
	beego.Controller
}

func (c *UserController) Authenticate() *models.User {
	var token = c.Ctx.Request.Header.Get("X-Access-Token")
	var uid, ok = lib.ValidateAccessToken(token)
	if !ok {
		c.Abort("401")
	}

	var u = new(models.User)
	var userErr = u.LoadById(lib.GetDB(), uid)
	if userErr != nil {
		c.Abort("403")
	}
	return u
}

func (c *UserController) RequirePerm(perm ...string) *models.User {
	var u = c.Authenticate()
	if !checkPerm(u, perm...) {
		c.Abort("403")
	}
	return u
}

func (c *UserController) RequireOwnerOrPerm(uid uint64, perm ...string) *models.User {
	var u = c.Authenticate()
	if uint64(u.Id) == uid {
		return u
	}
	if !checkPerm(u, perm...) {
		c.Abort("403")
	}
	return u
}

func (c *UserController) LoadFilter(fields ...string) *orm.Condition {
	var q = c.GetString("filter")
	if q == "" {
		return orm.NewCondition()
	}
	var cond, condErr = filter.Filter.ParseToOrm(q, fields...)
	if condErr != nil {
		c.Abort("400")
	}
	return cond
}

func checkPerm(u *models.User, perm ...string) bool {
	var granted = false
	for _, permItem := range perm {
		for _, permSlot := range u.Permissions {
			if permItem == permSlot.Title {
				granted = true
			}
		}
	}
	return granted
}
