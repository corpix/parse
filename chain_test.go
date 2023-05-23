package parse

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestChainName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewChain(
				"sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"sample chain",
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"another sample chain",
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
		})
	}
}

func TestChainShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewChain(
				"sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"none",
			"*parse.Chain(name: sample chain)(none)",
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"different",
			"*parse.Chain(name: another sample chain)(different)",
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

func TestChainString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewChain(
				"sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"*parse.Chain(name: sample chain)(\n  *parse.testRuleFinite(name: inner1)(), \n  *parse.testRuleFinite(name: inner2)()\n)",
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
			"*parse.Chain(name: another sample chain)(\n  *parse.testRuleFinite(name: inner1)(), \n  *parse.testRuleFinite(name: inner2)(), \n  *parse.testRuleFinite(name: inner3)()\n)",
		},
		{
			NewChain(
				"sample chain of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"*parse.Chain(name: sample chain of two)(\n  *parse.testRuleFinite(name: inner1)(), \n  *parse.testRuleFinite(name: inner2)()\n)",
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

func TestChainGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewChain(
				"sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			},
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			},
		},
		{
			NewChain(
				"sample long chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
				newTestRuleFinite("inner4"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
				newTestRuleFinite("inner4"),
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

func TestChainGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewChain(
				"sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			RuleParameters{"name": "sample chain"},
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
			RuleParameters{"name": "another sample chain"},
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

func TestChainIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewChain(
			"foo",
			newTestRuleFinite("inner1"),
			newTestRuleFinite("inner2"),
		).IsFinite(),
		"Chain is not a finite entity",
	)
}

func TestChainAdd(t *testing.T) {
	samples := []struct {
		chainFactory func() *Chain
		add          Rules
		result       Rule
	}{
		{
			func() *Chain {
				return NewChain(
					"chain",
					newTestRuleFinite("inner1"),
					newTestRuleFinite("inner2"),
				)
			},
			Rules{newTestRuleFinite("inner3")},
			NewChain(
				"chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
		},
		{
			func() *Chain {
				return NewChain(
					"chain",
					newTestRuleFinite("inner1"),
					newTestRuleFinite("inner2"),
					newTestRuleFinite("inner3"),
				)
			},
			Rules{
				newTestRuleFinite("inner4"),
				newTestRuleFinite("inner5"),
			},
			NewChain(
				"chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
				newTestRuleFinite("inner4"),
				newTestRuleFinite("inner5"),
			),
		},
	}

	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)

			rule := sample.chainFactory()
			for _, inner := range sample.add {
				rule.Add(inner)
			}

			assert.EqualValues(
				t,
				rule,
				sample.result,
				msg,
			)

		})
	}
}

func TestChain(t *testing.T) {
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
			NewChain("some"),
			nil,
			NewErrEmptyRule(NewChain("some"), nil),
			DefaultParser,
		},
		{
			"",
			NewChain(
				"some",
				NewTerminal("foo", "foo"),
				NewTerminal("bar", "bar"),
			),
			nil,
			NewErrUnexpectedEOF(
				NewTerminal("foo", "foo"),
				&Location{Depth: 1},
			),
			DefaultParser,
		},
		{
			"bar",
			NewChain(
				"foo",
				NewTerminal("foo", "foo"),
				NewTerminal("bar", "bar"),
			),
			nil,
			NewErrUnexpectedToken(
				NewTerminal("foo", "foo"),
				&Location{Depth: 1},
				[]byte("bar"),
			),
			DefaultParser,
		},
		{
			"foobarbaz",
			NewChain(
				"foo",
				NewTerminal("foo", "foo"),
				NewTerminal("bar", "bar"),
			),
			nil,
			NewErrUnexpectedToken(
				NewChain(
					"foo",
					NewTerminal("foo", "foo"),
					NewTerminal("bar", "bar"),
				),
				&Location{Position: 6, Line: 0, Column: 6},
				[]byte("baz"),
			),
			DefaultParser,
		},

		// success

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
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   7,
				},
				Data: []byte("foo bar"),
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
						Rule: NewTerminal("space", " "),
						Location: &Location{
							Position: 3,
							Column:   3,
							Depth:    1,
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
							Column:   4,
							Depth:    1,
						},
						Region: &Region{
							Start: 4,
							End:   7,
						},
						Data: []byte("bar"),
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
