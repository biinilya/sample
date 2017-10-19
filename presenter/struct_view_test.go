package presenter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type A1 struct {
	A    string
	AAA  float64
	AaBB int
}

func TestAutoSerializerStruct2Map(t *testing.T) {
	var testMap = []struct {
		Scenario   string
		Serializer *structView
		Sample     A1
		Result     map[string]interface{}
	}{
		{
			"When default configured",
			new(structView),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"A": "A1test", "AAA": -12.25, "AaBB": 12},
		},
		{
			"When snake case requested",
			new(structView).TranslateCase(),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"a": "A1test", "a_a_a": -12.25, "aa_b_b": 12},
		},
		{
			"When excluded fields exist",
			new(structView).Exclude("A", "AAA"),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"AaBB": 12},
		},
		{
			"When included fields exist",
			new(structView).Include("A", "AAA"),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"A": "A1test", "AAA": -12.25},
		},
		{
			"When both included and excluded fields exist",
			new(structView).Include("A", "AAA").Exclude("A"),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"AAA": -12.25},
		},
	}

	for _, test := range testMap {
		var scenario, serializer, sample, result = test.Scenario, test.Serializer, test.Sample, test.Result
		Convey(scenario, t, func() {
			So(serializer.Struct2Map(sample), ShouldResemble, result)
		})
	}
}

func TestAutoSerializerMap2Struct(t *testing.T) {
	var testMap = []struct {
		Scenario   string
		Serializer *structView
		Sample     A1
		Result     map[string]interface{}
	}{
		{
			"When default configured",
			new(structView),
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"A": "A1test", "AAA": -12.25, "AaBB": 12},
		},
		{
			"When snake case requested",
			&structView{translateCase: true},
			A1{"A1test", -12.25, 12},
			map[string]interface{}{"a": "A1test", "a_a_a": -12.25, "aa_b_b": 12},
		},
		{
			"When excluded fields exist",
			new(structView).Exclude("A", "AAA"),
			A1{"", 0, 12},
			map[string]interface{}{"A": "A1test", "AAA": -12.25, "AaBB": 12},
		},
		{
			"When included fields exist",
			new(structView).Include("A", "AAA"),
			A1{"A1test", -12.25, 0},
			map[string]interface{}{"A": "A1test", "AAA": -12.25, "AaBB": 12},
		},
		{
			"When both included and excluded fields exist",
			new(structView).Include("A", "AAA").Exclude("A"),
			A1{"", -12.25, 0},
			map[string]interface{}{"A": "A1test", "AAA": -12.25, "AaBB": 12},
		},
	}

	for _, test := range testMap {
		var scenario, serializer, sample, result = test.Scenario, test.Serializer, test.Sample, test.Result
		Convey(scenario, t, func() {
			var sampleRes A1
			serializer.Map2Struct(result, &sampleRes)
			So(sample, ShouldResemble, sampleRes)
		})
	}
}
