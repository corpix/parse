parse
----------

Simple parser constructor inspired by [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) which
describes grammars with composition of go types.

It generates syntax tree which could be used to fold the data inside into something useful.

This project is in **alpha state**, API may change in future.

## Limitations

- Has no concept of `string literal`, you could parse `"foo\"bar"` but you should fold the AST by hands
  This may change in future, I think we could introduce a separate rule type for this.
- Line reporting in AST is not implemented at this time, it reports only position in the string.
  This will change in the future, I think we could introduce an option to create a parser
  which will configure the line-break symbols.

## Example

> Example from `examples/` directory.

``` go
package main

// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
(*parse.Tree)(0xc4200163c0)({
 Rule: (*parse.Repetition)(0xc420010480)({
  <max depth reached>
 }),
 Start: (int) 0,
 End: (int) 25,
 Childs: ([]*parse.Tree) (len=25 cap=32) {
  <max depth reached>
 },
 Data: ([]uint8) (len=25 cap=32) {
  <max depth reached>
 }
})
```
