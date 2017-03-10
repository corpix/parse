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
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScan(t *testing.T) {
	samples := []struct {
		grammar Rule
		input   string
		result  []string
		err     error
	}{
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
				),
			),
			input: "12",
			result: []string{
				"numbers",
				"number",
				"number",
				"1",
				"2",
			},
			err: nil,
		},
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
				),
			),
			input: "1",
			result: []string{
				"numbers",
				"number",
				"1",
			},
			err: nil,
		},
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
					NewTerminal("3", "3"),
				),
			),
			input: "123",
			result: []string{
				"numbers",
				"number",
				"number",
				"number",
				"1",
				"2",
				"3",
			},
			err: nil,
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		tree, err := Parse(
			sample.grammar,
			[]byte(sample.input),
		)
		assert.Equal(t, sample.err, err, msg)

		buf := []string{}
		Scan(
			tree,
			func(tree *Tree) { buf = append(buf, tree.ID()) },
		)
		assert.EqualValues(t, sample.result, buf, msg)
	}
}

func TestScanAlterChilds(t *testing.T) {
	samples := []struct {
		grammar Rule
		input   string
		result  []string
		err     error
	}{
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
				),
			),
			input: "12",
			result: []string{
				"numbers",
				"number",
				"number",
				"1",
				"2",
			},
			err: nil,
		},
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
				),
			),
			input: "1",
			result: []string{
				"numbers",
				"number",
				"1",
			},
			err: nil,
		},
		{
			grammar: NewRepetition(
				"numbers",
				NewEither(
					"number",
					NewTerminal("1", "1"),
					NewTerminal("2", "2"),
					NewTerminal("3", "3"),
				),
			),
			input: "123",
			result: []string{
				"numbers",
				"number",
				"number",
				"number",
				"1",
				"2",
				"3",
			},
			err: nil,
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		tree, err := Parse(
			sample.grammar,
			[]byte(sample.input),
		)
		assert.Equal(t, sample.err, err, msg)

		buf := []string{}
		Scan(
			tree,
			func(tree *Tree) {
				buf = append(buf, tree.ID())
				tree.Childs = nil
			},
		)
		assert.EqualValues(t, sample.result, buf, msg)
	}
}
