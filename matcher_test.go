package parse

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
