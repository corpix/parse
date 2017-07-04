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

type testMatcher []string

func (m testMatcher) Match(chain []string) bool {
	return EqualSlicesFold(m, chain)
}

//

func TestAllMatcher(t *testing.T) {
	samples := []struct {
		a     []Matcher
		b     []string
		match bool
	}{
		{
			[]Matcher{},
			[]string{"foo", "bar"},
			false,
		},
		{
			[]Matcher{},
			[]string{},
			false,
		},
		{
			[]Matcher{
				testMatcher([]string{"foo", "bar"}),
			},
			[]string{},
			false,
		},
		{
			[]Matcher{
				testMatcher([]string{"foo", "bar"}),
			},
			[]string{"foo", "bar"},
			true,
		},
		{
			[]Matcher{
				testMatcher([]string{"foo", "bar"}),
				testMatcher([]string{"bar", "foo"}),
			},
			[]string{"foo", "bar"},
			false,
		},
		{
			[]Matcher{
				testMatcher([]string{"bar", "foo"}),
				testMatcher([]string{"foo", "bar"}),
			},
			[]string{"foo", "bar"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewAllMatcher(sample.a).Match(sample.b),
			msg,
		)
	}
}

func TestSomeMatcher(t *testing.T) {
	samples := []struct {
		a     []Matcher
		b     []string
		match bool
	}{
		{
			[]Matcher{},
			[]string{"foo", "bar"},
			false,
		},
		{
			[]Matcher{},
			[]string{},
			false,
		},
		{
			[]Matcher{
				testMatcher([]string{"foo", "bar"}),
				testMatcher([]string{"bar", "foo"}),
			},
			[]string{},
			false,
		},
		{
			[]Matcher{
				testMatcher([]string{"bar", "foo"}),
				testMatcher([]string{"foo", "bar"}),
			},
			[]string{"foo", "bar"},
			true,
		},
		{
			[]Matcher{
				testMatcher([]string{"foo", "bar"}),
				testMatcher([]string{"bar", "foo"}),
			},
			[]string{"foo", "bar"},
			true,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewSomeMatcher(sample.a).Match(sample.b),
			msg,
		)
	}
}

func TestPrefixMatcher(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		match bool
	}{
		{
			[]string{},
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{"foo"},
			true,
		},
		{
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "baz"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewPrefixMatcher(sample.b).Match(sample.a),
			msg,
		)
	}
}

func TestSuffixMatcher(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		match bool
	}{
		{
			[]string{},
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{"bar"},
			true,
		},
		{
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "baz"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewSuffixMatcher(sample.b).Match(sample.a),
			msg,
		)
	}
}

func TestEqualMatcher(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		match bool
	}{
		{
			[]string{},
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{"foo", "bar"},
			true,
		},
		{
			[]string{"foo", "bar"},
			[]string{"foo", "baz"},
			false,
		},
		{
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "baz"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewEqualMatcher(sample.b).Match(sample.a),
			msg,
		)
	}
}

func TestLengthMatcher(t *testing.T) {
	samples := []struct {
		a     []string
		b     int
		match bool
	}{
		{
			[]string{},
			0,
			true,
		},
		{
			[]string{"foo", "bar"},
			2,
			true,
		},
		{
			[]string{"foo", "bar", "baz"},
			2,
			false,
		},
		{
			[]string{"foo", "bar"},
			3,
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.match,
			NewLengthMatcher(sample.b).Match(sample.a),
			msg,
		)
	}
}
