// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate peg -switch -inline filter.peg

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/linkedliststack"
)

type Operator string
type Comprator string
type ParameterType uint8

const (
	OpAnd Operator = "AND"
	OpNot Operator = "NOT"
	OpOr  Operator = "OR"

	CmpGt  Comprator = ">"
	CmpGte Comprator = ">="
	CmpLt  Comprator = "<"
	CmpLte Comprator = "<="
	CmpEq  Comprator = "="

	TypeInt ParameterType = iota
	TypeString
)

type Operation struct {
	Statement  string
	Parameters []interface{}
}

type Expression struct {
	opStack    stacks.Stack
	identifier string
	comparator Comprator
	parameter  interface{}
	fieldMap   map[string]string
	fieldErr   error
}

func (e *Expression) Init(fieldMap map[string]string) {
	e.opStack = linkedliststack.New()
	e.fieldMap = fieldMap
	e.fieldErr = nil
}

func (e *Expression) AddOperator2(operator Operator) {
	var opL, opR Operation
	if op, ok := e.opStack.Pop(); ok {
		opR = op.(Operation)
	} else {
		panic("Broken operation stack")
	}
	if op, ok := e.opStack.Pop(); ok {
		opL = op.(Operation)
	} else {
		panic("Broken operation stack")
	}
	var merged Operation
	merged.Statement = fmt.Sprintf(`(%s) %s (%s)`, opL.Statement, operator, opR.Statement)
	merged.Parameters = append(opL.Parameters, opR.Parameters...)
	e.opStack.Push(merged)
}

func (e *Expression) AddOperator1(operator Operator) {
	if e.fieldErr != nil {
		return
	}

	var opS Operation
	if op, ok := e.opStack.Pop(); ok {
		opS = op.(Operation)
	} else {
		panic("Broken operation stack")
	}
	var merged Operation
	merged.Statement = fmt.Sprintf(`%s (%s)`, operator, opS.Statement)
	merged.Parameters = opS.Parameters
	e.opStack.Push(merged)
}

func (e *Expression) AddComparator(comparator Comprator) {
	if e.fieldErr != nil {
		return
	}

	e.comparator = comparator
}

func (e *Expression) AddIdentifier(identifier string) {
	if e.fieldErr != nil {
		return
	}

	e.identifier = identifier
	if _, found := e.fieldMap[identifier]; !found {
		e.fieldErr = errors.New(fmt.Sprintf("Unknown field <%s>", identifier))
	}
}

func (e *Expression) AddParameter(t ParameterType, parameter string) {
	if e.fieldErr != nil {
		return
	}

	var op Operation
	op.Statement = fmt.Sprintf(`"%s" %s ?::%s`, e.identifier, e.comparator, e.fieldMap[e.identifier])

	switch t {
	case TypeInt:
		if intV, intErr := strconv.ParseInt(parameter, 10, 64); intErr != nil {
			panic(intErr)
		} else {
			op.Parameters = []interface{}{intV}
		}
	case TypeString:
		op.Parameters = []interface{}{parameter}
	default:
		panic("Unknown parameter type")
	}
	e.opStack.Push(op)
}

func (e *Expression) Render() (string, []interface{}, error) {
	if e.fieldErr != nil {
		return "", nil, e.fieldErr
	}

	var opS Operation
	if op, ok := e.opStack.Peek(); ok {
		opS = op.(Operation)
	} else {
		panic("Broken operation stack")
	}
	return opS.Statement, opS.Parameters, nil
}

func ApplyFilter(
	fieldMap map[string]string,
	rawQ string,
	q orm.QuerySeter,
) (filtered orm.QuerySeter, err error) {
	filter := &Filter{}
	filter.Buffer = rawQ
	filter.Init()
	filter.Expression.Init(fieldMap)
	if filterErr := filter.Parse(); filterErr != nil {
		return q, filterErr
	}

	filter.Execute()
	var qS, qA, qErr = filter.Render()
	if qErr != nil {
		return q, qErr
	}
	return q.Filter(qS, qA...), nil
}
