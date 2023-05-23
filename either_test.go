package parse

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestEitherName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewEither(
				"sample either",
				newTestRuleFinite("foo"),
				newTestRuleFinite("bar"),
			),
			"sample either",
		},
		{
			NewEither(
				"another sample either",
				newTestRuleFinite("foo"),
				newTestRuleFinite("bar"),
			),
			"another sample either",
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
		})
	}
}

func TestEitherShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewEither(
				"sample either",
				newTestRuleFinite("foo"),
				newTestRuleFinite("bar"),
			),
			"none",
			"*parse.Either(name: sample either)(none)",
		},
		{
			NewEither(
				"another sample either",
				newTestRuleFinite("foo"),
				newTestRuleFinite("bar"),
			),
			"different",
			"*parse.Either(name: another sample either)(different)",
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

func TestEitherString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewEither(
				"sample either of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"*parse.Either(name: sample either of two)(\n  *parse.testRuleFinite(name: inner1)(), \n  *parse.testRuleFinite(name: inner2)()\n)",
		},
		{
			NewEither(
				"sample either of three",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
			"*parse.Either(name: sample either of three)(\n  *parse.testRuleFinite(name: inner1)(), \n  *parse.testRuleFinite(name: inner2)(), \n  *parse.testRuleFinite(name: inner3)()\n)",
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

func TestEitherGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewEither(
				"another sample either",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			},
		},
		{
			NewEither(
				"sample either of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			},
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

func TestEitherGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewEither(
				"sample either",
				newTestRuleFinite("foo"),
				newTestRuleFinite("bar"),
			),
			RuleParameters{
				"name": "sample either",
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

func TestEitherIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewEither(
			"foo",
			newTestRuleFinite("foo"),
			newTestRuleFinite("bar"),
		).IsFinite(),
		"Either is not a finite entity",
	)
}

func TestEitherAdd(t *testing.T) {
	either := NewEither(
		"either",
		newTestRuleFinite("foo"),
		newTestRuleFinite("bar"),
	)
	either.Add(newTestRuleFinite("baz"))

	assert.EqualValues(
		t,
		NewEither(
			"either",
			newTestRuleFinite("foo"),
			newTestRuleFinite("bar"),
			newTestRuleFinite("baz"),
		),
		either,
	)

	either.Add(newTestRuleFinite("daz"))
	assert.EqualValues(
		t,
		NewEither(
			"either",
			newTestRuleFinite("foo"),
			newTestRuleFinite("bar"),
			newTestRuleFinite("baz"),
			newTestRuleFinite("daz"),
		),
		either,
	)
}

func TestEither(t *testing.T) {
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
			NewEither("number"),
			nil,
			NewErrEmptyRule(NewEither("number"), nil),
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
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
					NewTerminal("three", "3"),
				),
				&Location{},
				[]byte("4"),
			),
			DefaultParser,
		},
		{
			"12",
			NewEither(
				"number",
				NewTerminal("one", "1"),
				NewTerminal("two", "2"),
			),
			nil,
			NewErrUnexpectedToken(
				NewEither(
					"number",
					NewTerminal("one", "1"),
					NewTerminal("two", "2"),
				),
				&Location{
					Position: 1,
					Column:   1,
				},
				[]byte("2"),
			),
			DefaultParser,
		},

		// success

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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("1"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("one", "1"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("1"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("1"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("one", "1"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("1"),
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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("2"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("two", "2"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("2"),
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
			if sample.err == nil && err != nil {
				t.Error(err)
			} else {
				assert.EqualValues(t, sample.err, err, msg)
			}
			assert.EqualValues(t, sample.tree, tree, msg)

		})
	}
}

func TestASCIIRange(t *testing.T) {
	samples := []struct {
		text   string
		rule   Rule
		tree   *Tree
		err    error
		parser *Parser
	}{
		{
			"0",
			NewASCIIRange(
				"numbers",
				'0', '3',
			),
			&Tree{
				Rule: NewEither(
					"numbers",
					NewTerminal("0", "0"),
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
					NewTerminal("3", "3"),
				),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("0"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("0", "0"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("0"),
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"1",
			NewASCIIRange(
				"numbers",
				'0', '3',
			),
			&Tree{
				Rule: NewEither(
					"numbers",
					NewTerminal("0", "0"),
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
					NewTerminal("3", "3"),
				),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("1"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("1", "1"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("1"),
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"3",
			NewASCIIRange(
				"numbers",
				'0', '3',
			),
			&Tree{
				Rule: NewEither(
					"numbers",
					NewTerminal("0", "0"),
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
					NewTerminal("3", "3"),
				),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   1,
				},
				Data: []byte("3"),
				Childs: []*Tree{
					{
						Rule:     NewTerminal("3", "3"),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   1,
						},
						Data: []byte("3"),
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
			if sample.err == nil && err != nil {
				t.Error(err)
			} else {
				assert.EqualValues(t, sample.err, err, msg)
			}
			assert.EqualValues(t, sample.tree, tree, msg)
		})
	}
}
