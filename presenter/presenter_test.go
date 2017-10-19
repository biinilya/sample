package presenter

import (
	"testing"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonSerializerEncode(t *testing.T) {
	var ser = New().AddHook(A1{}, new(structView).TranslateCase().Include("A"))

	var testMap = []struct {
		Scenario string
		Sample   interface{}
		Result   string
	}{
		{
			"When simple struct serialize",
			A1{"test1", -12.5, 13},
			`{"a":"test1"}`,
		},
		{
			"When nested struct serialize",
			map[string]A1{"xxx": {"test1", -12.5, 13}},
			`{"xxx":{"a":"test1"}}`,
		},
		{
			"When nested pointer struct serialize",
			map[string]*A1{"xxx": {"test1", -12.5, 13}},
			`{"xxx":{"a":"test1"}}`,
		},
	}

	for _, test := range testMap {
		var scenario, sample, result = test.Scenario, test.Sample, test.Result

		var data, dataErr = json.Marshal(ser.AsJson(sample))
		Convey(scenario+"no errors occured", t, func() {
			So(dataErr, ShouldEqual, nil)
		})
		Convey(scenario, t, func() {
			So(string(data), ShouldEqual, result)
		})
	}
}

func TestJsonSerializerDecode(t *testing.T) {
	var ser = New().AddHook(A1{}, new(structView).TranslateCase().Include("A"))

	var test = func(scenario string, js string, emptyData interface{}, sample interface{}) {
		var dataErr = json.Unmarshal([]byte(js), ser.AsJson(emptyData))
		Convey(scenario+"no errors occured", t, func() {
			So(dataErr, ShouldEqual, nil)
		})
		Convey(scenario, t, func() {
			So(emptyData, ShouldResemble, sample)
		})
	}

	test(
		"When simple struct deserialize",
		`{"a":"test1"}`,
		&A1{},
		&A1{A: "test1"},
	)

	test(
		"When nested struct deserialize",
		`{"xxx":{"a":"test1"}}`,
		map[string]A1{},
		map[string]A1{"xxx": {A: "test1"}},
	)

	test(
		"When nested pointer struct deserialize",
		`{"xxx":{"a":"test1"}}`,
		map[string]*A1{},
		map[string]*A1{"xxx": {A: "test1"}},
	)
}
