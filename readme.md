parse
----------

Simple parser constructor inspired by [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) which
describes grammars with composition of go types.

It generates syntax tree which could be used to fold the data inside into something useful.

This project is **experimental**. If you need industrial-grade parsers look at:

- [tree-sitter](https://tree-sitter.github.io/tree-sitter/)
- [ANTLR](https://www.antlr.org/)

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
	tree, err := Parse(
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
  expression{
    rule: *parse.Either(name: expression)(...)
    start: 1
    end: 2
    data: +
  }(
    operator{
      rule: *parse.Either(name: operator)(...)
      start: 1
      end: 2
      data: +
    }(
      +{
        rule: *parse.Terminal(name: +, value: +)()
        start: 1
        end: 2
        data: +
      }()
    )
  ),
  expression{
    rule: *parse.Either(name: expression)(...)
    start: 2
    end: 3
    data: (
  }(
    leftBracket{
      rule: *parse.Terminal(name: leftBracket, value: ()()
      start: 2
      end: 3
      data: (
    }()
  ),
  expression{
    rule: *parse.Either(name: expression)(...)
    start: 3
    end: 4
    data: 3
  }(
    numbers{
      rule: *parse.Repetition(name: numbers, times: 1, variadic: true)(...)
      start: 3
      end: 4
      data: 3
    }(
      number{
        rule: *parse.Either(name: number)(...)
        start: 3
        end: 4
        data: 3
      }(
        3{
          rule: *parse.Terminal(name: 3, value: 3)()
          start: 3
          end: 4
          data: 3
        }()
      )
    )
  ),
  expression{
    rule: *parse.Either(name: expression)(...)
    start: 4
    end: 5
    data: *
  }(
    operator{
      rule: *parse.Either(name: operator)(...)
      start: 4
      end: 5
      data: *
    }(
      *{
        rule: *parse.Terminal(name: *, value: *)()
        start: 4
        end: 5
        data: *
      }()
    )
  ),
  expression{
    rule: *parse.Either(name: expression)(...)
    start: 5
    end: 6
    data: 2
  }(
    numbers{
      rule: *parse.Repetition(name: numbers, times: 1, variadic: true)(...)
      start: 5
      end: 6
      data: 2
    }(
      number{
        rule: *parse.Either(name: number)(...)
        start: 5
        end: 6
        data: 2
      }(
        2{
          rule: *parse.Terminal(name: 2, value: 2)()
          start: 5
          end: 6
          data: 2
        }()
      )
    )
  ),
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
