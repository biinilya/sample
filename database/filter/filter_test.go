package filter

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	filter := &FilterPeg{}
	filter.Init()

	var testExpression = func(expression string, qString string, qArgs ...interface{}) {
		Convey(fmt.Sprintf("%s -> %s", expression, qString), func() {
			filter.Buffer = expression
			filter.Reset()
			filter.AST.Init()

			filterErr := filter.Parse()
			Convey("Parsed without errors", func() {
				So(filterErr, ShouldBeNil)
			})
			filter.Execute()
			var qS, qA, _ = filter.Render()
			Convey("Result Query String Is Valid", func() {
				So(qS, ShouldEqual, qString)
			})
			Convey("Result Query Args Has The Same Length", func() {
				So(len(qA), ShouldResemble, len(qArgs))
			})
			if len(qA) == len(qArgs) {
				for idx := range qA {
					Convey(fmt.Sprintf("Result Query Args #%d Is Valid", idx), func() {
						So(qA[idx], ShouldEqual, qArgs[idx])
					})
				}
			}
		})
	}

	var testForbiddenFieldsExpression = func(expression string, errMsg string) {
		Convey(fmt.Sprintf("[Restricted] %s", expression), func() {
			filter.Buffer = expression
			filter.Reset()
			filter.AST.Init()
			filterErr := filter.Parse()
			Convey("Parsed without errors", func() {
				So(filterErr, ShouldBeNil)
			})
			filter.Execute()
			var _, _, qErr = filter.Render("created_at", "value", "death")
			Convey("Result Query String Is Valid", func() {
				So(qErr, ShouldNotBeNil)
			})
			if qErr != nil {
				Convey("Error message in proper format", func() {
					So(qErr.Error(), ShouldEqual, errMsg)
				})
			}
		})
	}

	var testBrokenExpression = func(expression string, errMsg string) {
		Convey(fmt.Sprintf("[Broken] expr %s", expression), func() {
			filter.Buffer = expression
			filter.Reset()
			filterErr := filter.Parse()
			Convey("Parsed with errors", func() {
				So(filterErr, ShouldNotBeNil)
			})
			if filterErr != nil {
				Convey("Error message in proper format", func() {
					So(formatParsingError(expression, filterErr).Error(), ShouldEqual, errMsg)
				})
			}
		})
	}

	Convey("Our calculator should work", t, func() {
		testExpression(
			`created_at>'2017-01-01'`,
			`"created_at" > ?`,
			`2017-01-01`,
		)
		testExpression(
			`(created_at > '2017-01-01') AND (value < 19)`,
			`("created_at" > ?) AND ("value" < ?)`,
			`2017-01-01`, 19,
		)
		testExpression(
			`(key = '032ae64df1ab4c4dafc139f6c26cb647')`,
			`"key" = ?`,
			`032ae64df1ab4c4dafc139f6c26cb647`,
		)
		testExpression(
			`(created_at > '2017-01-01' AND value < 19) OR death = 1`,
			`(("created_at" > ?) AND ("value" < ?)) OR ("death" = ?)`,
			`2017-01-01`, 19, 1,
		)
		testExpression(
			`created_at > '2017-01-01' AND value < 19 OR death = 1`,
			`(("created_at" > ?) AND ("value" < ?)) OR ("death" = ?)`,
			`2017-01-01`, 19, 1,
		)
		testExpression(
			`(created_at > '2017-01-01' AND value < 19) OR NOT death = 1`,
			`(("created_at" > ?) AND ("value" < ?)) OR (NOT ("death" = ?))`,
			`2017-01-01`, 19, 1,
		)
		testBrokenExpression(
			`created_at >> '2017-01-01'`,
			`created_at > ||>||  '2017-01-01'`,
		)
		testBrokenExpression(
			`(created_at > '2017-01-01') AND (value < 19)x`,
			`(created_at > '2017-01-01') AND (value < 19) ||x|| `,
		)
		testForbiddenFieldsExpression(
			`created > '2017-01-01'`,
			`field 'created' is not allowed to filter in`,
		)
	})
}
