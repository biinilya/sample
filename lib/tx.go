package lib

import "github.com/astaxie/beego/orm"

func WithTx(o orm.Ormer, handlers ...func(orm.Ormer) error) (opErr error) {
	if opErr = o.Begin(); opErr != nil {
		return
	}

	for _, handler := range handlers {
		if opErr = handler(o); opErr != nil {
			o.Rollback()
			return
		}
	}

	opErr = o.Commit()
	return
}
