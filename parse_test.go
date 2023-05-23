package parse

import (
	"fmt"
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
			"faaaaaaaabfaaaabfaaaaaaaaaaaaaabfab",
			NewRepetition(
				"fabs",
				NewRegexp("fab", "fa+b"),
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
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)

			input := make([]byte, len(sample.input))
			copy(input, []byte(sample.input))
			_, err := Parse(sample.grammar, input)

			assert.Equal(t, []byte(sample.input), input, msg)
			assert.Equal(t, sample.err, err, msg)
		})
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
			nil,
			nil,
			NewErrEmptyRule(Rule(nil), nil),
			DefaultParser,
		},
		{
			"foo",
			NewRepetition(
				"deep",
				NewRepetition(
					"deep",
					NewRepetition(
						"deep",
						NewRepetition(
							"deep",
							NewTerminal("foo", "foo"),
						),
					),
				),
			),
			nil,
			NewErrNestingTooDeep(4, 0),
			NewParser(ParserOptionMaxDepth(3)),
		},

		// success

		{
			"",
			newTestRuleFinite("unsupported"),
			&Tree{
				Rule:     newTestRuleFinite("unsupported"),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   0,
				},
				Data:   []byte{},
				Childs: []*Tree{},
			},
			nil,
			DefaultParser,
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   6,
				},
				Data: []byte("foobar"),
				Childs: []*Tree{
					{
						Rule: NewEither(
							"foo and bar",
							NewTerminal("foo", "foo"),
							NewTerminal("bar", "bar"),
						),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   3,
						},
						Data: []byte("foo"),
						Childs: []*Tree{
							{
								Rule:     NewTerminal("foo", "foo"),
								Location: &Location{Depth: 2},
								Region: &Region{
									Start: 0,
									End:   3,
								},
								Data: []byte("foo"),
							},
						},
					},
					{
						Rule: NewEither(
							"foo and bar",
							NewTerminal("foo", "foo"),
							NewTerminal("bar", "bar"),
						),
						Location: &Location{
							Position: 3,
							Column:   3,
							Depth:    1,
						},
						Region: &Region{
							Start: 3,
							End:   6,
						},
						Data: []byte("bar"),
						Childs: []*Tree{
							{
								Rule: NewTerminal("bar", "bar"),
								Location: &Location{
									Position: 3,
									Column:   3,
									Depth:    2,
								},
								Region: &Region{
									Start: 3,
									End:   6,
								},
								Data: []byte("bar"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   9,
				},
				Data: []byte("foo(1234)"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("foo", "foo"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   3,
						},
						Data: []byte("foo"),
					},
					{
						Rule: NewTerminal("left bracket", "("),
						Location: &Location{
							Position: 3,
							Column:   3,
							Depth:    1,
						},
						Region: &Region{
							Start: 3,
							End:   4,
						},
						Data: []byte("("),
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
						Location: &Location{
							Position: 4,
							Column:   4,
							Depth:    1,
						},
						Region: &Region{
							Start: 4,
							End:   8,
						},
						Data: []byte("1234"),
						Childs: []*Tree{
							{
								Rule: NewEither(
									"number",
									NewTerminal("four", "4"),
									NewTerminal("three", "3"),
									NewTerminal("two", "2"),
									NewTerminal("one", "1"),
								),
								Location: &Location{
									Position: 4,
									Column:   4,
									Depth:    2,
								},
								Region: &Region{
									Start: 4,
									End:   5,
								},
								Data: []byte("1"),
								Childs: []*Tree{
									{
										Rule: NewTerminal("one", "1"),
										Location: &Location{
											Position: 4,
											Column:   4,
											Depth:    3,
										},
										Region: &Region{
											Start: 4,
											End:   5,
										},
										Data: []byte("1"),
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
								Location: &Location{
									Position: 5,
									Column:   5,
									Depth:    2,
								},
								Region: &Region{
									Start: 5,
									End:   6,
								},
								Data: []byte("2"),
								Childs: []*Tree{
									{
										Rule: NewTerminal("two", "2"),
										Location: &Location{
											Position: 5,
											Column:   5,
											Depth:    3,
										},
										Region: &Region{
											Start: 5,
											End:   6,
										},
										Data: []byte("2"),
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
								Location: &Location{
									Position: 6,
									Column:   6,
									Depth:    2,
								},
								Region: &Region{
									Start: 6,
									End:   7,
								},
								Data: []byte("3"),
								Childs: []*Tree{
									{
										Rule: NewTerminal("three", "3"),
										Location: &Location{
											Position: 6,
											Column:   6,
											Depth:    3,
										},
										Region: &Region{
											Start: 6,
											End:   7,
										},
										Data: []byte("3"),
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
								Location: &Location{
									Position: 7,
									Column:   7,
									Depth:    2,
								},
								Region: &Region{
									Start: 7,
									End:   8,
								},
								Data: []byte("4"),
								Childs: []*Tree{
									{
										Rule: NewTerminal("four", "4"),
										Location: &Location{
											Position: 7,
											Column:   7,
											Depth:    3,
										},
										Region: &Region{
											Start: 7,
											End:   8,
										},
										Data: []byte("4"),
									},
								},
							},
						},
					},
					{
						Rule: NewTerminal("right bracket", ")"),
						Location: &Location{
							Position: 8,
							Column:   8,
							Depth:    1,
						},
						Region: &Region{
							Start: 8,
							End:   9,
						},
						Data: []byte(")"),
					},
				},
			},
			nil,
			DefaultParser,
		},
	}

	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
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
		})
	}
}

func TestParserLineBreaksLocate(t *testing.T) {
	samples := []struct {
		text   string
		parser *Parser
		locs   []*Location
	}{
		{
			"foo\nbar\nbaz",
			DefaultParser,
			[]*Location{
				{0, 0, 0, 0},
				{5, 1, 1, 0},
				{6, 1, 2, 0},
				{7, 1, 3, 0},
				{10, 2, 2, 0},
				{666, 2, 2, 0},
			},
		},
		{
			"foo\nbar\nbaz\r\nqux",
			DefaultParser,
			[]*Location{
				{5, 1, 1, 0},
				{6, 1, 2, 0},
				{7, 1, 3, 0},
				{10, 2, 2, 0},
				{11, 2, 3, 0},
				{12, 3, 0, 0},
				{13, 3, 0, 0},
			},
		},
		{
			"foo bar",
			DefaultParser,
			[]*Location{
				{5, 0, 5, 0},
				{6, 0, 6, 0},
				{7, 0, 6, 0},
				{10, 0, 6, 0},
			},
		},
	}

	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			regions := sample.parser.LineRegions([]byte(sample.text))
			sample.parser.LineIndex = regions
			for n, loc := range sample.locs {
				t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
					line, col := sample.parser.Locate(loc.Position)
					msg := spew.Sdump(loc)
					_ = msg
					//spew.Dump(regions)
					assert.Equal(
						t, loc.Line, line,
						// fmt.Sprintf(
						// 	"line check failed: %q -> %q, %s",
						// 	sample.text, sample.text[loc.Position:], msg,
						// ),
					)
					assert.Equal(
						t, loc.Column, col,
						// fmt.Sprintf(
						// 	"column check failed: %q -> %q, %s",
						// 	sample.text, sample.text[loc.Position:], msg,
						// ),
					)
				})
			}
		})
	}
}
