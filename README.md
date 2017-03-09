parse
----------

Simple parser constructor inspired by [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) which
describes grammars with composition of go types.

It generates syntax tree which could be used to fold the data inside into something useful.

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

	spew.Dump(tree.Rule.ID())
	spew.Printf("Rule is terminal? %#v\n", tree.Rule.IsTerminal())
	spew.Dump(tree.Data)
}
```

Save this code into some file, `main.go` for example and run:

``` shell
go run main.go
(string) (len=11) "expressions"
Rule is terminal? (bool)false
([]uint8) (len=25 cap=32) {
 00000000  35 20 2b 20 33 20 2a 20  28 34 20 2d 20 33 29 20  |5 + 3 * (4 - 3) |
 00000010  2f 20 28 37 20 2b 20 33  29                       |/ (7 + 3)|
}
```
