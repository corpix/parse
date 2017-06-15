parse
----------

[![Build Status](https://travis-ci.org/corpix/parse.svg?branch=master)](https://travis-ci.org/corpix/parse)

Simple parser constructor inspired by [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) which
describes grammars with composition of go types.

It generates syntax tree which could be used to fold the data inside into something useful.

This project is in **alpha state**, API may change in future.

## Limitations

- Has no concept of `string literal`, you could parse `"foo\"bar"` but you should fold the AST by hands.
  This may change in future, I think we could introduce a separate rule type for this.
- Line reporting in AST is not implemented at this time, it reports only position in the string.
  This will change in the future, I think we could introduce an option to create a parser
  which will configure the line-break symbols.

## Example

> Example from `examples/` directory.

``` go
package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/corpix/parse"
)

var (
	expression parse.Rule
)

func init() {
	numbers := parse.NewRepetition(
		"numbers",
		parse.NewEither(
			"number",
			parse.NewTerminal("1", "1"),
			parse.NewTerminal("2", "2"),
			parse.NewTerminal("3", "3"),
			parse.NewTerminal("4", "4"),
			parse.NewTerminal("5", "5"),
			parse.NewTerminal("6", "6"),
			parse.NewTerminal("7", "7"),
			parse.NewTerminal("8", "8"),
			parse.NewTerminal("9", "9"),
			parse.NewTerminal("0", "0"),
		),
	)

	operator := parse.NewEither(
		"operator",
		parse.NewTerminal("+", "+"),
		parse.NewTerminal("-", "-"),
		parse.NewTerminal("*", "*"),
		parse.NewTerminal("/", "/"),
		parse.NewTerminal("mod", "mod"),
	)

	whitespace := parse.NewEither(
		"whitespace",
		parse.NewTerminal("space", " "),
		parse.NewTerminal("tab", "\t"),
		parse.NewTerminal("line-break", "\n"),
	)

	leftBracket := parse.NewTerminal("(", "(")
	rightBracket := parse.NewTerminal(")", ")")

	expression = parse.NewRepetition(
		"expressions",
		parse.NewEither(
			"expression",
			numbers,
			whitespace,
			operator,
			leftBracket,
			rightBracket,
		),
	)

}

func main() {
	tree, err := parse.Parse(
		expression,
		[]byte("5 + 3 * (4 - 3) / (7 + 3)"),
	)
	if err != nil {
		panic(err)
	}

	s := spew.NewDefaultConfig()
	s.MaxDepth = 1
	s.Dump(tree)
}
```

Save this code into some file, `main.go` for example.

Run it:

``` shell
go run main.go
(*parse.Tree)(0xc4200163c0)(*parse.Repetition(ID: expressions, Times: 1, Variadic: true)(
  *parse.Either(ID: expression)(
        *parse.Repetition(ID: numbers, Times: 1, Variadic: true)(
                  *parse.Either(ID: number)(
                                *parse.Terminal(ID: 1, Value: [49])(),
                                *parse.Terminal(ID: 2, Value: [50])(),
                                *parse.Terminal(ID: 3, Value: [51])(),
                                *parse.Terminal(ID: 4, Value: [52])(),
                                *parse.Terminal(ID: 5, Value: [53])(),
                                *parse.Terminal(ID: 6, Value: [54])(),
                                *parse.Terminal(ID: 7, Value: [55])(),
                                *parse.Terminal(ID: 8, Value: [56])(),
                                *parse.Terminal(ID: 9, Value: [57])(),
                                *parse.Terminal(ID: 0, Value: [48])()
                  )
        ),
        *parse.Either(ID: whitespace)(
                  *parse.Terminal(ID: space, Value: [32])(),
                  *parse.Terminal(ID: tab, Value: [9])(),
                  *parse.Terminal(ID: line-break, Value: [10])()
        ),
        *parse.Either(ID: operator)(
                  *parse.Terminal(ID: +, Value: [43])(),
                  *parse.Terminal(ID: -, Value: [45])(),
                  *parse.Terminal(ID: *, Value: [42])(),
                  *parse.Terminal(ID: /, Value: [47])(),
                  *parse.Terminal(ID: mod, Value: [109 111 100])()
        ),
        *parse.Terminal(ID: (, Value: [40])(),
        *parse.Terminal(ID: ), Value: [41])()
  )
))
```
