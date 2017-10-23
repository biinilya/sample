package models

import (
	"time"

	"toptal/lib"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id          int64         `orm:"auto"`
	Key         string        `orm:"size(128);index"`
	Secret      string        `orm:"size(128)"`
	Created     time.Time     `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time     `orm:"auto_now;type(datetime)"`
	Permissions []*Permission `orm:"rel(m2m)"`
	Records     []*Record     `orm:"reverse(many)"`
}

func (u *User) AddPermission(o orm.Ormer, permTitle string) (err error) {
	var perm *Permission
	if perm, err = PermissionByTitle(o, permTitle); err != nil {
		return
	}
	for _, uPerm := range u.Permissions {
		if perm.Id == uPerm.Id {
			return
		}
	}

	if _, err = o.QueryM2M(u, "Permissions").Add(perm); err == nil {
		u.Permissions = append(u.Permissions, perm)
	}
	return
}

func (u *User) HasPermission(permCode string) bool {
	for _, uPerm := range u.Permissions {
		if permCode == uPerm.Title {
			return true
		}
	}
	return false
}

func (u *User) New(o orm.Ormer, extraPerm ...string) (secret string, opErr error) {
	secret = lib.RndKey()
	u.Key = lib.RndKey()
	u.Secret = lib.Encode(secret)

	if u.Id, opErr = o.Insert(u); opErr != nil {
		return
	}
	for _, permTitle := range append(extraPerm, PERM_USER) {
		if opErr = u.AddPermission(o, permTitle); opErr != nil {
			return
		}
	}

	return
}

func (u *User) NewCred(o orm.Ormer) (secret string, opErr error) {
	secret = lib.RndKey()
	u.Key = lib.RndKey()
	u.Secret = lib.Encode(secret)

	if _, opErr = o.Update(u, "Key", "Secret"); opErr != nil {
		return
	}

	return
}

func (u *User) LoadByKey(o orm.Ormer, key string) (opErr error) {
	if opErr = o.QueryTable(u).Filter("Key", key).One(u); opErr != nil {
		return
	}
	_, opErr = o.LoadRelated(u, "Permissions")
	return
}

func (u *User) LoadById(o orm.Ormer, id uint64) (opErr error) {
	if opErr = o.QueryTable(u).Filter("Id", id).One(u); opErr != nil {
		return
	}
	_, opErr = o.LoadRelated(u, "Permissions")
	return
}

func (u *User) Delete(o orm.Ormer) (opErr error) {
	if _, opErr = o.QueryM2M(u, "Permissions").Clear(); opErr == nil {
		u.Permissions = []*Permission{}
	}
	if _, opErr = o.QueryTable(u).Filter("Id", u.Id).Delete(); opErr != nil {
		return
	}
	return
}

func UsersGetAll(o orm.Ormer) (users []*User, opErr error) {
	if _, opErr = o.QueryTable(User{}).All(&users); opErr != nil {
		return
	}
	return
}

type UserCredentialsView struct {
	Id     uint64 `json:"id"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type UserAccessTokenView struct {
	AccessToken string `json:"access-token"`
}

type UserInfoView struct {
	Id      int64     `json:"id"`
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func init() {
	orm.RegisterModel(new(User))
}
