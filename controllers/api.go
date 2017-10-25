package controllers

import (
	"toptal/database/filter"
	"toptal/lib"
	"toptal/models"

	"fmt"

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
		c.AbortWith(401, "Invalid access token")
	}

	var u = new(models.User)
	var userErr = u.LoadById(lib.GetDB(), uid)
	if userErr != nil {
		c.AbortWith(403, "No such user")
	}
	return u
}

func (c *UserController) RequirePerm(perm ...string) *models.User {
	var u = c.Authenticate()
	if !checkPerm(u, perm...) {
		c.AbortWith(403, "Permission denied")
	}
	return u
}

func (c *UserController) RequireOwnerOrPerm(uid uint64, perm ...string) *models.User {
	var u = c.Authenticate()
	if uint64(u.Id) == uid {
		return u
	}
	if !checkPerm(u, perm...) {
		c.AbortWith(403, "Permission denied")
	}
	return u
}

func (c *UserController) AbortWith(code int, message interface{}) {
	var err models.RequestError
	err.Message = fmt.Sprintf("%+v", message)
	c.Data["json"] = err
	c.Ctx.Output.SetStatus(code)
	c.ServeJSON()
	c.StopRun()
}

func (c *UserController) LoadFilter(fields ...string) *orm.Condition {
	var q = c.GetString("filter")
	if q == "" {
		return orm.NewCondition()
	}
	var cond, condErr = filter.Filter.ParseToOrm(q, fields...)
	if condErr != nil {
		c.AbortWith(400, "Bad filter: "+condErr.Error())
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

// Each time entry when entered has a date, distance, time, and location.
// Based on the provided date and location, API should connect to a weather API provider and get the weather conditions for the run, and store that with each run.
// The API must create a report on average speed & distance per week.
