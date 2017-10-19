package lib

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	_ "toptal/tests"
)

func TestAccessTokenFunc(t *testing.T) {
	Convey("With AccessToken Library", t, func() {
		var uid = uint64(rand.Int63())
		var key = GenerateAccessToken(uid)
		Convey("When Access Token Is Generated", func() {
			Convey("It Is Not Empty", func() {
				So(key, ShouldNotBeBlank)
			})
			var key2 = GenerateAccessToken(uid)
			Convey("It Is Not The Same When Generated Again", func() {
				So(key2, ShouldNotEqual, key)
			})
			uid_decoded, valid := ValidateAccessToken(key)
			Convey("It Can Be Successfully Validated", func() {
				So(valid, ShouldBeTrue)
			})
			Convey("It Matches The Same User It Was Generated For", func() {
				So(uid_decoded, ShouldEqual, uid)
			})
		})
		Convey("When Access Token Is Wrong", func() {
			_, valid := ValidateAccessToken("123123123")
			Convey("It Can Not Be Successfully Validated", func() {
				So(valid, ShouldBeFalse)
			})
		})
	})
}
