package presenter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ugorji/go/codec"
)

func TestJsonSerializerEncode(t *testing.T) {
	var ser = Presenter{}
	var a1Handler = new(StructView).TranslateCase().Include("A")
	ser.AddHook(A1{}, a1Handler)

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

		var data []byte
		var dataErr = codec.NewEncoderBytes(&data, &ser.JsonHandle).Encode(sample)
		Convey(scenario+"no errors occured", t, func() {
			So(dataErr, ShouldEqual, nil)
		})
		Convey(scenario, t, func() {
			So(string(data), ShouldEqual, result)
		})
	}
}

func TestJsonSerializerDecode(t *testing.T) {
	var ser = Presenter{}
	var a1Handler = new(StructView).TranslateCase().Include("A")
	ser.AddHook(A1{}, a1Handler)

	var test = func(scenario string, js string, emptyData interface{}, sample interface{}) {
		var dataErr = codec.NewDecoderBytes([]byte(js), &ser.JsonHandle).Decode(emptyData)
		Convey(scenario+"no errors occured", t, func() {
			So(dataErr, ShouldEqual, nil)
		})
		Convey(scenario, t, func() {
			So(emptyData, ShouldResemble, sample)
		})
	}

	test(
		"When simple struct serialize",
		`{"a":"test1"}`,
		&A1{},
		&A1{A: "test1"},
	)

	test(
		"When nested struct serialize",
		`{"xxx":{"a":"test1"}}`,
		map[string]A1{},
		map[string]A1{"xxx": {A: "test1"}},
	)

	test(
		"When nested pointer struct serialize",
		`{"xxx":{"a":"test1"}}`,
		map[string]*A1{},
		map[string]*A1{"xxx": {A: "test1"}},
	)
}
