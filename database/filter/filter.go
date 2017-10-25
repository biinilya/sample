package filter

import (
	"sync"

	"regexp"

	"errors"

	"strconv"

	"fmt"

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
		return "", nil, formatParsingError(rawQ, parseErr)
	}
	f.filter.Execute()
	return f.filter.topExpr().renderToSql(fields...)
}

func (f *filter) ParseToOrm(rawQ string, fields ...string) (*orm.Condition, error) {
	f.init()

	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.reset(rawQ)
	if parseErr := f.filter.Parse(); parseErr != nil {
		return nil, formatParsingError(rawQ, parseErr)
	}
	f.filter.Execute()
	return f.filter.topExpr().renderToOrm(fields...)
}

func formatParsingError(rawQ string, parseErr error) error {
	var msg = parseErr.Error()
	var sq = fmtRe.FindStringSubmatch(msg)
	if len(sq) < 3 {
		return parseErr
	}
	var errPoint1, _ = strconv.ParseInt(sq[1], 10, 64)
	var errPoint2, _ = strconv.ParseInt(sq[2], 10, 64)
	var errStr = fmt.Sprintf("%s ||%s|| %s", rawQ[0:errPoint1], rawQ[errPoint1:errPoint2], rawQ[errPoint2:])
	return errors.New(errStr)
}

var fmtRe = regexp.MustCompile("symbol (\\d+).*symbol (\\d+)")

var Filter filter
