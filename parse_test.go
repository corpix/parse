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
package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	samples := []struct {
		text   string
		rule   Rule
		tree   *Tree
		err    error
		parser *Parser
	}{

		// Terminal

		{
			"",
			NewTerminal(""),
			nil,
			NewErrEmptyRule(NewTerminal("")),
			DefaultParser,
		},
		{
			"foo",
			NewTerminal("foo"),
			&Tree{
				Rule:  NewTerminal("foo"),
				Data:  []byte("foo"),
				Start: 0,
				End:   3,
			},
			nil,
			DefaultParser,
		},
		{
			"bar",
			NewTerminal("foo"),
			nil,
			NewErrUnexpectedToken(
				[]byte("b"),
				1,
				NewTerminal("foo"),
			),
			DefaultParser,
		},
		{
			"foobar",
			NewTerminal("foo"),
			nil,
			NewErrUnexpectedToken(
				[]byte("b"),
				4,
				NewTerminal("foo"),
			),
			DefaultParser,
		},

		// Chain

		{
			"",
			NewChain(),
			nil,
			NewErrEmptyRule(NewChain()),
			DefaultParser,
		},
		{
			"foo bar",
			NewChain(
				NewTerminal("foo"),
				NewTerminal(" "),
				NewTerminal("bar"),
			),
			&Tree{
				Rule: NewChain(
					NewTerminal("foo"),
					NewTerminal(" "),
					NewTerminal("bar"),
				),
				Data:  []byte("foo bar"),
				Start: 0,
				End:   7,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("foo"),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
					},
					{
						Rule:  NewTerminal(" "),
						Data:  []byte(" "),
						Start: 3,
						End:   4,
					},
					{
						Rule:  NewTerminal("bar"),
						Data:  []byte("bar"),
						Start: 4,
						End:   7,
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"foo",
			NewChain(
				NewChain(
					NewChain(
						NewChain(
							NewTerminal("foo"),
						),
					),
				),
			),
			nil,
			NewErrNestingTooDeep(4, 1),
			NewParser(3),
		},

		// Either

		{
			"",
			NewEither(),
			nil,
			NewErrEmptyRule(NewEither()),
			DefaultParser,
		},
		{
			"1",
			NewEither(
				NewTerminal("1"),
				NewTerminal("2"),
				NewTerminal("3"),
			),
			&Tree{
				Rule: NewEither(
					NewTerminal("1"),
					NewTerminal("2"),
					NewTerminal("3"),
				),
				Data:  []byte("1"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("1"),
						Data:  []byte("1"),
						Start: 0,
						End:   1,
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"1",
			NewEither(
				NewTerminal("3"),
				NewTerminal("2"),
				NewTerminal("1"),
			),
			&Tree{
				Rule: NewEither(
					NewTerminal("3"),
					NewTerminal("2"),
					NewTerminal("1"),
				),
				Data:  []byte("1"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("1"),
						Data:  []byte("1"),
						Start: 0,
						End:   1,
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"2",
			NewEither(
				NewTerminal("1"),
				NewTerminal("2"),
				NewTerminal("3"),
			),
			&Tree{
				Rule: NewEither(
					NewTerminal("1"),
					NewTerminal("2"),
					NewTerminal("3"),
				),
				Data:  []byte("2"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("2"),
						Data:  []byte("2"),
						Start: 0,
						End:   1,
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"4",
			NewEither(
				NewTerminal("1"),
				NewTerminal("2"),
				NewTerminal("3"),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("4"),
				1,
				NewTerminal("3"),
			),
			DefaultParser,
		},

		// Repetition

		{
			"123",
			NewRepetition(
				NewEither(
					NewTerminal("1"),
					NewTerminal("2"),
					NewTerminal("3"),
				),
			),
			&Tree{
				Rule: NewRepetition(
					NewEither(
						NewTerminal("1"),
						NewTerminal("2"),
						NewTerminal("3"),
					),
				),
				Data:  []byte("123"),
				Start: 0,
				End:   3,
				Childs: []*Tree{
					{
						Rule: NewEither(
							NewTerminal("1"),
							NewTerminal("2"),
							NewTerminal("3"),
						),
						Data:  []byte("1"),
						Start: 0,
						End:   1,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("1"),
								Data:  []byte("1"),
								Start: 0,
								End:   1,
							},
						},
					},
					{
						Rule: NewEither(
							NewTerminal("1"),
							NewTerminal("2"),
							NewTerminal("3"),
						),
						Data:  []byte("2"),
						Start: 1,
						End:   2,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("2"),
								Data:  []byte("2"),
								Start: 1,
								End:   2,
							},
						},
					},
					{
						Rule: NewEither(
							NewTerminal("1"),
							NewTerminal("2"),
							NewTerminal("3"),
						),
						Data:  []byte("3"),
						Start: 2,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("3"),
								Data:  []byte("3"),
								Start: 2,
								End:   3,
							},
						},
					},
				},
			},
			nil,
			DefaultParser,
		},

		// mixed

		{
			"foobar",
			NewChain(
				NewEither(
					NewTerminal("foo"),
					NewTerminal("bar"),
				),
				NewEither(
					NewTerminal("foo"),
					NewTerminal("bar"),
				),
			),
			&Tree{
				Rule: NewChain(
					NewEither(
						NewTerminal("foo"),
						NewTerminal("bar"),
					),
					NewEither(
						NewTerminal("foo"),
						NewTerminal("bar"),
					),
				),
				Data:  []byte("foobar"),
				Start: 0,
				End:   6,
				Childs: []*Tree{
					{
						Rule: NewEither(
							NewTerminal("foo"),
							NewTerminal("bar"),
						),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("foo"),
								Data:  []byte("foo"),
								Start: 0,
								End:   3,
							},
						},
					},
					{
						Rule: NewEither(
							NewTerminal("foo"),
							NewTerminal("bar"),
						),
						Data:  []byte("bar"),
						Start: 3,
						End:   6,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("bar"),
								Data:  []byte("bar"),
								Start: 3,
								End:   6,
							},
						},
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"foo(1234)",
			NewChain(
				NewTerminal("foo"),
				NewTerminal("("),
				NewRepetition(
					NewEither(
						NewTerminal("4"),
						NewTerminal("3"),
						NewTerminal("2"),
						NewTerminal("1"),
					),
				),
				NewTerminal(")"),
			),
			&Tree{
				Rule: NewChain(
					NewTerminal("foo"),
					NewTerminal("("),
					NewRepetition(
						NewEither(
							NewTerminal("4"),
							NewTerminal("3"),
							NewTerminal("2"),
							NewTerminal("1"),
						),
					),
					NewTerminal(")"),
				),
				Data:  []byte("foo(1234)"),
				Start: 0,
				End:   9,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("foo"),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
					},
					{
						Rule:  NewTerminal("("),
						Data:  []byte("("),
						Start: 3,
						End:   4,
					},
					{
						Rule: NewRepetition(
							NewEither(
								NewTerminal("4"),
								NewTerminal("3"),
								NewTerminal("2"),
								NewTerminal("1"),
							),
						),
						Data:  []byte("1234"),
						Start: 4,
						End:   8,
						Childs: []*Tree{
							{
								Rule: NewEither(
									NewTerminal("4"),
									NewTerminal("3"),
									NewTerminal("2"),
									NewTerminal("1"),
								),
								Data:  []byte("1234"),
								Start: 4,
								End:   8,
								Childs: []*Tree{
									{
										Rule:  NewTerminal("1"),
										Data:  []byte("1"),
										Start: 4,
										End:   5,
									},
									{
										Rule:  NewTerminal("2"),
										Data:  []byte("2"),
										Start: 5,
										End:   6,
									},
									{
										Rule:  NewTerminal("3"),
										Data:  []byte("3"),
										Start: 6,
										End:   7,
									},
									{
										Rule:  NewTerminal("4"),
										Data:  []byte("4"),
										Start: 7,
										End:   8,
									},
								},
							},
						},
					},
					{
						Rule:  NewTerminal(")"),
						Data:  []byte(")"),
						Start: 8,
						End:   9,
					},
				},
			},
			nil,
			DefaultParser,
		},
	}

	for k, sample := range samples {
		tree, err := sample.parser.Parse(
			sample.rule,
			[]byte(sample.text),
		)
		msg := spew.Sdump(
			k,
			sample.rule,
			sample.text,
		)
		assert.EqualValues(t, sample.err, err, msg)
		assert.EqualValues(t, sample.tree, tree, msg)
	}
}

