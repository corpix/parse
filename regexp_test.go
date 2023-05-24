package parse

import (
	"fmt"
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
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			msg := spew.Sdump(k, sample)
			assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
		})
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
				NewRegexp("foo", "foo"),
				&Location{Path: DefaultParserPath},
				ShowInput([]byte("")),
			),
			DefaultParser,
		},
		{
			"bar",
			NewRegexp("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				NewRegexp("foo", "foo"),
				&Location{Path: DefaultParserPath},
				ShowInput([]byte("bar")),
			),
			DefaultParser,
		},
		{
			"foobar",
			NewRegexp("foo", "foo"),
			nil,
			NewErrUnexpectedToken(
				NewRegexp("foo", "foo"),
				&Location{
					Path:     DefaultParserPath,
					Position: 3,
					Column:   3,
				},
				ShowInput([]byte("bar")),
			),
			DefaultParser,
		},

		// success

		{
			"",
			NewRegexp("empty", ""),
			&Tree{
				Rule:     NewRegexp("empty", ""),
				Location: &Location{Path: DefaultParserPath},
				Region: &Region{
					Start: 0,
					End:   0,
				},
				Data: []byte(""),
			},
			nil,
			DefaultParser,
		},
		{
			"foo",
			NewRegexp("foo", "foo"),
			&Tree{
				Rule:     NewRegexp("foo", "foo"),
				Location: &Location{Path: DefaultParserPath},
				Region: &Region{
					Start: 0,
					End:   3,
				},
				Data: []byte("foo"),
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
