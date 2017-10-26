package presenter

import (
	"reflect"

	"time"

	"fmt"

	"errors"

	"github.com/serenize/snaker"
)

type structView struct {
	translateCase bool

	include map[string]bool
	exclude map[string]bool
}

func (as *structView) TranslateCase() *structView {
	as.translateCase = true
	return as
}

func (as *structView) Include(list ...string) *structView {
	as.include = make(map[string]bool)
	for _, l := range list {
		as.include[l] = true
	}
	return as
}

func (as *structView) Exclude(list ...string) *structView {
	as.exclude = make(map[string]bool)
	for _, l := range list {
		as.exclude[l] = true
	}
	return as
}

func (as *structView) camelToSnake(key string) string {
	if !as.translateCase {
		return key
	}
	return snaker.CamelToSnake(key)
}

func (as *structView) snakeToCamel(key string) string {
	if !as.translateCase {
		return key
	}
	return snaker.SnakeToCamel(key)
}

func (as *structView) Struct2Map(model interface{}) map[string]interface{} {
	val := reflect.ValueOf(model)

	m := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if as.exclude != nil && as.exclude[field.Name] {
			continue
		}
		if as.include == nil || as.include[field.Name] {
			m[as.camelToSnake(field.Name)] = val.Field(i).Interface()
		}
	}
	return m
}

func (as *structView) Map2Struct(in map[string]interface{}, model interface{}) {
	val := reflect.ValueOf(model).Elem()
	valType := reflect.TypeOf(model).Elem()
	for key, in_val_iface := range in {
		key = as.snakeToCamel(key)
		if as.exclude != nil && as.exclude[key] {
			continue
		}
		if as.include == nil || as.include[key] {
			if field, ok := valType.FieldByName(key); ok {
				val.FieldByName(key).Set(convertValue(in_val_iface, field.Type))
			}
		}
	}
}

func convertValue(in interface{}, out reflect.Type) reflect.Value {
	var inType = reflect.TypeOf(in)
	var inValue = reflect.ValueOf(in)
	if inType == out {
		return inValue
	}
	switch {
	case out == reflect.TypeOf(time.Time{}):
		switch inType.Kind() {
		case reflect.String:
			for _, format := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"} {
				if d, dErr := time.Parse(format, in.(string)); dErr == nil {
					return reflect.ValueOf(d)
				}
			}
			panic(errors.New(fmt.Sprintf("Cannot parse '%s' as 'time.Time'", in.(string))))
		}
	case out.Kind() == reflect.Float32:
		switch {
		case inType.Kind() == reflect.Int ||
			inType.Kind() == reflect.Int8 ||
			inType.Kind() == reflect.Int16 ||
			inType.Kind() == reflect.Int32 ||
			inType.Kind() == reflect.Int64:
			return reflect.ValueOf(float32(inValue.Int()))
		}
	case out.Kind() == reflect.Float64:
		switch {
		case inType.Kind() == reflect.Int ||
			inType.Kind() == reflect.Int8 ||
			inType.Kind() == reflect.Int16 ||
			inType.Kind() == reflect.Int32 ||
			inType.Kind() == reflect.Int64:
			return reflect.ValueOf(float64(inValue.Int()))
		}
	}
	return inValue
}

func StructView() *structView {
	return &structView{}
}
