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
)

var (
	DefaultParser = NewParser(128)
)

type Parser struct {
	maxDepth int
}

func (p *Parser) Parse(rule Rule, input []byte) (*Tree, error) {
	tree, err := p.parse(rule, input, 0, 0)
	if err != nil {
		return nil, err
	}

	if tree.End < len(input) {
		return nil, NewErrUnexpectedToken(
			input[tree.End:tree.End+1],
			p.humanizePosition(tree.End),
		)
	}

	return tree, nil
}

func (p *Parser) parse(rule Rule, input []byte, position int, depth int) (*Tree, error) {
	var (
		tree    *Tree
		subTree *Tree
		buf     []byte
		bufLen  int
		bounds  *Bounds
		err     error
	)

	if depth > p.maxDepth {
		return nil, NewErrNestingTooDeep(
			depth,
			p.humanizePosition(position),
		)
	}

	tree = &Tree{}

	switch v := rule.(type) {
	case Terminal:
		buf = input[:len(v)]
		bufLen = len(buf)

		if bufLen < len(v) {
			return nil, NewErrUnexpectedEOF(
				p.humanizePosition(position),
			)
		}

		if !bytes.EqualFold(buf, []byte(v)) {
			return nil, NewErrUnexpectedToken(
				buf[:1],
				p.humanizePosition(position),
			)
		}

		tree.Rule = v
		tree.Data = buf
		tree.Start = position
		tree.End = position + bufLen
	case Either:
	case Chain:
		for _, r := range v {
			subTree, err = p.parse(
				r,
				input[position:],
				position,
				depth+1,
			)
			if err != nil {
				return nil, err
			}
			position = subTree.End

			tree.Childs = append(
				tree.Childs,
				subTree,
			)
		}

		bounds = GetTreesBounds(tree.Childs)
		tree.Rule = v
		tree.Data = input[bounds.Starting:bounds.Closing]
		tree.Start = bounds.Starting
		tree.End = bounds.Closing
	case Repetition:
	default:
		return nil, NewErrUnsupportedRule(rule)
	}

	return tree, nil
}

// FIXME: We need position starting from 0
// to simplify code, but human readable errors and
// dumps should contain position starting from 1
// I am not sure this is ok to implement it in this manner
// but at this time let it be so.
func (p *Parser) humanizePosition(position int) int {
	return position + 1
}

func NewParser(maxDepth int) *Parser {
	return &Parser{maxDepth}
}