// func TestTemporary(t *testing.T) {
// 	number := NewEither(
// 		NewTerminal("0"),
// 		NewTerminal("1"),
// 		NewTerminal("2"),
// 		NewTerminal("3"),
// 		NewTerminal("4"),
// 		NewTerminal("5"),
// 		NewTerminal("6"),
// 		NewTerminal("7"),
// 		NewTerminal("8"),
// 		NewTerminal("9"),
// 	)
// 	leftBracket := NewTerminal("(")
// 	rightBracket := NewTerminal(")")

// 	typeDefinition := NewEither()

// 	wrap := NewChain(
// 		NewTerminal("Wrap"),
// 		leftBracket,
// 		typeDefinition,
// 		rightBracket,
// 	)
// 	node := NewChain(
// 		NewTerminal("Node"),
// 		leftBracket,
// 		typeDefinition,
// 		rightBracket,
// 	)
// 	integer := NewChain(
// 		NewTerminal("Int"),
// 		leftBracket,
// 		NewRepetition(number),
// 		rightBracket,
// 	)

// 	*typeDefinition = append(*typeDefinition, wrap)
// 	*typeDefinition = append(*typeDefinition, node)
// 	*typeDefinition = append(*typeDefinition, integer)

// 	spew.Dump(
// 		DefaultParser.Parse(
// 			typeDefinition,
// 			[]byte("Wrap(Node(Int(5)))"),
// 		),
// 	)
// }
