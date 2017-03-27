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

func TestRuleString(t *testing.T) {
	samples := []struct {
		grammar Rule
		result  string
	}{
		{
			NewChain(
				"foo-bar",
				NewTerminal("foo", "Foo"),
				NewTerminal("bar", "Bar"),
			),
			"*parse.Chain(ID: foo-bar)(\n  *parse.Terminal(ID: foo, Value: [70 111 111])(), \n  *parse.Terminal(ID: bar, Value: [66 97 114])()\n)",
		},
		{
			NewChain(
				"foo-bar",
				NewEither(
					"foo",
					NewRepetition(
						"foo",
						NewTerminal("foo", "foo"),
					),
				),
				NewTerminal("bar", "Bar"),
			),
			"*parse.Chain(ID: foo-bar)(\n  *parse.Either(ID: foo)(\n    *parse.Repetition(ID: foo, Times: 1, Variadic: true)(*parse.Terminal(ID: foo, Value: [102 111 111])())\n  ), \n  *parse.Terminal(ID: bar, Value: [66 97 114])()\n)",
		},
		{
			func() Rule {
				foo := NewEither(
					"foo",
					NewRepetition(
						"foo",
						NewTerminal("foo", "foo"),
					),
				)
				foo.Add(foo)

				return NewChain(
					"foo-bar",
					foo,
					NewTerminal("bar", "Bar"),
				)
			}(),
			"*parse.Chain(ID: foo-bar)(\n  *parse.Either(ID: foo)(\n    *parse.Repetition(ID: foo, Times: 1, Variadic: true)(*parse.Terminal(ID: foo, Value: [102 111 111])()), \n    *parse.Either(ID: foo)(<circular>)\n  ), \n  *parse.Terminal(ID: bar, Value: [66 97 114])()\n)",
		},
		{
			NewChain(
				"nil",
				nil,
			),
			"*parse.Chain(ID: nil)(<nil>)",
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.Equal(
			t,
			sample.result,
			sample.grammar.String(),
			msg,
		)
		assert.Equal(
			t,
			sample.result,
			RulePrettyString(sample.grammar),
			msg,
		)
	}
}
