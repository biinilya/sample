package controllers

import (
	"toptal/lib"
	"toptal/models"
	"toptal/presenter"

	"encoding/json"

	"github.com/astaxie/beego/orm"
)

//  UserController operations for User
type UserController struct {
	ApiController
}

var userView = presenter.New().AddHook(
	models.User{}, presenter.StructView().TranslateCase().Include("Id", "Key", "Created", "Updated"),
).AddHook(
	models.Permission{}, presenter.StructView().TranslateCase().Include("Title", "Description"),
)

// @Title Post
// @Description create new User, returns key and secret used to authorize
// @Success 201 {object} models.UserCredentialsView
// @router / [post]
func (c *UserController) Post() {
	var resp models.UserCredentialsView
	if opErr := lib.WithTx(lib.GetDB(), func(o orm.Ormer) (err error) {
		var v models.User
		resp.Secret, err = v.New(o)
		resp.Key = v.Key
		resp.Id = uint64(v.Id)
		return
	}); opErr != nil {
		c.AbortWith(500, opErr)
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title SignIn
// @Description Use this method to receive Access Token which is required for API access
// @Param	body		body 	models.UserCredentialsData	true		"Credentials to Sign In"
// @Success 200 {object} models.UserAccessTokenView
// @Failure 400 bad request
// @Failure 403 forbidden
// @router /sign_in [post]
func (c *UserController) SignIn() {
	var cred models.UserCredentialsView
	if dataInErr := json.Unmarshal(c.Ctx.Input.RequestBody, &cred); dataInErr != nil {
		c.AbortWith(400, dataInErr)
	}
	if cred.Key == "" || cred.Secret == "" {
		c.AbortWith(400, "Key and Secret should not be empty")
	}

	var v models.User
	if opErr := v.LoadByKey(lib.GetDB(), cred.Key); opErr != nil {
		c.AbortWith(403, "Invalid credentials")
	}
	if lib.Encode(cred.Secret) != v.Secret {
		c.AbortWith(403, "Invalid credentials")
	}

	var resp models.UserAccessTokenView
	resp.AccessToken = lib.GenerateAccessToken(uint64(v.Id))

	c.Data["json"] = resp
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Success 200 {object} models.UserInfoView
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid [get]
func (c *UserController) GetOne() {
	var u = c.Authenticate()

	c.Data["json"] = userView.AsJson(u)
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description generate new User credentials
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Success 200 {object} models.UserCredentialsView
// @Failure 400 bad request (uid is missing or is not a number)
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid/credentials [put]
func (c *UserController) Put() {
	c.RequireOwnerOrPerm(models.PERM_MANAGER, models.PERM_ADMIN)
	var u = c.LoadUser()

	var secret, dbErr = u.NewCred(lib.GetDB())
	if dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Data["json"] = models.UserCredentialsView{
		Id:     uint64(u.Id),
		Key:    u.Key,
		Secret: secret,
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Success 200 {object} models.UserInfoView
// @Failure 400 bad request (uid is missing or is not a number)
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid [delete]
func (c *UserController) Delete() {
	c.RequireOwnerOrPerm(models.PERM_MANAGER, models.PERM_ADMIN)
	var u = c.LoadUser()

	var dbErr = u.Delete(lib.GetDB())
	if dbErr != nil {
		c.AbortWith(500, dbErr)
	}
	c.Data["json"] = userView.AsJson(u)
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description Users Directory
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	filter		query 	string	false		"Filter users e.x. (key eq 'xxx')"
// @Param	offset		query 	int  	false		"Offset in records"
// @Param	limit		query 	int  	false		"Limit number of records to (default 50)"
// @Success 200 {object} []models.UserInfoView
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router / [get]
func (c *UserController) GetAll() {
	var f, offset, limit = c.LoadFilter("key", "created", "updated")
	c.RequirePerm(models.PERM_MANAGER, models.PERM_ADMIN)

	var users []*models.User
	var opErr error
	if users, opErr = models.UsersGetAll(lib.GetDB(), f, offset, limit); opErr != nil {
		c.AbortWith(500, opErr)
	}

	c.Data["json"] = userView.AsJson(users)
	c.ServeJSON()
}

// ListPerm ...
// @Title Add Permission
// @Description delete the User
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Success 200 {object} []models.PermissionView
// @Failure 400 bad request (uid is missing or is not a number)
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid/permission [get]
func (c *UserController) ListPerm() {
	c.RequirePerm(models.PERM_ADMIN)
	var u = c.LoadUser()

	c.Data["json"] = userView.AsJson(u.Permissions)
	c.ServeJSON()
}

// AddPerm ...
// @Title Add Permission
// @Description delete the User
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Param	title	path 	string	true		"Title of Permission to Add to user"
// @Success 200 {object} []models.PermissionView
// @Failure 400 bad request (uid is missing or is not a number)
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid/permission/:title [post]
func (c *UserController) AddPerm() {
	c.RequirePerm(models.PERM_ADMIN)
	var u = c.LoadUser()
	if dbErr := u.AddPermission(lib.GetDB(), c.GetString(":title")); dbErr != nil {
		c.AbortWith(400, dbErr)
	}

	c.Data["json"] = userView.AsJson(u.Permissions)
	c.ServeJSON()
}

// DelPerm ...
// @Title Del Permission
// @Description delete the User
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Param	uid		path 	uint64	true		"User ID"
// @Param	title	path 	string	true		"Title of Permission to Add to user"
// @Success 200 {object} []models.PermissionView
// @Failure 400 bad request (uid is missing or is not a number)
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:uid/permission/:title [delete]
func (c *UserController) DelPerm() {
	c.RequirePerm(models.PERM_ADMIN)
	var u = c.LoadUser()
	if dbErr := u.DelPermission(lib.GetDB(), c.GetString(":title")); dbErr != nil {
		c.AbortWith(400, dbErr)
	}

	c.Data["json"] = userView.AsJson(u.Permissions)
	c.ServeJSON()
}
