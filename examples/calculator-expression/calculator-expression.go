package main

import (
	"fmt"

	. "github.com/corpix/parse"
)

var (
	expression Rule
)

func init() {
	numbers := NewRepetition(
		"numbers",
		NewEither(
			"number",
			NewTerminal("1", "1"),
			NewTerminal("2", "2"),
			NewTerminal("3", "3"),
			NewTerminal("4", "4"),
			NewTerminal("5", "5"),
			NewTerminal("6", "6"),
			NewTerminal("7", "7"),
			NewTerminal("8", "8"),
			NewTerminal("9", "9"),
			NewTerminal("0", "0"),
		),
	)

	operator := NewEither(
		"operator",
		NewTerminal("+", "+"),
		NewTerminal("-", "-"),
		NewTerminal("*", "*"),
		NewTerminal("/", "/"),
		NewTerminal("mod", "mod"),
	)

	whitespace := NewEither(
		"whitespace",
		NewTerminal("space", " "),
		NewTerminal("tab", "\t"),
		NewTerminal("line-break", "\n"),
	)

	leftBracket := NewTerminal("leftBracket", "(")
	rightBracket := NewTerminal("rightBracket", ")")

	expression = NewRepetition(
		"expressions",
		NewEither(
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
	p := NewParser(ParserOptionPath("calculator-expression.go"))
	tree, err := p.Parse(
		expression,
		[]byte("5+(3*2)"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(tree)
}
