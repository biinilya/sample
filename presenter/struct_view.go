package presenter

import (
	"reflect"

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
		var inVal = reflect.ValueOf(in_val_iface)
		key = as.snakeToCamel(key)
		if as.exclude != nil && as.exclude[key] {
			continue
		}
		if as.include == nil || as.include[key] {
			if _, ok := valType.FieldByName(key); ok {
				val.FieldByName(key).Set(inVal)
			}
		}
	}
}

func StructView() *structView {
	return &structView{}
}
