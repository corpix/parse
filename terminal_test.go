package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestTerminalName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewTerminal("sample terminal", "hello"),
			"sample terminal",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
	}
}

func TestTerminalShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewTerminal("sample terminal", "hello"),
			"none",
			"*parse.Terminal(name: sample terminal, value: hello)(none)",
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

func TestTerminalString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewTerminal("sample terminal", "hello"),
			"*parse.Terminal(name: sample terminal, value: hello)()",
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

func TestTerminalGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewTerminal("sample terminal", "hello"),
			Treers(nil),
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

func TestTerminalGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewTerminal("sample terminal", "hello"),
			RuleParameters{
				"name":  "sample terminal",
				"value": "hello",
			},
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

func TestTerminalIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		true,
		NewTerminal(
			"foo",
			"foo",
		).IsFinite(),
		"Terminal is a finite entity",
	)
}

func TestTerminal(t *testing.T) {
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
			NewTerminal("empty", ""),
			nil,
			NewErrEmptyRule(
				NewTerminal("empty", ""),
				nil,
			),
			DefaultParser,
		},
		{
			"",
			NewTerminal("foo", "foo"),
			nil,
			NewErrUnexpectedEOF(
				0,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"o",
			NewTerminal("o", "oó"),
			nil,
			NewErrUnexpectedEOF(
				0,
				NewTerminal("o", "oó"),
			),
			DefaultParser,
		},
		{
			"oo",
			NewTerminal("o", "oó"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("oo")),
				0,
				NewTerminal("o", "oó"),
			),
			DefaultParser,
		},
		{
			"bar",
			NewTerminal("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("bar")),
				0,
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
				3,
				NewTerminal("foo", "foo"),
			),
			DefaultParser,
		},

		// success

		{
			"foo",
			NewTerminal("foo", "foo"),
			&Tree{
				Rule:  NewTerminal("foo", "foo"),
				Location: &Location{},
				Region: &Region{
					Start: 0,
					End:   3,
				},
				Data:  []byte("foo"),
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
