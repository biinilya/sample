package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	filter := &Filter{}
	filter.Init()
	filter.Expression.Init()

	var testExpression = func(expression string, qString string, qArgs ...interface{}) {
		Convey(fmt.Sprintf("%s", expression), func() {
			filter.Buffer = expression
			filter.Reset()
			filterErr := filter.Parse()
			Convey("Parsed without errors", func() {
				So(filterErr, ShouldBeNil)
			})
			filter.Execute()
			filter.Render()

			var qS, qA = filter.Render()
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

	Convey("Our calculator should work", t, func() {
		testExpression(
			"created_at > '2017-01-01'",
			`"created_at" > ?`,
			"2017-01-01",
		)
		testExpression(
			"(created_at > '2017-01-01') AND (value < 19)",
			`("created_at" > ?) AND ("value" < ?)`,
			"2017-01-01", 19,
		)
		testExpression(
			"(created_at > '2017-01-01' AND value < 19) OR death = 1",
			`(("created_at" > ?) AND ("value" < ?)) OR ("death" = ?)`,
			"2017-01-01", 19, 1,
		)
		testExpression(
			"(created_at > '2017-01-01' AND value < 19) OR NOT death = 1",
			`(("created_at" > ?) AND ("value" < ?)) OR (NOT ("death" = ?))`,
			"2017-01-01", 19, 1,
		)
	})
}
