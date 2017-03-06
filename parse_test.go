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
package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	samples := []struct {
		text   string
		syntax Rule
		tree   *Tree
		err    error
		parser *Parser
	}{
		{
			"",
			Terminal(""),
			&Tree{
				Rule:  Terminal(""),
				Data:  []byte{},
				Start: 0,
				End:   0,
			},
			nil,
			DefaultParser,
		},
		{
			"foo",
			Terminal("foo"),
			&Tree{
				Rule:  Terminal("foo"),
				Data:  []byte("foo"),
				Start: 0,
				End:   3,
			},
			nil,
			DefaultParser,
		},
		{
			"bar",
			Terminal("foo"),
			nil,
			NewErrUnexpectedToken([]byte("b"), 1),
			DefaultParser,
		},
		{
			"foobar",
			Terminal("foo"),
			nil,
			NewErrUnexpectedToken([]byte("b"), 4),
			DefaultParser,
		},
		{
			"foo bar",
			Chain{
				Terminal("foo"),
				Terminal(" "),
				Terminal("bar"),
			},
			&Tree{
				Rule: Chain{
					Terminal("foo"),
					Terminal(" "),
					Terminal("bar"),
				},
				Data:  []byte("foo bar"),
				Start: 0,
				End:   7,
				Childs: []*Tree{
					{
						Rule:  Terminal("foo"),
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
					},
					{
						Rule:  Terminal(" "),
						Data:  []byte(" "),
						Start: 3,
						End:   4,
					},
					{
						Rule:  Terminal("bar"),
						Data:  []byte("bar"),
						Start: 4,
						End:   7,
					},
				},
			},
			nil,
			DefaultParser,
		},
		{
			"foo",
			Chain{Chain{Chain{Chain{Terminal("foo")}}}},
			nil,
			NewErrNestingTooDeep(4, 1),
			NewParser(3),
		},

		//

		{
			"foobar",
			Chain{
				Either{Terminal("foo"), Terminal("bar")},
				Either{Terminal("foo"), Terminal("bar")},
			},
			&Tree{
				Rule: Chain{
					Either{Terminal("foo"), Terminal("bar")},
					Either{Terminal("foo"), Terminal("bar")},
				},
				Data:  []byte("foobar"),
				Start: 0,
				End:   6,
				Childs: []*Tree{
					{
						Rule: Either{
							Terminal("foo"),
							Terminal("bar"),
						},
						Data:  []byte("foo"),
						Start: 0,
						End:   3,
						Childs: []*Tree{
							{
								Rule:  Terminal("foo"),
								Data:  []byte("foo"),
								Start: 0,
								End:   3,
							},
						},
					},
					{
						Rule: Either{
							Terminal("foo"),
							Terminal("bar"),
						},
						Data:  []byte("bar"),
						Start: 3,
						End:   6,
						Childs: []*Tree{
							{
								Rule:  Terminal("bar"),
								Data:  []byte("bar"),
								Start: 3,
								End:   6,
							},
						},
					},
				},
			},
			nil,
			DefaultParser,
		},
	}

	for k, sample := range samples {
		tree, err := sample.parser.Parse(
			sample.syntax,
			[]byte(sample.text),
		)
		msg := spew.Sdump(k, sample.text)
		assert.Equal(t, sample.err, err, msg)
		assert.Equal(t, sample.tree, tree, msg)
	}
}
