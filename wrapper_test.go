package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestWrapperName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewWrapper(
				"sample wrapper",
				newTestRuleFinite("inner"),
			),
			"sample wrapper",
		},
		{
			NewWrapper(
				"another sample wrapper",
				newTestRuleFinite("inner"),
			),
			"another sample wrapper",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
	}
}

func TestWrapperShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewWrapper(
				"sample wrapper",
				newTestRuleFinite("inner"),
			),
			"none",
			"*parse.Wrapper(name: sample wrapper)(none)",
		},
		{
			NewWrapper(
				"another sample wrapper",
				newTestRuleFinite("inner"),
			),
			"different",
			"*parse.Wrapper(name: another sample wrapper)(different)",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.show,
			sample.rule.Show(sample.childs),
			msg,
		)
	}
}

func TestWrapperString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewWrapper(
				"sample wrapper",
				newTestRuleFinite("inner"),
			),
			"*parse.Wrapper(name: sample wrapper)(\n  *parse.testRuleFinite(name: inner)()\n)",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.stringified,
			sample.rule.String(),
			msg,
		)
	}
}

func TestWrapperGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewWrapper(
				"sample wrapper",
				newTestRuleFinite("inner"),
			),
			Treers{
				newTestRuleFinite("inner"),
			},
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.childs,
			sample.rule.GetChilds(),
			msg,
		)
	}
}

func TestWrapperGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewWrapper(
				"sample wrapper",
				newTestRuleFinite("inner"),
			),
			RuleParameters{"name": "sample wrapper"},
		},
		{
			NewWrapper(
				"another sample wrapper",
				newTestRuleFinite("inner"),
			),
			RuleParameters{"name": "another sample wrapper"},
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.params,
			sample.rule.GetParameters(),
			msg,
		)
	}
}

func TestWrapperIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewWrapper(
			"foo",
			newTestRuleFinite("inner"),
		).IsFinite(),
		"Wrapper is not a finite entity",
	)
}

func TestWrapper(t *testing.T) {
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
			NewWrapper(
				"some",
				nil,
			),
			nil,
			NewErrEmptyRule(
				NewWrapper(
					"some",
					nil,
				),
				nil,
			),
			DefaultParser,
		},
		{
			"",
			NewWrapper(
				"some",
				NewChain(
					"terminals",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
			),
			nil,
			NewErrUnexpectedEOF(
				0,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"bar",
			NewWrapper(
				"foo",
				NewChain(
					"terminals",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("bar"),
				0,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"foobarbaz",
			NewWrapper(
				"foo",
				NewChain(
					"terminals",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
			),
			nil,
			NewErrUnexpectedToken(
				[]byte("baz"),
				6,
				NewWrapper(
					"foo",
					NewChain(
						"terminals",
						NewTerminal("foo", "foo"),
						NewTerminal("bar", "bar"),
					),
				),
			),
			DefaultParser,
		},

		// success

		{
			"foo bar",
			NewWrapper(
				"wrapper of foo bar",
				NewChain(
					"terminals",
					NewTerminal("foo", "foo"),
					NewTerminal("space", " "),
					NewTerminal("bar", "bar"),
				),
			),
			&Tree{
				Rule: NewWrapper(
					"wrapper of foo bar",
					NewChain(
						"terminals",
						NewTerminal("foo", "foo"),
						NewTerminal("space", " "),
						NewTerminal("bar", "bar"),
					),
				),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   7,
				},
				Data: []byte("foo bar"),
				Childs: []*Tree{
					{
						Rule: NewChain(
							"terminals",
							NewTerminal("foo", "foo"),
							NewTerminal("space", " "),
							NewTerminal("bar", "bar"),
						),
						Location: &Location{Depth: 1},
						Region: &Region{
							Start: 0,
							End:   7,
						},
						Data: []byte("foo bar"),
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
							{
								Rule: NewTerminal("space", " "),
								Location: &Location{
									Position: 3,
									Depth:    2,
								},
								Region: &Region{
									Start: 3,
									End:   4,
								},
								Data: []byte(" "),
							},
							{
								Rule: NewTerminal("bar", "bar"),
								Location: &Location{
									Position: 4,
									Depth:    2,
								},
								Region: &Region{
									Start: 4,
									End:   7,
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
		if sample.err == nil && err != nil {
			t.Error(err)
		} else {
			assert.EqualValues(t, sample.err, err, msg)
		}
		assert.EqualValues(t, sample.tree, tree, msg)
	}
}
