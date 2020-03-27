package main

import (
	"fmt"
	"go/token"
)

type prefixParselet func(Token) exprNode

type infixParselet func(left exprNode, tok Token) exprNode

type exprParser struct {
	lexer lexer

	prefixParselets map[token.Token]prefixParselet
	infixParselets  map[token.Token]infixParselet

	prefixPrecedenceTab map[token.Token]int
	infixPrecedenceTab  map[token.Token]int
}

func newExprParser() *exprParser {
	p := &exprParser{
		prefixParselets:     make(map[token.Token]prefixParselet),
		infixParselets:      make(map[token.Token]infixParselet),
		prefixPrecedenceTab: make(map[token.Token]int),
		infixPrecedenceTab:  make(map[token.Token]int),
	}

	// Helper functions to bind the parselets.
	addPrefixParselet := func(tok token.Token, precedence int, parselet prefixParselet) {
		p.prefixParselets[tok] = parselet
		p.prefixPrecedenceTab[tok] = precedence
	}
	addInfixParselet := func(tok token.Token, precedence int, parselet infixParselet) {
		p.infixParselets[tok] = parselet
		p.infixPrecedenceTab[tok] = precedence
	}
	prefixExpr := func(precedence int, kinds ...token.Token) {
		for _, kind := range kinds {
			addPrefixParselet(kind, precedence, p.parsePrefixExpr)
		}
	}
	postfixExpr := func(precedence int, kinds ...token.Token) {
		for _, kind := range kinds {
			addInfixParselet(kind, precedence, p.parsePostfixExpr)
		}
	}
	leftAssocBinaryExpr := func(precedence int, kinds ...token.Token) {
		for _, kind := range kinds {
			addInfixParselet(kind, precedence, p.parseBinaryExpr)
		}
	}
	rightAssocBinaryExpr := func(precedence int, kinds ...token.Token) {
		for _, kind := range kinds {
			addInfixParselet(kind, precedence, p.rparseBinaryExpr)
		}
	}

	// Initialization of the parser tables.
	addPrefixParselet(token.LPAREN, 0, p.parseParenExpr)
	addPrefixParselet(token.IDENT, 0, p.parseNameExpr)
	addInfixParselet(token.LPAREN, 8, p.parseCallExpr)
	prefixExpr(6,
		token.ADD,
		token.SUB,
		token.INC,
		token.DEC,
	)
	postfixExpr(7,
		token.INC,
		token.DEC,
	)
	leftAssocBinaryExpr(3,
		token.ADD,
		token.SUB,
	)
	leftAssocBinaryExpr(4,
		token.MUL,
		token.QUO,
		token.REM,
	)
	rightAssocBinaryExpr(3,
		token.SHL,
	)

	return p
}

func (p *exprParser) ParseExpr(src []byte) (result exprNode, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if err2, ok := r.(error); ok {
			err = err2
			return
		}
		panic(r)
	}()
	p.lexer.Init(src)
	result = p.parseExpr(0)
	return result, nil
}

func (p *exprParser) parseExpr(precedence int) exprNode {
	tok := p.lexer.Consume()
	prefix, ok := p.prefixParselets[tok.kind]
	if !ok {
		panic(fmt.Errorf("unexpected token: %v", tok.kind))
	}
	left := prefix(tok)

	for precedence < p.infixPrecedenceTab[p.lexer.Peek().kind] {
		tok := p.lexer.Consume()
		infix := p.infixParselets[tok.kind]
		left = infix(left, tok)
	}

	return left
}

func (p *exprParser) parseNameExpr(tok Token) exprNode {
	return &nameExpr{Value: tok.value}
}

func (p *exprParser) parseParenExpr(tok Token) exprNode {
	x := p.parseExpr(0)
	p.expect(token.RPAREN)
	return x
}

func (p *exprParser) parsePrefixExpr(tok Token) exprNode {
	arg := p.parseExpr(p.prefixPrecedenceTab[tok.kind])
	return &prefixExpr{Op: tok.kind, Arg: arg}
}

func (p *exprParser) parsePostfixExpr(left exprNode, tok Token) exprNode {
	return &postfixExpr{Op: tok.kind, Arg: left}
}

func (p *exprParser) parseBinaryExpr(left exprNode, tok Token) exprNode {
	right := p.parseExpr(p.infixPrecedenceTab[tok.kind])
	return &binaryExpr{Op: tok.kind, Left: left, Right: right}
}

func (p *exprParser) rparseBinaryExpr(left exprNode, tok Token) exprNode {
	right := p.parseExpr(p.infixPrecedenceTab[tok.kind] - 1)
	return &binaryExpr{Op: tok.kind, Left: left, Right: right}
}

func (p *exprParser) parseCallExpr(left exprNode, tok Token) exprNode {
	if p.lexer.Peek().kind == token.RPAREN {
		// A call without arguments.
		p.lexer.Consume()
		return &callExpr{fn: left}
	}

	var args []exprNode
	for {
		args = append(args, p.parseExpr(0))
		if p.lexer.Peek().kind != token.COMMA {
			break
		}
		p.lexer.Consume()
	}
	p.expect(token.RPAREN)
	return &callExpr{fn: left, args: args}
}

func (p *exprParser) expect(want token.Token) {
	have := p.lexer.Peek().kind
	if have != want {
		panic(fmt.Errorf("expected %v, found %v", want, have))
	}
	p.lexer.Consume()
}
