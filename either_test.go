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

func TestEitherName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewEither("sample either"),
			"sample either",
		},
		{
			NewEither("another sample either", newTestRuleFinite("inner")),
			"another sample either",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
	}
}

func TestEitherShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewEither("sample either"),
			"none",
			"*parse.Either(name: sample either)(none)",
		},
		{
			NewEither("another sample either", newTestRuleFinite("inner")),
			"different",
			"*parse.Either(name: another sample either)(different)",
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

func TestEitherString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewEither("sample either"),
			"*parse.Either(name: sample either)()",
		},
		{
			NewEither("another sample either", newTestRuleFinite("inner")),
			"*parse.Either(name: another sample either)(*parse.testRuleFinite(name: inner)())",
		},
		{
			NewEither(
				"sample either of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			"*parse.Either(name: sample either of two)(*parse.testRuleFinite(name: inner1)(), \n*parse.testRuleFinite(name: inner2)())",
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

func TestEitherGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewEither("sample either"),
			Treers{},
		},
		{
			NewEither("another sample either", newTestRuleFinite("inner")),
			Treers{newTestRuleFinite("inner")},
		},
		{
			NewEither(
				"sample either of two",
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
			),
			Treers{
				newTestRuleFinite("inner1"),
				newTestRuleFinite("inner2"),
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

func TestEitherGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewEither("sample either"),
			RuleParameters{"name": "sample either"},
		},
		{
			NewEither("another sample either", newTestRuleFinite("inner")),
			RuleParameters{"name": "another sample either"},
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

func TestEitherIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewEither("foo").IsFinite(),
		"Either is not a finite entity",
	)
}

func TestEitherAdd(t *testing.T) {
	either := NewEither("either")
	either.Add(newTestRuleFinite("terminal"))

	assert.EqualValues(
		t,
		NewEither(
			"either",
			newTestRuleFinite("terminal"),
		),
		either,
	)

	either.Add(newTestRuleFinite("otherTerminal"))
	assert.EqualValues(
		t,
		NewEither(
			"either",
			newTestRuleFinite("terminal"),
			newTestRuleFinite("otherTerminal"),
		),
		either,
	)
}
