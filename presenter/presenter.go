package presenter

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

type jsonView struct {
	presenter *presenter
	item      interface{}
}

func (v *jsonView) MarshalJSON() (data []byte, dataErr error) {
	data = []byte{}
	dataErr = codec.NewEncoderBytes(&data, &v.presenter.handler).Encode(v.item)
	return
}

func (v *jsonView) UnmarshalJSON(data []byte) (dataErr error) {
	dataErr = codec.NewDecoderBytes(data, &v.presenter.handler).Decode(v.item)
	return
}

type presenter struct {
	handler codec.JsonHandle
	tag     uint64
}

func (js *presenter) AddHook(item interface{}, handler *structView) *presenter {
	js.handler.SetInterfaceExt(
		reflect.TypeOf(item),
		js.tag,
		interfaceExt{handler},
	)
	js.tag++
	return js
}

func (js *presenter) AsJson(item interface{}) *jsonView {
	return &jsonView{presenter: js, item: item}
}

func New() *presenter {
	return &presenter{}
}
