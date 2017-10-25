package filter

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/orm"
)

func (c Comprator) toOrmComparator() string {
	switch c {
	case CmpEq:
		return ""
	case CmpGt:
		return "__gt"
	case CmpGte:
		return "__gt"
	case CmpLt:
		return "__lt"
	case CmpLte:
		return "__lt"
	default:
		panic("Unknown comparator")
	}
}

type expr struct {
	idenitifier string
	comparator  Comprator
	args        []interface{}
	agg         *aggregation
}

type aggregation struct {
	exprList []expr
	operator Operator
}

func (e expr) renderToSql(fields ...string) (string, []interface{}, error) {
	var fieldsAllowed = map[string]bool{}
	for _, f := range fields {
		fieldsAllowed[f] = true
	}
	if e.agg == nil {
		if !fieldsAllowed[e.idenitifier] && len(fields) > 0 {
			return "", nil, errors.New(fmt.Sprintf(`field <%s> is not allowd to filter in`, e.idenitifier))
		}
		return fmt.Sprintf(`"%s" %s ?`, e.idenitifier, e.comparator), e.args, nil
	}
	if e.agg.operator.ArgsNum() == 1 {
		var exprQ, exprArgs, exprErr = e.agg.exprList[0].renderToSql()
		if exprErr != nil {
			return exprQ, exprArgs, exprErr
		}
		return fmt.Sprintf(`%s (%s)`, e.agg.operator, exprQ), exprArgs, nil
	}

	if e.agg.operator.ArgsNum() == 2 {
		var exprQ1, exprArgs1, exprErr1 = e.agg.exprList[0].renderToSql()
		if exprErr1 != nil {
			return exprQ1, exprArgs1, exprErr1
		}
		var exprQ2, exprArgs2, exprErr2 = e.agg.exprList[1].renderToSql()
		if exprErr2 != nil {
			return exprQ2, exprArgs2, exprErr2
		}
		return fmt.Sprintf(`(%s) %s (%s)`, exprQ2, e.agg.operator, exprQ1), append(exprArgs2, exprArgs1...), nil
	}

	return "", nil, errors.New("Invalid operator ")
}

func (e expr) renderToOrm(fields ...string) (*orm.Condition, error) {
	var fieldsAllowed = map[string]bool{}
	for _, f := range fields {
		fieldsAllowed[f] = true
	}
	if e.agg == nil {
		if !fieldsAllowed[e.idenitifier] && len(fields) > 0 {
			return nil, errors.New(fmt.Sprintf(`field <%s> is not allowd to filter in`, e.idenitifier))
		}
		return new(orm.Condition).And(
			fmt.Sprintf(`%s%s`, e.idenitifier, e.comparator.toOrmComparator()),
			e.args...,
		), nil
	}
	if e.agg.operator.ArgsNum() == 1 {
		var exprQ, exprErr = e.agg.exprList[0].renderToOrm()
		if exprErr != nil {
			return exprQ, exprErr
		}
		return new(orm.Condition).AndNotCond(exprQ), nil
	}

	if e.agg.operator.ArgsNum() == 2 {
		var exprQ1, exprErr1 = e.agg.exprList[0].renderToOrm()
		if exprErr1 != nil {
			return exprQ1, exprErr1
		}
		var exprQ2, exprErr2 = e.agg.exprList[1].renderToOrm()
		if exprErr2 != nil {
			return exprQ2, exprErr2
		}
		switch e.agg.operator {
		case OpAnd:
			return exprQ2.AndCond(exprQ1), nil
		case OpOr:
			return exprQ2.OrCond(exprQ1), nil
		default:
			panic("Invalid operator")
		}
	}

	return new(orm.Condition), errors.New("Invalid operator")
}
