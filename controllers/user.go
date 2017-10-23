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
		c.Abort("500")
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title SignIn
// @Description Use this method to receive Access Token which is required for API access
// @Param	body		body 	models.UserCredentialsView	true		"Credentials to Sign In"
// @Success 200 {object} models.UserAccessTokenView
// @Failure 400 bad request
// @Failure 401 wrong credentials
// @router /sign_in [post]
func (c *UserController) SignIn() {
	var cred models.UserCredentialsView
	if dataInErr := json.Unmarshal(c.Ctx.Input.RequestBody, &cred); dataInErr != nil {
		c.Abort("400")
	}
	if cred.Key == "" || cred.Secret == "" {
		c.Abort("400")
	}

	var v models.User
	if opErr := v.LoadByKey(lib.GetDB(), cred.Key); opErr != nil {
		c.Abort("401")
	}
	if lib.Encode(cred.Secret) != v.Secret {
		c.Abort("401")
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
	var uid, uidErr = c.GetUint64(":uid")
	if uidErr != nil {
		c.Abort("400")
	}
	var u = c.RequireOwnerOrPerm(uid, models.PERM_MANAGER, models.PERM_ADMIN)

	var secret, dbErr = u.NewCred(lib.GetDB())
	if dbErr != nil {
		c.Abort("500")
	}

	c.Data["json"] = models.UserCredentialsView{
		Id:     uint64(u.Id),
		Key:    u.Key,
		Secret: secret,
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description Users Directory
// @Param	X-Access-Token		header 	string	true		"Access Token"
// @Success 200 {object} []models.UserInfoView
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router / [get]
func (c *UserController) GetAll() {
	c.RequirePerm(models.PERM_MANAGER, models.PERM_ADMIN)

	var users []*models.User
	var opErr error
	if users, opErr = models.UsersGetAll(lib.GetDB()); opErr != nil {
		c.Abort("500")
	}

	c.Data["json"] = userView.AsJson(users)
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
	var uid, uidErr = c.GetUint64(":uid")
	if uidErr != nil {
		c.Abort("400")
	}
	var u = c.RequireOwnerOrPerm(uid, models.PERM_MANAGER, models.PERM_ADMIN)

	var dbErr = u.Delete(lib.GetDB())
	if dbErr != nil {
		c.Abort("500")
	}
	c.Data["json"] = userView.AsJson(u)
	c.ServeJSON()
}