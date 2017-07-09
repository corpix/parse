package parse


import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestEqualSlicesFoldPrefix(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		equal bool
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
			[]string{"foo", "bar"},
			true,
		},
		{
			[]string{"A", "a"},
			[]string{"a", "a"},
			true,
		},
		{
			[]string{"a"},
			[]string{"b"},
			false,
		},
		{
			[]string{"a", "a"},
			[]string{"a"},
			true,
		},
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b"},
			true,
		},
		{
			[]string{"a"},
			[]string{"a", "b"},
			false,
		},
		{
			[]string{"a", "c"},
			[]string{"a", "b"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.equal,
			EqualSlicesFoldPrefix(sample.a, sample.b),
			msg,
		)
	}
}

func TestEqualSlicesFoldSuffix(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		equal bool
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
			[]string{"foo", "bar"},
			true,
		},
		{
			[]string{"A", "a"},
			[]string{"a", "a"},
			true,
		},
		{
			[]string{"a"},
			[]string{"b"},
			false,
		},
		{
			[]string{"a", "a"},
			[]string{"a"},
			true,
		},
		{
			[]string{"a", "b", "c"},
			[]string{"b", "c"},
			true,
		},
		{
			[]string{"a"},
			[]string{"b", "a"},
			false,
		},
		{
			[]string{"a", "c"},
			[]string{"a", "b"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.equal,
			EqualSlicesFoldSuffix(sample.a, sample.b),
			msg,
		)
	}
}

func TestEqualSlicesFold(t *testing.T) {
	samples := []struct {
		a     []string
		b     []string
		equal bool
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
			[]string{"A", "a"},
			[]string{"a", "a"},
			true,
		},
		{
			[]string{"a"},
			[]string{"b"},
			false,
		},
		{
			[]string{"a", "a"},
			[]string{"a"},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.equal,
			EqualSlicesFold(sample.a, sample.b),
			msg,
		)
	}
}

func TestEqualSlicesFoldSome(t *testing.T) {
	samples := []struct {
		a     []string
		b     [][]string
		equal bool
	}{
		{
			[]string{},
			[][]string{},
			true,
		},
		{
			[]string{"foo", "bar"},
			[][]string{{"foo", "bar"}},
			true,
		},
		{
			[]string{"A", "a"},
			[][]string{{"a", "a"}},
			true,
		},
		{
			[]string{"a"},
			[][]string{{"b"}},
			false,
		},
		{
			[]string{"a", "a"},
			[][]string{{"a"}},
			false,
		},
		{
			[]string{"a", "a"},
			[][]string{{"a"}, {"a"}},
			false,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.EqualValues(
			t,
			sample.equal,
			EqualSlicesFoldSome(sample.a, sample.b...),
			msg,
		)
	}

}
