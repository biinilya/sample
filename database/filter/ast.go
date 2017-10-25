// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate peg -switch -inline filter.peg

package main

import (
	"strconv"

	"fmt"

	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/linkedliststack"
)

type Operator string
type Comprator string
type ParameterType uint8

const (
	OpAnd Operator = "AND"
	OpOr  Operator = "OR"
	OpNot Operator = "NOT"

	CmpGt  Comprator = ">"
	CmpGte Comprator = ">="
	CmpLt  Comprator = "<"
	CmpLte Comprator = "<="
	CmpEq  Comprator = "="
)

const (
	TypeInt ParameterType = iota
	TypeString
)

func (op Operator) ArgsNum() int {
	if op == OpNot {
		return 1
	}
	return 2
}

type AST struct {
	exprStack stacks.Stack
	curExpr   expr
}

func (e *AST) Init() {
	e.exprStack = linkedliststack.New()
	e.curExpr = expr{}
}

func (e *AST) AddOperator(operator Operator) {
	fmt.Println("AddOperator", operator, e.exprStack.Size())
	var exprList = make([]expr, operator.ArgsNum())
	for idx, _ := range exprList {
		if e, found := e.exprStack.Pop(); !found {
			panic(operator)
		} else {
			exprList[idx] = e.(expr)
		}
	}
	e.exprStack.Push(expr{
		agg: &aggregation{
			exprList: exprList,
			operator: operator,
		},
	})
}

func (e *AST) AddComparator(comparator Comprator) {
	fmt.Println("AddComparator", comparator)
	e.curExpr.comparator = comparator
}

func (e *AST) AddIdentifier(identifier string) {
	fmt.Println("AddIdentifier", identifier)
	e.curExpr.idenitifier = identifier
}

func (e *AST) AddArgument(t ParameterType, parameter string) {
	fmt.Println("AddArgument", parameter)
	switch t {
	case TypeInt:
		if intV, intErr := strconv.ParseInt(parameter, 10, 64); intErr != nil {
			panic(intErr)
		} else {
			e.curExpr.args = append(e.curExpr.args, intV)
		}
	case TypeString:
		e.curExpr.args = append(e.curExpr.args, parameter)
	default:
		panic("Unknown parameter type")
	}
}

func (e *AST) AddExpression() {
	e.exprStack.Push(e.curExpr)
	fmt.Println("AddExpression", e.exprStack.Size())
	e.curExpr = expr{}
}

func (e *AST) topExpr() expr {
	if e.exprStack.Size() != 1 {
		panic(e.exprStack.Size())
	}
	var lExpr, _ = e.exprStack.Peek()
	return lExpr.(expr)
}

func (e *AST) Render(fields ...string) (string, []interface{}, error) {
	return e.topExpr().renderToSql(fields...)
}
