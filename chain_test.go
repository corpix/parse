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
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
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
			"*parse.Chain(name: sample chain)(*parse.testRuleFinite(name: inner1)(), \n*parse.testRuleFinite(name: inner2)())",
		},
		{
			NewChain(
				"another sample chain",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
				newTestRuleFinite("inner3"),
			),
			"*parse.Chain(name: another sample chain)(*parse.testRuleFinite(name: inner1)(), \n*parse.testRuleFinite(name: inner2)(), \n*parse.testRuleFinite(name: inner3)())",
		},
		{
			NewChain(
				"sample chain of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"*parse.Chain(name: sample chain of two)(*parse.testRuleFinite(name: inner1)(), \n*parse.testRuleFinite(name: inner2)())",
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
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.childs,
			sample.rule.GetChilds(),
			msg,
		)
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
		msg := spew.Sdump(k, sample)
		assert.EqualValues(
			t,
			sample.params,
			sample.rule.GetParameters(),
			msg,
		)
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
	}
}
