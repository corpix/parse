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
	"bytes"
)

var (
	// DefaultParser is a Parser with default settings.
	DefaultParser = NewParser(128)
)

// Parser represents a parser which use Rule's
// to parse the input.
type Parser struct {
	maxDepth int
}

// Parse parses input with Rule's.
func (p *Parser) Parse(rule Rule, input []byte) (*Tree, error) {
	tree, err := p.parse(rule, input, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	if tree.End < len(input) {
		return nil, NewErrUnexpectedToken(
			ShowInput(input[tree.End:]),
			p.humanizePosition(tree.End),
			rule,
		)
	}

	return tree, nil
}

func (p *Parser) parse(rule Rule, input []byte, parent Rule, position int, depth int) (*Tree, error) {
	var (
		tree    *Tree
		subTree *Tree
		buf     []byte
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

	if rule == nil {
		return nil, NewErrEmptyRule(rule, parent)
	}

	switch v := rule.(type) {
	case *Terminal:
		length := len(v.Value)
		if length == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		if len(input) < length {
			return nil, NewErrUnexpectedEOF(
				p.humanizePosition(position),
			)
		}
		buf = input[:length]

		if !bytes.EqualFold(buf, v.Value) {
			return nil, NewErrUnexpectedToken(
				ShowInput(buf),
				p.humanizePosition(position),
				v,
			)
		}

		tree.Start = position
		tree.End = position + length
		tree.Rule = v
		tree.Data = buf
	case *Either:
		if len(v.Rules) == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		for _, r := range v.Rules {
			subTree, err = p.parse(
				r,
				input,
				v,
				position,
				depth+1,
			)
			if err != nil {
				switch err.(type) {
				case *ErrUnexpectedToken, *ErrUnexpectedEOF:
					continue
				default:
					return nil, err
				}
			}

			break
		}
		if err != nil {
			return nil, err
		}

		tree.Childs = make([]*Tree, 1)
		tree.Childs[0] = subTree
		tree.Start = subTree.Start
		tree.End = subTree.End
		tree.Rule = v
		tree.Data = input[:subTree.End-subTree.Start]
	case *Chain:
		if len(v.Rules) == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		buf = input
		for _, r := range v.Rules {
			subTree, err = p.parse(
				r,
				buf,
				v,
				position,
				depth+1,
			)
			if err != nil {
				return nil, err
			}
			position = subTree.End
			buf = buf[subTree.End-subTree.Start:]

			tree.Childs = append(
				tree.Childs,
				subTree,
			)
		}

		bounds = GetTreesBounds(tree.Childs)
		tree.Start = bounds.Starting
		tree.End = bounds.Closing
		tree.Rule = v
		tree.Data = input[:bounds.Closing-bounds.Starting]
	case *Repetition:
		seen := 0
		buf = input

	repetitionLoop:
		for {
			subTree, err = p.parse(
				v.Rule,
				buf,
				v,
				position,
				depth+1,
			)
			if err != nil {
				switch err.(type) {
				case *ErrUnexpectedToken, *ErrUnexpectedEOF:
					break repetitionLoop
				default:
					return nil, err
				}
			}
			seen++

			movePos := subTree.End - subTree.Start
			if !v.Variadic && seen > v.Times {
				return nil, NewErrUnexpectedToken(
					ShowInput(input[position:]),
					position+movePos,
					v,
				)
			}

			position += movePos
			buf = buf[movePos:]

			tree.Childs = append(
				tree.Childs,
				subTree,
			)
		}
		if seen < v.Times {
			if err != nil {
				return nil, err
			}
			return nil, NewErrUnexpectedToken(
				input,
				position,
				v,
			)
		}

		bounds = GetTreesBounds(tree.Childs)
		tree.Start = bounds.Starting
		tree.End = bounds.Closing
		tree.Rule = v
		tree.Data = input[:bounds.Closing-bounds.Starting]
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

// Parse is a shortcut to call the DefaultParser.Parse().
func Parse(rule Rule, input []byte) (*Tree, error) {
	return DefaultParser.Parse(rule, input)
}

// NewParser constructs new *Parser.
func NewParser(maxDepth int) *Parser {
	return &Parser{maxDepth}
}
