package controllers

import (
	"sample/database/filter"
	"sample/lib"
	"sample/models"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// ApiController defines basic operations for API resources
type ApiController struct {
	beego.Controller
}

func (c *ApiController) Authenticate() *models.User {
	var token = c.Ctx.Request.Header.Get("X-Access-Token")
	var uid, ok = lib.ValidateAccessToken(token)
	if !ok {
		c.AbortWith(401, "Invalid access token")
	}

	var u = new(models.User)
	var userErr = u.LoadById(lib.GetDB(), uid)
	if userErr != nil {
		c.AbortWith(403, "No such user")
	}
	return u
}

func (c *ApiController) LoadUser() *models.User {
	var uid, uidErr = c.GetUint64(":uid")
	if uidErr != nil {
		c.AbortWith(400, "Invalid :uid")
	}
	var u = new(models.User)
	var userErr = u.LoadById(lib.GetDB(), uid)
	if userErr != nil {
		c.AbortWith(403, "No such user")
	}
	return u
}

func (c *ApiController) RequirePerm(perm ...string) {
	var u = c.Authenticate()
	if !checkPerm(u, perm...) {
		c.AbortWith(403, "Permission denied")
	}
	return
}

func (c *ApiController) RequireOwnerOrPerm(perm ...string) {
	var uid, uidErr = c.GetUint64(":uid")
	if uidErr != nil {
		c.AbortWith(400, "Invalid :uid")
	}

	var u = c.Authenticate()
	if uint64(u.Id) == uid {
		return
	}
	if !checkPerm(u, perm...) {
		c.AbortWith(403, "Permission denied")
	}
	return
}

func (c *ApiController) AbortWith(code int, message interface{}) {
	var err models.RequestError
	err.Message = fmt.Sprintf("%+v", message)
	c.Data["json"] = err
	c.Ctx.Output.SetStatus(code)
	c.ServeJSON()
	c.StopRun()
}

func (c *ApiController) LoadFilter(fields ...string) (*orm.Condition, int, int) {
	var limit, limitErr = c.GetUint64("limit", 50)
	if limitErr != nil {
		c.AbortWith(400, "Bad limit")
	}
	var offset, offsetErr = c.GetUint64("offset", 0)
	if offsetErr != nil {
		c.AbortWith(400, "Bad offset")
	}
	var q = c.GetString("filter")
	if q == "" {
		return orm.NewCondition(), int(offset), int(limit)
	}
	var cond, condErr = filter.Filter.ParseToOrm(q, fields...)
	if condErr != nil {
		c.AbortWith(400, "Bad filter: "+condErr.Error())
	}
	return cond, int(offset), int(limit)
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
