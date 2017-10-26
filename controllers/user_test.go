package controllers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"bytes"
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

			Convey("When Calling With Broken uid", func() {
				w := tests.BeegoCallWithHeader("PUT", "/api/v1/user/xxx/credentials", nil, headers)
				Convey("Status Code Should Be 400", func() {
					So(w.Code, ShouldEqual, 400)
				})
			})
			Convey("When Calling With Broken Non-Existing uid", func() {
				qUrl = fmt.Sprintf("/api/v1/user/%d/credentials", user.Id-1)
				w := tests.BeegoCallWithHeader("PUT", qUrl, nil, headers)
				Convey("Status Code Should Be 403", func() {
					So(w.Code, ShouldEqual, 403)
				})
			})

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
					user2.LoadById(lib.GetDB(), uint64(user.Id))
					So(user2.Id, ShouldEqual, 0)
				}
			}

			Convey("When Calling With Broken uid", func() {
				w := tests.BeegoCallWithHeader("DELETE", "/api/v1/user/xxx/", nil, headers)
				Convey("Status Code Should Be 400", func() {
					So(w.Code, ShouldEqual, 400)
				})
			})
			Convey("When Calling With Broken Non-Existing uid", func() {
				w := tests.BeegoCallWithHeader("DELETE", "/api/v1/user/0/", nil, headers)
				Convey("Status Code Should Be 403", func() {
					So(w.Code, ShouldEqual, 403)
				})
			})
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

		Convey("SignIn Method", func() {
			var user models.User
			var secret, _ = user.New(lib.GetDB())
			var credentials = models.UserCredentialsView{Key: user.Key, Secret: secret}

			Convey("When Invalid Body", func() {
				w := tests.BeegoCall("POST", "/api/v1/user/sign_in", nil)
				Convey("Status Code Should Be 400", func() {
					So(w.Code, ShouldEqual, 400)
				})
			})
			Convey("When Empty Body", func() {
				w := tests.BeegoCall("POST", "/api/v1/user/sign_in", bytes.NewReader(
					tests.ToJson(models.UserCredentialsView{})),
				)
				Convey("Status Code Should Be 400", func() {
					So(w.Code, ShouldEqual, 400)
				})
			})
			Convey("When Invalid Key", func() {
				credentials.Key += "SomeRandomString"
				w := tests.BeegoCall("POST", "/api/v1/user/sign_in", bytes.NewReader(
					tests.ToJson(credentials)),
				)
				Convey("Status Code Should Be 403", func() {
					So(w.Code, ShouldEqual, 403)
				})
			})
			Convey("When Invalid Secret", func() {
				credentials.Secret += "SomeRandomString"
				w := tests.BeegoCall("POST", "/api/v1/user/sign_in", bytes.NewReader(
					tests.ToJson(credentials)),
				)
				Convey("Status Code Should Be 403", func() {
					So(w.Code, ShouldEqual, 403)
				})
			})
			Convey("When Credentials Are Valid", func() {
				w := tests.BeegoCall("POST", "/api/v1/user/sign_in", bytes.NewReader(
					tests.ToJson(credentials)),
				)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("AccessToken Should Be Valid", func() {
					var at models.UserAccessTokenView
					tests.FromJson(w.Body.Bytes(), &at)
					Convey("With Validation Passed", func() {
						var uid, decOk = lib.ValidateAccessToken(at.AccessToken)
						So(decOk, ShouldBeTrue)
						So(uid, ShouldEqual, user.Id)
						Convey("With Valid User", func() {
							So(uid, ShouldEqual, user.Id)
						})
					})
				})
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
			Convey("When Applying Filter And User Found", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				qUrl := fmt.Sprintf(`/api/v1/user/?filter=(key eq '%s')`, user.Key)
				w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				if w.Code == 200 {
					Convey("The Result Should Be A List Of Prepared Users", func() {
						var res = []models.UserInfoView{}
						tests.FromJson(w.Body.Bytes(), &res)
						So(len(res), ShouldEqual, 1)
						So(res[0].Key, ShouldEqual, user.Key)
						So(res[0].Id, ShouldEqual, user.Id)
					})
				}
			})
			Convey("When Applying Filter And User Not Found", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				qUrl := fmt.Sprintf(`/api/v1/user/?filter=(key ne '%s')`, user.Key)
				w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				if w.Code == 200 {
					Convey("The Result Should Be Empty", func() {
						var res = []models.UserInfoView{}
						tests.FromJson(w.Body.Bytes(), &res)
						So(len(res), ShouldEqual, 0)
					})
				}
			})
			Convey("When Applying Filter And Filter Is Broken", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				Convey("When Invalid Key In Filter", func() {
					qUrl := fmt.Sprintf(`/api/v1/user/?filter=(keyxx = '%sabc')`, user.Key)
					w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
					Convey("Status Code Should Be 400", func() {
						So(w.Code, ShouldEqual, 400)
					})
				})
				Convey("When Invalid Syntax In Filter", func() {
					qUrl := fmt.Sprintf(`/api/v1/user/?filter=(key == '%sabc')`, user.Key)
					w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
					Convey("Status Code Should Be 400", func() {
						So(w.Code, ShouldEqual, 400)
					})
				})
			})
		})

		Convey("Get One Method", func() {
			var user models.User
			user.New(lib.GetDB())
			var token = lib.GenerateAccessToken(uint64(user.Id))
			var headers = make(http.Header)
			headers.Set("X-Access-Token", token)
			var qUrl = fmt.Sprintf("/api/v1/user/%d/", user.Id)

			var infoCheck = func(w *httptest.ResponseRecorder) func() {
				return func() {
					var uv models.UserInfoView
					tests.FromJson(w.Body.Bytes(), &uv)
					So(uv.Id, ShouldEqual, user.Id)
					So(uv.Key, ShouldEqual, user.Key)
					So(uv.Created, ShouldNotBeEmpty)
					So(uv.Updated, ShouldNotBeEmpty)
				}
			}

			Convey("When Calling Unauthenticated", func() {
				w := tests.BeegoCall("GET", qUrl, nil)
				Convey("Status Code Should Be 401", func() {
					So(w.Code, ShouldEqual, 401)
				})
			})
			Convey("When Calling As A Regular User", func() {
				w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", infoCheck(w))
			})
			Convey("When Calling As A Manager", func() {
				user.AddPermission(lib.GetDB(), models.PERM_MANAGER)
				w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", infoCheck(w))
			})
			Convey("When Calling As An Admin", func() {
				user.AddPermission(lib.GetDB(), models.PERM_ADMIN)
				w := tests.BeegoCallWithHeader("GET", qUrl, nil, headers)
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey("Should Remove The User", infoCheck(w))
			})
		})
	})
}
