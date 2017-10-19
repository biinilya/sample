package controllers

import (
	"toptal/lib"

	"toptal/models"

	"github.com/astaxie/beego"
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
