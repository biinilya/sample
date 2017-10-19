package controllers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"net/http"
	"net/http/httptest"
	"toptal/lib"
	"toptal/models"
	"toptal/tests"
	"toptal/tests/controllers"
)

func TestUserController(t *testing.T) {
	Convey("User Controller", t, func() {
		controllers.SetUp("/user", &UserController{})
		defer controllers.TearDown()

		Convey("Post Method", func() {
			Convey("When We Create User", func() {
				w := tests.BeegoCall("POST", "/api/v1/user/", nil)

				user := models.UserCredentialsView{}
				tests.FromJson(w.Body.Bytes(), &user)
				Convey("Status Code Should Be 201", func() {
					So(w.Code, ShouldEqual, 201)
				})
				Convey("It Should Be Found In DB", func() {
					var foundUser models.User
					loadErr := foundUser.LoadByKey(lib.GetDB(), user.Key)
					Convey("Without Errors", func() {
						So(loadErr, ShouldEqual, nil)
					})
					Convey("With Proper Key", func() {
						So(user.Key, ShouldEqual, foundUser.Key)
					})
					Convey("With Proper Secret", func() {
						So(lib.Encode(user.Secret), ShouldEqual, foundUser.Secret)
					})
					Convey("With User Permission Granted", func() {
						So(foundUser.HasPermission(models.PERM_USER), ShouldEqual, true)
					})
					Convey("Without Manager Permission Granted", func() {
						So(foundUser.HasPermission(models.PERM_MANAGER), ShouldEqual, false)
					})
					Convey("Without Admin Permission Granted", func() {
						So(foundUser.HasPermission(models.PERM_ADMIN), ShouldEqual, false)
					})
				})
			})
		})

		Convey("Put Credentials Method", func() {
			var user models.User
			user.New(lib.GetDB())
			var token = lib.GenerateAccessToken(uint64(user.Id))
			var headers = make(http.Header)
			headers.Set("X-Access-Token", token)
			var qUrl = fmt.Sprintf("/api/v1/user/%d/credentials", user.Id)

			var credentialsCheck = func(w *httptest.ResponseRecorder) func() {
				return func() {
					var res models.UserCredentialsView
					tests.FromJson(w.Body.Bytes(), &res)
					So(res.Id, ShouldEqual, user.Id)
					So(res.Key, ShouldNotBeBlank)
					So(res.Secret, ShouldNotBeBlank)
				}
			}

			Convey("When Calling Unauthenticated", func() {
				w := tests.BeegoCall("PUT", qUrl, nil)
				Convey("Status Code Should Be 401", func() {
					So(w.Code, ShouldEqual, 401)
				})
			})
			Convey("When Calling As A Regular User", func() {
				w := tests.BeegoCallWithHeader("PUT", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Return Proper Credentials", credentialsCheck(w))
			})
			Convey("When Calling As A Manager", func() {
				user.AddPermission(lib.GetDB(), models.PERM_MANAGER)
				w := tests.BeegoCallWithHeader("PUT", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Return Proper Credentials", credentialsCheck(w))
			})
			Convey("When Calling As An Admin", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				w := tests.BeegoCallWithHeader("PUT", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Return Proper Credentials", credentialsCheck(w))
			})
		})

		Convey("Delete Method", func() {
			var user models.User
			user.New(lib.GetDB())
			var token = lib.GenerateAccessToken(uint64(user.Id))
			var headers = make(http.Header)
			headers.Set("X-Access-Token", token)
			var qUrl = fmt.Sprintf("/api/v1/user/%d/", user.Id)

			var removedCheck = func(w *httptest.ResponseRecorder) func() {
				return func() {
					var user2 models.User
					var opErr = user2.LoadById(lib.GetDB(), uint64(user.Id))
					So(opErr, ShouldEqual, nil)
				}
			}

			Convey("When Calling Unauthenticated", func() {
				w := tests.BeegoCall("DELETE", qUrl, nil)
				Convey("Status Code Should Be 401", func() {
					So(w.Code, ShouldEqual, 401)
				})
			})
			Convey("When Calling As A Regular User", func() {
				w := tests.BeegoCallWithHeader("DELETE", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", removedCheck(w))
			})
			Convey("When Calling As A Manager", func() {
				user.AddPermission(lib.GetDB(), models.PERM_MANAGER)
				w := tests.BeegoCallWithHeader("DELETE", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", removedCheck(w))
			})
			Convey("When Calling As An Admin", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				w := tests.BeegoCallWithHeader("DELETE", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", removedCheck(w))
			})
		})

		Convey("GetAll Method", func() {
			var user models.User
			user.New(lib.GetDB())
			var token = lib.GenerateAccessToken(uint64(user.Id))
			var headers = make(http.Header)
			headers.Set("X-Access-Token", token)

			Convey("When Calling Unauthenticated", func() {
				w := tests.BeegoCall("GET", "/api/v1/user/", nil)
				Convey("Status Code Should Be 401", func() {
					So(w.Code, ShouldEqual, 401)
				})
			})
			Convey("When Calling As A Regular User", func() {
				w := tests.BeegoCallWithHeader("GET", "/api/v1/user/", nil, headers)
				Convey("Status Code Should Be 403", func() {
					So(w.Code, ShouldEqual, 403)
				})
			})
			Convey("When Calling As A Manager", func() {
				user.AddPermission(lib.GetDB(), models.PERM_MANAGER)
				w := tests.BeegoCallWithHeader("GET", "/api/v1/user/", nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
			})
			Convey("When Calling As An Admin", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				w := tests.BeegoCallWithHeader("GET", "/api/v1/user/", nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("The Result Should Be A List Of Prepared Users", func() {
					var res = []models.UserInfoView{}
					tests.FromJson(w.Body.Bytes(), &res)
					So(len(res), ShouldEqual, 1)
					So(res[0].Key, ShouldEqual, user.Key)
					So(res[0].Id, ShouldEqual, user.Id)
				})
			})
		})
	})
}
