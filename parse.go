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
	"io"
)

func Parse(rule Rule, input io.Reader) (*Tree, error) {
	tree, err := parse(rule, input, 1, 0)
	if err != nil {
		return nil, err
	}

	// FIXME: Could we do better?
	buf := make([]byte, 1)
	read, err := input.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if read > 0 {
		return nil, NewErrUnexpectedToken(buf, tree.End)
	}

	return tree, nil
}

func parse(rule Rule, input io.Reader, position int, depth int) (*Tree, error) {
	var (
		tree    *Tree
		subTree *Tree
		err     error
		buf     []byte
		read    int
	)

	tree = &Tree{}

	switch v := rule.(type) {
	case Terminal:
		buf = make([]byte, len(v))
		read, err = input.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if read < len(buf) {
			return nil, NewErrUnexpectedEOF(position)
		}

		if !bytes.EqualFold(buf, []byte(v)) {
			return nil, NewErrUnexpectedToken(buf, position)
		}

		tree.Rule = v
		tree.Data = buf
		tree.Start = position
		tree.End = position + read

		position = tree.End
	case Either:
	case Chain:
		for _, r := range v {
			subTree, err = parse(
				r,
				input,
				position,
				depth+1,
			)
			if err != nil {
				return nil, err
			}
			tree.Child = append(
				tree.Child,
				subTree,
			)
		}
	case Repetition:
	default:
		return nil, NewErrUnsupportedRule(rule)
	}

	return tree, nil
}
