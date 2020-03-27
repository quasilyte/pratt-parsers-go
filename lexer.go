package main

import (
	"go/scanner"
	"go/token"
)

type Token struct {
	kind  token.Token
	value string
}

type lexer struct {
	s scanner.Scanner

	peeked bool
	tok    Token
}

func (l *lexer) Init(src []byte) {
	l.peeked = false
	fset := token.NewFileSet()
	file := fset.AddFile("<expr.go>", fset.Base(), len(src))
	l.s.Init(file, src, nil, scanner.ScanComments)
}

func (l *lexer) Consume() Token {
	if l.peeked {
		l.peeked = false
		return l.tok
	}
	_, kind, value := l.s.Scan()
	return Token{kind: kind, value: value}
}

func (l *lexer) Peek() Token {
	l.tok = l.Consume()
	l.peeked = true
	return l.tok
}
