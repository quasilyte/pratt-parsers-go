package main

import (
	"fmt"
	"go/token"
	"strings"
)

type exprNode interface {
	expr()
	String() string
}

type nameExpr struct {
	Value string
}

type prefixExpr struct {
	Op  token.Token
	Arg exprNode
}

type postfixExpr struct {
	Op  token.Token
	Arg exprNode
}

type binaryExpr struct {
	Op    token.Token
	Left  exprNode
	Right exprNode
}

type callExpr struct {
	fn   exprNode
	args []exprNode
}

func (e *nameExpr) expr()    {}
func (e *prefixExpr) expr()  {}
func (e *postfixExpr) expr() {}
func (e *binaryExpr) expr()  {}
func (e *callExpr) expr()    {}

func (e *nameExpr) String() string    { return e.Value }
func (e *prefixExpr) String() string  { return fmt.Sprintf("(prefix %s %s)", e.Op, e.Arg) }
func (e *postfixExpr) String() string { return fmt.Sprintf("(postfix %s %s)", e.Op, e.Arg) }
func (e *binaryExpr) String() string  { return fmt.Sprintf("(%s %s %s)", e.Op, e.Left, e.Right) }

func (e *callExpr) String() string {
	parts := make([]string, 0, len(e.args)+1)
	parts = append(parts, e.fn.String())
	for _, arg := range e.args {
		parts = append(parts, arg.String())
	}
	return fmt.Sprintf("(call %s)", strings.Join(parts, " "))
}
