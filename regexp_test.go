package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestRegexpName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewRegexp("sample regexp", "hello"),
			"sample regexp",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
	}
}

func TestRegexpShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewRegexp("sample regexp", "hello"),
			"none",
			"*parse.Regexp(expr: hello, name: sample regexp)(none)",
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

func TestRegexpString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewRegexp("sample regexp", "hello"),
			"*parse.Regexp(expr: hello, name: sample regexp)()",
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

func TestRegexpGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewRegexp("sample regexp", "hello"),
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

func TestRegexpGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewRegexp("sample regexp", "hello"),
			RuleParameters{
				"name": "sample regexp",
				"expr": "hello",
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

func TestRegexpIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		true,
		NewRegexp(
			"foo",
			"foo",
		).IsFinite(),
		"Regexp is a finite entity",
	)
}

func TestRegexp(t *testing.T) {
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
			NewRegexp("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("")),
				1,
				NewRegexp("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"bar",
			NewRegexp("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("bar")),
				1,
				NewRegexp("foo", "foo"),
			),
			DefaultParser,
		},
		{
			"foobar",
			NewRegexp("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				ShowInput([]byte("bar")),
				4,
				NewRegexp("foo", "foo"),
			),
			DefaultParser,
		},

		// success

		{
			"",
			NewRegexp("empty", ""),
			&Tree{
				Rule:  NewRegexp("empty", ""),
				Data:  []byte(""),
				Start: 0,
				End:   0,
			},
			nil,
			DefaultParser,
		},
		{
			"foo",
			NewRegexp("foo", "foo"),
			&Tree{
				Rule:  NewRegexp("foo", "foo"),
				Data:  []byte("foo"),
				Start: 0,
				End:   3,
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
