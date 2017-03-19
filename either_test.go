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

func TestEitherAdd(t *testing.T) {
	chain := NewEither("either")
	chain.Add(NewTerminal("terminal", "hello"))

	assert.EqualValues(
		t,
		NewEither(
			"either",
			NewTerminal("terminal", "hello"),
		),
		chain,
	)

	chain.Add(NewTerminal("otherTerminal", "hello"))
	assert.EqualValues(
		t,
		NewEither(
			"either",
			NewTerminal("terminal", "hello"),
			NewTerminal("otherTerminal", "hello"),
		),
		chain,
	)
}

func TestEitherString(t *testing.T) {
	samples := []struct {
		grammar Rule
		result  string
	}{
		{
			NewEither(
				"foo/bar",
				NewTerminal("foo", "Foo"),
				NewTerminal("bar", "Bar"),
			),
			"foo/bar{foo(Foo), bar(Bar)}",
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		assert.Equal(t, sample.result, sample.grammar.String(), msg)
	}
}
