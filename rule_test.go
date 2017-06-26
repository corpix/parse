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

type testRuleFinite string

func (r *testRuleFinite) Name() string { return string(*r) }
func (r *testRuleFinite) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().Show(),
		childs,
	)
}
func (r *testRuleFinite) String() string                { return TreerString(r) }
func (r *testRuleFinite) GetChilds() Treers             { return nil }
func (r *testRuleFinite) GetParameters() RuleParameters { return RuleParameters{"name": string(*r)} }
func (r *testRuleFinite) IsFinite() bool                { return true }

func newTestRuleFinite(name string) *testRuleFinite {
	r := testRuleFinite(name)
	return &r
}

//

type testRuleNonFinite struct {
	Rule
	rules Treers
}

func (r *testRuleNonFinite) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().Show(),
		childs,
	)
}
func (r *testRuleNonFinite) String() string    { return TreerString(r) }
func (r *testRuleNonFinite) GetChilds() Treers { return r.rules }
func (r *testRuleNonFinite) IsFinite() bool    { return false }

func newTestRuleNonFinite(name string, rules ...Treer) *testRuleNonFinite {
	return &testRuleNonFinite{
		newTestRuleFinite(name),
		rules,
	}
}

//

func TestRuleShow(t *testing.T) {
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
			"*parse.Chain(name: foo-bar)(\n  *parse.Terminal(name: foo, value: Foo)(), \n  *parse.Terminal(name: bar, value: Bar)()\n)",
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
					NewRepetition(
						"bar",
						NewTerminal("bar", "bar"),
					),
				),
				NewTerminal("bar", "Bar"),
			),
			"*parse.Chain(name: foo-bar)(\n  *parse.Either(name: foo)(\n    *parse.Repetition(name: foo, times: 1, variadic: true)(\n      *parse.Terminal(name: foo, value: foo)()\n    ), \n    *parse.Repetition(name: bar, times: 1, variadic: true)(\n      *parse.Terminal(name: bar, value: bar)()\n    )\n  ), \n  *parse.Terminal(name: bar, value: Bar)()\n)",
		},
		{
			func() Rule {
				foo := NewEither(
					"foo",
					NewRepetition(
						"foo",
						NewTerminal("foo", "foo"),
					),
					NewRepetition(
						"bar",
						NewTerminal("bar", "bar"),
					),
				)
				foo.Add(foo)

				return NewChain(
					"foo-bar",
					foo,
					NewTerminal("bar", "Bar"),
				)
			}(),
			"*parse.Chain(name: foo-bar)(\n  *parse.Either(name: foo)(\n    *parse.Repetition(name: foo, times: 1, variadic: true)(\n      *parse.Terminal(name: foo, value: foo)()\n    ), \n    *parse.Repetition(name: bar, times: 1, variadic: true)(\n      *parse.Terminal(name: bar, value: bar)()\n    ), \n    *parse.Either(name: foo)(<circular>)\n  ), \n  *parse.Terminal(name: bar, value: Bar)()\n)",
		},
		{
			NewChain(
				"nil",
				nil,
				nil,
			),
			"*parse.Chain(name: nil)(\n  <nil>, \n  <nil>\n)",
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
			TreerString(sample.grammar),
			msg,
		)
	}
}
