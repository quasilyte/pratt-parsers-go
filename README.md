# Pratt Parsers: Expression Parsing Made Easy

## Go edition

This implementation is inspired by the [Bantam](https://github.com/munificent/bantam) parser and
[Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) article.

In Go edition, we parse a Go dialect instead, using the [go/scanner](https://golang.org/pkg/go/scanner/) and [go/token](https://golang.org/pkg/go/token/)
from the standard library.

Unlike normal Go, the dialect we're going to parse:

* Defines `++` and `--` as expressions
* Has prefix forms of increment and decrement, `--x` and `++x` is valid
* Left bitwise shift is right-associative (for the demonstration purposes)

This repository has associated articles: [RU](https://habr.com/ru/post/494316/), EN (TODO).
