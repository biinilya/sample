package main

import (
	"sync"

	"github.com/astaxie/beego/orm"
)

type filter struct {
	once   sync.Once
	mutex  sync.Mutex
	filter FilterPeg
}

func (f *filter) init() {
	f.once.Do(func() {
		f.filter.Init()
	})
}

func (f *filter) reset(rawQ string) {
	f.filter.Buffer = rawQ
	f.filter.Reset()
	f.filter.AST.Init()
}

func (f *filter) ParseToSql(rawQ string, fields ...string) (string, []interface{}, error) {
	f.init()

	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.reset(rawQ)
	if parseErr := f.filter.Parse(); parseErr != nil {
		return "", nil, parseErr
	}

	return f.filter.topExpr().renderToSql(fields...)
}

func (f *filter) ParseToOrm(rawQ string, fields ...string) (*orm.Condition, error) {
	f.init()

	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.reset(rawQ)
	if parseErr := f.filter.Parse(); parseErr != nil {
		return nil, parseErr
	}

	return f.filter.topExpr().renderToOrm(fields...)
}

var Filter filter
