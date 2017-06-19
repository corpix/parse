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

func TestRepetitionName(t *testing.T) {
	samples := []struct {
		rule Rule
		name string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("foo")),
			"sample repetition",
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)
		assert.EqualValues(t, sample.name, sample.rule.Name(), msg)
	}
}

func TestRepetitionShow(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs string
		show   string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("foo")),
			"none",
			"*parse.Repetition(name: sample repetition)(none)",
		},
		{
			NewRepetition("another sample repetition", newTestRuleFinite("inner")),
			"different",
			"*parse.Repetition(name: another sample repetition)(different)",
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

func TestRepetitionString(t *testing.T) {
	samples := []struct {
		rule        Rule
		stringified string
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			"*parse.Repetition(name: sample repetition)(*parse.testRuleFinite(name: inner)())",
		},
		{
			NewRepetition("another sample repetition", newTestRuleNonFinite("inner")),
			"*parse.Repetition(name: another sample repetition)(*parse.testRuleNonFinite(name: inner)())",
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

func TestRepetitionGetChilds(t *testing.T) {
	samples := []struct {
		rule   Rule
		childs Treers
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			Treers{newTestRuleFinite("inner")},
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

func TestRepetitionGetParameters(t *testing.T) {
	samples := []struct {
		rule   Rule
		params RuleParameters
	}{
		{
			NewRepetition("sample repetition", newTestRuleFinite("inner")),
			RuleParameters{"name": "sample repetition"},
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

func TestRepetitionIsFinite(t *testing.T) {
	assert.EqualValues(
		t,
		false,
		NewRepetition(
			"foo",
			newTestRuleFinite("inner"),
		).IsFinite(),
		"Repetition is not a finite entity",
	)
}
