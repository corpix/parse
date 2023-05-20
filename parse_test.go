package parse

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
			nil,
			nil,
			NewErrEmptyRule(
				Rule(nil),
				nil,
			),
			DefaultParser,
		},
		{
			"",
			newTestRuleFinite("unsupported"),
			nil,
			NewErrUnsupportedRule(
				newTestRuleFinite("unsupported"),
			),
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
			NewErrNestingTooDeep(4, 1),
			NewParser(3),
		},

		// success

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
