package parse

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestRepetitionName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("foo")),
			"sample repetition",
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
		})
	}
}

func TestRepetitionShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("foo")),
			"none",
			"*parse.Repetition(name: sample repetition, times: 1, variadic: true)(none)",
		},
		{
			NewRepetition("another sample repetition", newTestRuleFinite("inner")),
			"different",
			"*parse.Repetition(name: another sample repetition, times: 1, variadic: true)(different)",
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(
				t,
				sample.show,
				sample.rule.Show(sample.childs),
				msg,
			)
		})
	}
}

func TestRepetitionString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			"*parse.Repetition(name: sample repetition, times: 1, variadic: true)(\n  *parse.testRuleFinite(name: inner)()\n)",
		},
		{
			NewRepetition("another sample repetition", newTestRuleNonFinite("inner")),
			"*parse.Repetition(name: another sample repetition, times: 1, variadic: true)(\n  *parse.testRuleNonFinite(name: inner)()\n)",
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(
				t,
				sample.stringified,
				sample.rule.String(),
				msg,
			)
		})
	}
}

func TestRepetitionGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			Treers{newTestRuleFinite("inner")},
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(
				t,
				sample.childs,
				sample.rule.GetChilds(),
				msg,
			)
		})
	}
}

func TestRepetitionGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			RuleParameters{
				"name":     "sample repetition",
				"times":    1,
				"variadic": true,
			},
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(
				t,
				sample.params,
				sample.rule.GetParameters(),
				msg,
			)
		})
	}
}

func TestRepetitionIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewRepetition(
			"foo",
			newTestRuleFinite("inner"),
		).IsFinite(),
		"Repetition is not a finite entity",
	)
}

func TestRepetition(t *testing.T) {
	samples := []struct {
		text   string
		rule   Rule
		tree   *Tree
		err    error
		parser *Parser
	}{

		// errors

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
				3,
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

		// success

		{
			"",
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("t"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("terminal", "t"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("t"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("f"),
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
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("f"),
						Childs: []*Tree{
							{
								Rule:     NewTerminal("terminal", "f"),
								Location: &Location{Depth: 2},
								Region: &Region{
									Start: 0,
									End:   1,
								},
								Data: []byte("f"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   2,
				},
				Childs: []*Tree{
					{
						Rule:     NewTerminal("t", "t"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("t"),
					},
					{
						Rule: NewTerminal("f", "f"),
						Location: &Location{
							Position: 1,
							Depth:    1,
						},
						Region: &Region{
							Start: 1,
							End:   2,
						},
						Data: []byte("f"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   3,
				},
				Childs: []*Tree{
					{
						Rule:     NewTerminal("t", "t"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("t"),
					},
					{
						Rule: NewRepetitionTimesVariadic(
							"space",
							0,
							NewTerminal("space", " "),
						),
						Location: &Location{
							Position: 1,
							Depth:    1,
						},
						Region: &Region{
							Start: 1,
							End:   2,
						},
						Childs: []*Tree{
							{
								Rule: NewTerminal("space", " "),
								Location: &Location{
									Position: 1,
									Depth:    2,
								},
								Region: &Region{
									Start: 1,
									End:   2,
								},
								Data: []byte(" "),
							},
						},
						Data: []byte(" "),
					},
					{
						Rule: NewTerminal("f", "f"),
						Location: &Location{
							Position: 2,
							Depth:    1,
						},
						Region: &Region{
							Start: 2,
							End:   3,
						},
						Data: []byte("f"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   4,
				},
				Childs: []*Tree{
					{
						Rule:     NewTerminal("t", "t"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("t"),
					},
					{
						Rule: NewRepetitionTimesVariadic(
							"space",
							0,
							NewTerminal("space", " "),
						),
						Location: &Location{
							Position: 1,
							Depth:    1,
						},
						Region: &Region{
							Start: 1,
							End:   3,
						},
						Childs: []*Tree{
							{
								Rule: NewTerminal("space", " "),
								Location: &Location{
									Position: 1,
									Depth:    2,
								},
								Region: &Region{
									Start: 1,
									End:   2,
								},
								Data: []byte(" "),
							},
							{
								Rule: NewTerminal("space", " "),
								Location: &Location{
									Position: 2,
									Depth:    2,
								},
								Region: &Region{
									Start: 2,
									End:   3,
								},
								Data: []byte(" "),
							},
						},
						Data: []byte("  "),
					},
					{
						Rule: NewTerminal("f", "f"),
						Location: &Location{
							Position: 3,
							Depth:    1,
						},
						Region: &Region{
							Start: 3,
							End:   4,
						},
						Data: []byte("f"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   3,
				},
				Data: []byte("123"),
				Childs: []*Tree{
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("1"),
						Childs: []*Tree{
							{
								Rule:     NewTerminal("one", "1"),
								Location: &Location{Depth: 2},
								Region: &Region{
									Start: 0,
									End:   1,
								},
								Data: []byte("1"),
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
						Location: &Location{
							Position: 1,
							Depth:    1,
						},
						Region: &Region{
							Start: 1,
							End:   2,
						},
						Data: []byte("2"),
						Childs: []*Tree{
							{
								Rule: NewTerminal("two", "2"),
								Location: &Location{
									Position: 1,
									Depth:    2,
								},
								Region: &Region{
									Start: 1,
									End:   2,
								},
								Data: []byte("2"),
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
						Location: &Location{
							Position: 2,
							Depth:    1,
						},
						Region: &Region{
							Start: 2,
							End:   3,
						},
						Data: []byte("3"),
						Childs: []*Tree{
							{
								Rule: NewTerminal("three", "3"),
								Location: &Location{
									Position: 2,
									Depth:    2,
								},
								Region: &Region{
									Start: 2,
									End:   3,
								},
								Data: []byte("3"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   3,
				},
				Data: []byte("123"),
				Childs: []*Tree{
					{
						Rule: NewEither(
							"number",
							NewTerminal("one", "1"),
							NewTerminal("two", "2"),
							NewTerminal("three", "3"),
						),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("1"),
						Childs: []*Tree{
							{
								Rule:     NewTerminal("one", "1"),
								Location: &Location{Depth: 2},
								Region: &Region{
									Start: 0,
									End:   1,
								},
								Data: []byte("1"),
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
						Location: &Location{
							Position: 1,
							Depth:    1,
						},
						Region: &Region{
							Start: 1,
							End:   2,
						},
						Data: []byte("2"),
						Childs: []*Tree{
							{
								Rule: NewTerminal("two", "2"),
								Location: &Location{
									Position: 1,
									Depth:    2,
								},
								Region: &Region{
									Start: 1,
									End:   2,
								},
								Data: []byte("2"),
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
						Location: &Location{
							Position: 2,
							Depth:    1,
						},
						Region: &Region{
							Start: 2,
							End:   3,
						},
						Data: []byte("3"),
						Childs: []*Tree{
							{
								Rule: NewTerminal("three", "3"),
								Location: &Location{
									Position: 2,
									Depth:    2,
								},
								Region: &Region{
									Start: 2,
									End:   3,
								},
								Data: []byte("3"),
							},
						},
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
