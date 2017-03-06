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
	"bytes"
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
	}{
		{
			"",
			Terminal(""),
			&Tree{
				Rule:  Terminal(""),
				Data:  []byte{},
				Start: 1,
				End:   1,
			},
			nil,
		},
		{
			"foo",
			Terminal("foo"),
			&Tree{
				Rule:  Terminal("foo"),
				Data:  []byte("foo"),
				Start: 1,
				End:   4,
			},
			nil,
		},
	}

	for k, sample := range samples {
		reader := bytes.NewBufferString(sample.text)
		tree, err := Parse(sample.syntax, reader)

		msg := spew.Sdump(k, sample)
		assert.Equal(t, sample.err, err, msg)
		assert.Equal(t, sample.tree, tree, msg)
	}
}
