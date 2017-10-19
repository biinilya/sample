package presenter

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

type Presenter struct {
	codec.JsonHandle
	tag uint64
}

func (js *Presenter) AddHook(item interface{}, handler *StructView) {
	js.JsonHandle.SetInterfaceExt(
		reflect.TypeOf(item),
		js.tag,
		interfaceExt{handler},
	)
	js.tag++
}
