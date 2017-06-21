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

		// success

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
