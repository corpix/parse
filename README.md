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

``` go
package main

import (
	"fmt"

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

	leftBracket := parse.NewTerminal("leftBracket", "(")
	rightBracket := parse.NewTerminal("rightBracket", ")")

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
		[]byte("5+(3*2)"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(tree)
}
```

> Console output shortened

``` console
$ go run ./examples/calculator-expression/calculator-expression.go
expressions{
  rule: *parse.Repetition(name: expressions, times: 1, variadic: true)(...)
  start: 0
  end: 7
  data: 5+(3*2)
}(
    expression{
      rule: *parse.Either(name: expression)(...)
      start: 0
      end: 1
      data: 5
    }(
        numbers{
          rule: *parse.Repetition(name: numbers, times: 1, variadic: true)(...)
          start: 0
          end: 1
          data: 5
        }(
            number{
              rule: *parse.Either(name: number)(...)
              start: 0
              end: 1
              data: 5
            }(
                5{
                  rule: *parse.Terminal(name: 5, value: 5)()
                  start: 0
                  end: 1
                  data: 5
                }()
            )
        )
    ),

    ...

    expression{
      rule: *parse.Either(name: expression)(...)
      start: 6
      end: 7
      data: )
    }(
        rightBracket{
          rule: *parse.Terminal(name: rightBracket, value: ))()
          start: 6
          end: 7
          data: )
        }()
    )
)
```
