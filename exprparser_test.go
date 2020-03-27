package main

import (
	"testing"
)

func TestParseGoodExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`x`, `x`},

		// Tests for prefix expr.
		{`-x`, `(prefix - x)`},
		{`+x`, `(prefix + x)`},
		{`-(-x)`, `(prefix - (prefix - x))`},
		{`-(+x)`, `(prefix - (prefix + x))`},
		{`+(-x)`, `(prefix + (prefix - x))`},
		{`--x`, `(prefix -- x)`},
		{`++x`, `(prefix ++ x)`},

		// Tests for binary expr.
		{`x-y`, `(- x y)`},
		{`x-y-z`, `(- (- x y) z)`},
		{`x + y + z`, `(+ (+ x y) z)`},
		{`x+y*z`, `(+ x (* y z))`},
		{`x*y+z`, `(+ (* x y) z)`},

		// Tests for right-assoc binary expr.
		{`x<<y`, `(<< x y)`},
		{`x << y << z`, `(<< x (<< y z))`},

		// Tests for paren expr.
		{`x+(y*z)`, `(+ x (* y z))`},
		{`(x*y)+z`, `(+ (* x y) z)`},
		{`(x+y)*z`, `(* (+ x y) z)`},
		{`x*(y+z)`, `(* x (+ y z))`},

		// Tests for postfix expr.
		{`x++`, `(postfix ++ x)`},
		{`x--`, `(postfix -- x)`},
		{`x++ + y++`, `(+ (postfix ++ x) (postfix ++ y))`},

		// Tests for call expr.
		{`f()`, `(call f)`},
		{`f(x)`, `(call f x)`},
		{`f(x, y)`, `(call f x y)`},
		{`f(g(x, y), z)`, `(call f (call g x y) z)`},
		{`f(x)(y,z)`, `(call (call f x) y z)`},
		{`f(x++, y++)`, `(call f (postfix ++ x) (postfix ++ y))`},

		// Mixed tests.
		{`-x + -y`, `(+ (prefix - x) (prefix - y))`},
		{`+x + f(+y)`, `(+ (prefix + x) (call f (prefix + y)))`},
	}

	p := newExprParser()
	for _, test := range tests {
		expr, err := p.ParseExpr([]byte(test.input))
		if err != nil {
			t.Errorf("parse(%q): error: %v", test.input, err)
			continue
		}
		if expr.String() != test.want {
			t.Errorf("parse(%q):\nhave: %q\nwant: %q",
				test.input, expr.String(), test.want)
		}
	}
}

func TestParseBadExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`(x;`, `expected ), found ;`},
		{`f(x, y;`, `expected ), found ;`},
		{`x+`, `unexpected token: EOF`},
		{`<x`, `unexpected token: <`},
	}

	p := newExprParser()
	for _, test := range tests {
		_, err := p.ParseExpr([]byte(test.input))
		have := "<nil>"
		if err != nil {
			have = err.Error()
		}
		if have != test.want {
			t.Errorf("parse(%q) error:\nhave: %s\nwant: %s",
				test.input, have, test.want)
		}
	}
}
