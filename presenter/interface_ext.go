package presenter

import (
	"reflect"
)

type interfaceExt struct {
	*structView
}

// ConvertExt converts a value into a simpler interface for easy encoding e.g. convert time.Time to int64.
//
// Note: v *may* be a pointer to the extension type, if the extension type was a struct or array.
func (ie interfaceExt) ConvertExt(v interface{}) interface{} {
	var vType = reflect.TypeOf(v)
	if vType.Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem().Interface()
	}
	return ie.structView.Struct2Map(v)
}

// UpdateExt updates a value from a simpler interface for easy decoding e.g. convert int64 to time.Time.
func (ie interfaceExt) UpdateExt(dst interface{}, src interface{}) {
	var srcMap = src.(map[interface{}]interface{})
	var srcStrMap = make(map[string]interface{}, len(srcMap))
	for k, v := range srcMap {
		srcStrMap[k.(string)] = v
	}
	ie.structView.Map2Struct(srcStrMap, dst)
}
