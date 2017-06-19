package parse

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
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParseInputIntegrity(t *testing.T) {
	samples := []struct {
		input   string
		grammar Rule
		err     error
	}{
		{"foo", NewTerminal("foo", "foo"), nil},
		{
			"foo bar",
			NewChain(
				"chain of foo bar",
				NewTerminal("foo", "foo"),
				NewTerminal("space", " "),
				NewTerminal("bar", "bar"),
			),
			nil,
		},
		{
			"1",
			NewEither(
				"number",
				NewTerminal("one", "1"),
				NewTerminal("two", "2"),
				NewTerminal("three", "3"),
			),
			nil,
		},
		{
			"123",
			NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			nil,
		},
		{
			"123",
			NewRepetitionTimesVariadic(
				"numbers",
				2,
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			nil,
		},
		{
			"foobar",
			NewChain(
				"chain of foo bar",
				NewEither(
					"foo and bar",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
				NewEither(
					"foo and bar",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
			),
			nil,
		},
		{
			"foo(1234)",
			NewChain(
				"foo func",
				NewTerminal("foo", "foo"),
				NewTerminal("left bracket", "("),
				NewRepetition(
					"numbers",
					NewEither(
						"number",
						NewTerminal("four", "4"),
						NewTerminal("three", "3"),
						NewTerminal("two", "2"),
						NewTerminal("one", "1"),
					),
				),
				NewTerminal("right bracket", ")"),
			),
			nil,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		input := make([]byte, len(sample.input))
		copy(input, []byte(sample.input))
		_, err := Parse(sample.grammar, input)

		assert.Equal(t, []byte(sample.input), input, msg)
		assert.Equal(t, sample.err, err, msg)
	}
}

func TestParse(t *testing.T) {
	samples := []struct {
		text   string
		rule   Rule
		tree   *Tree
		err    error
		parser *Parser
	}{
		// errors

		{
			"",
			newTestRuleFinite("empty"),
			nil,
			NewErrEmptyRule(
				newTestRuleFinite("empty"),
				nil,
			),
			DefaultParser,
		},
		{
			"",
			newTestRuleNonFinite("empty", newTestRuleFinite("emptyInside")),
			nil,
			NewErrEmptyRule(
				newTestRuleFinite("emptyInside"),
				newTestRuleNonFinite("empty", newTestRuleFinite("emptyInside")),
			),
			DefaultParser,
		},
		{
			"foo",
			newTestRuleNonFinite(
				"deep",
				newTestRuleNonFinite(
					"deep",
					newTestRuleNonFinite(
						"deep",
						newTestRuleNonFinite(
							"deep",
							newTestRuleFinite("foo"),
						),
					),
				),
			),
			nil,
			NewErrNestingTooDeep(4, 1),
			NewParser(3),
		},

		// Terminal

		{
			"",
			NewTerminal("empty", ""),
			nil,
			NewErrEmptyRule(
				NewTerminal("empty", ""),
				nil,
			),
			DefaultParser,
		},
		{
			"foo",
			NewTerminal("foo", "foo"),
			&Tree{
				Rule:  NewTerminal("foo", "foo"),
				Data:  []byte("foo"),
				Start: 0,
				End:   3,
			},
			nil,
			DefaultParser,
		},
		{
			"bar",
			NewTerminal("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("bar")),
				1,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"foobar",
			NewTerminal("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("bar")),
				4,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},

		// Chain

		{
			"foo bar",
			NewChain(
				"chain of foo bar",
				NewTerminal("foo", "foo"),
				NewTerminal("space", " "),
				NewTerminal("bar", "bar"),
			),
			&Tree{
				Rule: NewChain(
					"chain of foo bar",
					NewTerminal("foo", "foo"),
					NewTerminal("space", " "),
					NewTerminal("bar", "bar"),
				),
				Data:  []byte("foo bar"),
				Start: 0,
				End:   7,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("foo", "foo"),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
					},
					{
						Rule:  NewTerminal("space", " "),
						Data:  []byte(" "),
						Start: 3,
						End:   4,
					},
					{
						Rule:  NewTerminal("bar", "bar"),
						Data:  []byte("bar"),
						Start: 4,
						End:   7,
					},
				},
			},
			nil,
			DefaultParser,
		},

		// Either

		{
			"",
			NewEither("empty"),
			nil,
			NewErrEmptyRule(
				NewEither("empty"),
				nil,
			),
			DefaultParser,
		},
		{
			"1",
			NewEither(
				"number",
				NewTerminal("one", "1"),
				NewTerminal("two", "2"),
				NewTerminal("three", "3"),
			),
			&Tree{
				Rule: NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
				Data:  []byte("1"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("one", "1"),
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
				"number",
				NewTerminal("three", "3"),
				NewTerminal("two", "2"),
				NewTerminal("one", "1"),
			),
			&Tree{
				Rule: NewEither(
					"number",
					NewTerminal("three", "3"),
					NewTerminal("two", "2"),
					NewTerminal("one", "1"),
				),
				Data:  []byte("1"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("one", "1"),
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
				"number",
				NewTerminal("one", "1"),
				NewTerminal("two", "2"),
				NewTerminal("three", "3"),
			),
			&Tree{
				Rule: NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
				Data:  []byte("2"),
				Start: 0,
				End:   1,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("two", "2"),
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
				"number",
				NewTerminal("one", "1"),
				NewTerminal("two", "2"),
				NewTerminal("three", "3"),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("4"),
				1,
				NewTerminal("three", "3"),
			),
			DefaultParser,
		},

		// Repetition

		{
			"",
			NewRepetitionTimesVariadic(
				"maybe",
				0,
				NewTerminal("terminal", "t"),
			),
			nil,
			NewErrUnexpectedEOF(
				1,
				NewRepetitionTimesVariadic(
					"maybe",
					0,
					NewTerminal("terminal", "t"),
				),
			),
			DefaultParser,
		},
		{
			"t",
			NewRepetitionTimesVariadic(
				"maybe",
				0,
				NewTerminal("terminal", "t"),
			),
			&Tree{
				Rule: NewRepetitionTimesVariadic(
					"maybe",
					0,
					NewTerminal("terminal", "t"),
				),
				Start: 0,
				End:   1,
				Data:  []byte("t"),
				Childs: []*Tree{
					{
						Rule:  NewTerminal("terminal", "t"),
						Start: 0,
						End:   1,
						Data:  []byte("t"),
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"f",
			NewRepetition(
				"set",
				NewEither(
					"variant",
					NewTerminal("terminal", "f"),
					NewRepetitionTimesVariadic(
						"maybe",
						0,
						NewTerminal("terminal", "t"),
					),
				),
			),
			&Tree{
				Rule: NewRepetition(
					"set",
					NewEither(
						"variant",
						NewTerminal("terminal", "f"),
						NewRepetitionTimesVariadic(
							"maybe",
							0,
							NewTerminal("terminal", "t"),
						),
					),
				),
				Start: 0,
				End:   1,
				Data:  []byte("f"),
				Childs: []*Tree{
					{
						Rule: NewEither(
							"variant",
							NewTerminal("terminal", "f"),
							NewRepetitionTimesVariadic(
								"maybe",
								0,
								NewTerminal("terminal", "t"),
							),
						),
						Start: 0,
						End:   1,
						Data:  []byte("f"),
						Childs: []*Tree{
							{
								Rule:  NewTerminal("terminal", "f"),
								Start: 0,
								End:   1,
								Data:  []byte("f"),
							},
						},
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"tf",
			NewChain(
				"chain",
				NewTerminal("t", "t"),
				NewRepetitionTimesVariadic(
					"space",
					0,
					NewTerminal("space", " "),
				),
				NewTerminal("f", "f"),
			),
			&Tree{
				Rule: NewChain(
					"chain",
					NewTerminal("t", "t"),
					NewRepetitionTimesVariadic(
						"space",
						0,
						NewTerminal("space", " "),
					),
					NewTerminal("f", "f"),
				),
				Start: 0,
				End:   2,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("t", "t"),
						Start: 0,
						End:   1,
						Data:  []byte("t"),
					},
					{
						Rule:  NewTerminal("f", "f"),
						Start: 1,
						End:   2,
						Data:  []byte("f"),
					},
				},
				Data: []byte("tf"),
			},
			nil,
			DefaultParser,
		},
		{
			"t f",
			NewChain(
				"chain",
				NewTerminal("t", "t"),
				NewRepetitionTimesVariadic(
					"space",
					0,
					NewTerminal("space", " "),
				),
				NewTerminal("f", "f"),
			),
			&Tree{
				Rule: NewChain(
					"chain",
					NewTerminal("t", "t"),
					NewRepetitionTimesVariadic(
						"space",
						0,
						NewTerminal("space", " "),
					),
					NewTerminal("f", "f"),
				),
				Start: 0,
				End:   3,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("t", "t"),
						Start: 0,
						End:   1,
						Data:  []byte("t"),
					},
					{
						Rule: NewRepetitionTimesVariadic(
							"space",
							0,
							NewTerminal("space", " "),
						),
						Start: 1,
						End:   2,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("space", " "),
								Start: 1,
								End:   2,
								Data:  []byte(" "),
							},
						},
						Data: []byte(" "),
					},

					{
						Rule:  NewTerminal("f", "f"),
						Start: 2,
						End:   3,
						Data:  []byte("f"),
					},
				},
				Data: []byte("t f"),
			},
			nil,
			DefaultParser,
		},
		{
			// with 2 spaces
			"t  f",
			NewChain(
				"chain",
				NewTerminal("t", "t"),
				NewRepetitionTimesVariadic(
					"space",
					0,
					NewTerminal("space", " "),
				),
				NewTerminal("f", "f"),
			),
			&Tree{
				Rule: NewChain(
					"chain",
					NewTerminal("t", "t"),
					NewRepetitionTimesVariadic(
						"space",
						0,
						NewTerminal("space", " "),
					),
					NewTerminal("f", "f"),
				),
				Start: 0,
				End:   4,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("t", "t"),
						Start: 0,
						End:   1,
						Data:  []byte("t"),
					},
					{
						Rule: NewRepetitionTimesVariadic(
							"space",
							0,
							NewTerminal("space", " "),
						),
						Start: 1,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("space", " "),
								Start: 1,
								End:   2,
								Data:  []byte(" "),
							},
							{
								Rule:  NewTerminal("space", " "),
								Start: 2,
								End:   3,
								Data:  []byte(" "),
							},
						},
						Data: []byte("  "),
					},
					{
						Rule:  NewTerminal("f", "f"),
						Start: 3,
						End:   4,
						Data:  []byte("f"),
					},
				},
				Data: []byte("t  f"),
			},
			nil,
			DefaultParser,
		},
		{
			"123",
			NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			&Tree{
				Rule: NewRepetition(
					"numbers",
					NewEither(
						"number",
						NewTerminal("one", "1"),
						NewTerminal("two", "2"),
						NewTerminal("three", "3"),
					),
				),
				Data:  []byte("123"),
				Start: 0,
				End:   3,
				Childs: []*Tree{
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("1"),
						Start: 0,
						End:   1,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("one", "1"),
								Data:  []byte("1"),
								Start: 0,
								End:   1,
							},
						},
					},
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("2"),
						Start: 1,
						End:   2,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("two", "2"),
								Data:  []byte("2"),
								Start: 1,
								End:   2,
							},
						},
					},
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("3"),
						Start: 2,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("three", "3"),
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
		{
			"123",
			NewRepetitionTimesVariadic(
				"numbers",
				2,
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			&Tree{
				Rule: NewRepetitionTimesVariadic(
					"numbers",
					2,
					NewEither(
						"number",
						NewTerminal("one", "1"),
						NewTerminal("two", "2"),
						NewTerminal("three", "3"),
					),
				),
				Data:  []byte("123"),
				Start: 0,
				End:   3,
				Childs: []*Tree{
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("1"),
						Start: 0,
						End:   1,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("one", "1"),
								Data:  []byte("1"),
								Start: 0,
								End:   1,
							},
						},
					},
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("2"),
						Start: 1,
						End:   2,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("two", "2"),
								Data:  []byte("2"),
								Start: 1,
								End:   2,
							},
						},
					},
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Data:  []byte("3"),
						Start: 2,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("three", "3"),
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
		{
			"1234",
			NewRepetitionTimes(
				"numbers",
				3,
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("4"),
				4,
				NewRepetitionTimes(
					"numbers",
					3,
					NewEither(
						"number",
						NewTerminal("one", "1"),
						NewTerminal("two", "2"),
						NewTerminal("three", "3"),
					),
				),
			),
			DefaultParser,
		},
		{
			"123",
			NewRepetitionTimes(
				"numbers",
				2,
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("3"),
				3,
				NewRepetitionTimes(
					"numbers",
					2,
					NewEither(
						"number",
						NewTerminal("one", "1"),
						NewTerminal("two", "2"),
						NewTerminal("three", "3"),
					),
				),
			),
			DefaultParser,
		},

		// mixed

		{
			"foobar",
			NewChain(
				"chain of foo bar",
				NewEither(
					"foo and bar",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
				NewEither(
					"foo and bar",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
			),
			&Tree{
				Rule: NewChain(
					"chain of foo bar",
					NewEither(
						"foo and bar",
						NewTerminal("foo", "foo"),
						NewTerminal("bar", "bar"),
					),
					NewEither(
						"foo and bar",
						NewTerminal("foo", "foo"),
						NewTerminal("bar", "bar"),
					),
				),
				Data:  []byte("foobar"),
				Start: 0,
				End:   6,
				Childs: []*Tree{
					{
						Rule: NewEither(
							"foo and bar",
							NewTerminal("foo", "foo"),
							NewTerminal("bar", "bar"),
						),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("foo", "foo"),
								Data:  []byte("foo"),
								Start: 0,
								End:   3,
							},
						},
					},
					{
						Rule: NewEither(
							"foo and bar",
							NewTerminal("foo", "foo"),
							NewTerminal("bar", "bar"),
						),
						Data:  []byte("bar"),
						Start: 3,
						End:   6,
						Childs: []*Tree{
							{
								Rule:  NewTerminal("bar", "bar"),
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
				"foo func",
				NewTerminal("foo", "foo"),
				NewTerminal("left bracket", "("),
				NewRepetition(
					"numbers",
					NewEither(
						"number",
						NewTerminal("four", "4"),
						NewTerminal("three", "3"),
						NewTerminal("two", "2"),
						NewTerminal("one", "1"),
					),
				),
				NewTerminal("right bracket", ")"),
			),
			&Tree{
				Rule: NewChain(
					"foo func",
					NewTerminal("foo", "foo"),
					NewTerminal("left bracket", "("),
					NewRepetition(
						"numbers",
						NewEither(
							"number",
							NewTerminal("four", "4"),
							NewTerminal("three", "3"),
							NewTerminal("two", "2"),
							NewTerminal("one", "1"),
						),
					),
					NewTerminal("right bracket", ")"),
				),
				Data:  []byte("foo(1234)"),
				Start: 0,
				End:   9,
				Childs: []*Tree{
					{
						Rule:  NewTerminal("foo", "foo"),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
					},
					{
						Rule:  NewTerminal("left bracket", "("),
						Data:  []byte("("),
						Start: 3,
						End:   4,
					},
					{
						Rule: NewRepetition(
							"numbers",
							NewEither(
								"number",
								NewTerminal("four", "4"),
								NewTerminal("three", "3"),
								NewTerminal("two", "2"),
								NewTerminal("one", "1"),
							),
						),
						Data:  []byte("1234"),
						Start: 4,
						End:   8,
						Childs: []*Tree{
							{
								Rule: NewEither(
									"number",
									NewTerminal("four", "4"),
									NewTerminal("three", "3"),
									NewTerminal("two", "2"),
									NewTerminal("one", "1"),
								),
								Data:  []byte("1"),
								Start: 4,
								End:   5,
								Childs: []*Tree{
									{
										Rule:  NewTerminal("one", "1"),
										Data:  []byte("1"),
										Start: 4,
										End:   5,
									},
								},
							},
							{
								Rule: NewEither(
									"number",
									NewTerminal("four", "4"),
									NewTerminal("three", "3"),
									NewTerminal("two", "2"),
									NewTerminal("one", "1"),
								),
								Data:  []byte("2"),
								Start: 5,
								End:   6,
								Childs: []*Tree{
									{
										Rule:  NewTerminal("two", "2"),
										Data:  []byte("2"),
										Start: 5,
										End:   6,
									},
								},
							},
							{
								Rule: NewEither(
									"number",
									NewTerminal("four", "4"),
									NewTerminal("three", "3"),
									NewTerminal("two", "2"),
									NewTerminal("one", "1"),
								),
								Data:  []byte("3"),
								Start: 6,
								End:   7,
								Childs: []*Tree{
									{
										Rule:  NewTerminal("three", "3"),
										Data:  []byte("3"),
										Start: 6,
										End:   7,
									},
								},
							},
							{
								Rule: NewEither(
									"number",
									NewTerminal("four", "4"),
									NewTerminal("three", "3"),
									NewTerminal("two", "2"),
									NewTerminal("one", "1"),
								),
								Data:  []byte("4"),
								Start: 7,
								End:   8,
								Childs: []*Tree{
									{
										Rule:  NewTerminal("four", "4"),
										Data:  []byte("4"),
										Start: 7,
										End:   8,
									},
								},
							},
						},
					},
					{
						Rule:  NewTerminal("right bracket", ")"),
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
