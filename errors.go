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
	"fmt"
)

type ErrBoundIncomplete struct {
	Starting []byte
	Closing  []byte
	Position int
}

func (e *ErrBoundIncomplete) Error() string {
	return fmt.Sprintf(
		"Bound start token '%s' found but close token '%s' is not, bound incomplete at position %d",
		string(e.Starting),
		string(e.Closing),
		e.Position,
	)
}

func NewErrBoundIncomplete(position int, starting, closing []byte) error {
	return &ErrBoundIncomplete{starting, closing, position}
}

//

type ErrUnsupportedRule struct {
	Rule
}

func (e *ErrUnsupportedRule) Error() string {
	return fmt.Sprintf(
		"Unsupported rule type '%T'",
		e.Rule,
	)
}

func NewErrUnsupportedRule(rule Rule) error {
	return &ErrUnsupportedRule{rule}
}

//

type ErrUnexpectedEOF struct {
	Position int
}

func (e *ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"Unexpected EOF at position '%d'",
		e.Position,
	)
}

func NewErrUnexpectedEOF(position int) error {
	return &ErrUnexpectedEOF{position}
}

//

type ErrUnexpectedToken struct {
	Token    []byte
	Position int
	Rule
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf(
		"Unexpected token '%s' at position '%d' while applying '%#v'",
		e.Token,
		e.Position,
		e.Rule,
	)
}

func NewErrUnexpectedToken(token []byte, position int, rule Rule) error {
	return &ErrUnexpectedToken{token, position, rule}
}

//

type ErrNestingTooDeep struct {
	Nesting  int
	Position int
}

func (e *ErrNestingTooDeep) Error() string {
	return fmt.Sprintf(
		"Nesting too deep, counted to '%d' levels at position %d",
		e.Nesting,
		e.Position,
	)
}

func NewErrNestingTooDeep(nesting int, position int) error {
	return &ErrNestingTooDeep{nesting, position}
}

//

type ErrEmptyRule struct {
	Rule
}

func (e *ErrEmptyRule) Error() string {
	return fmt.Sprintf(
		"Empty rule '%T' is '%#v' is not allowed",
		e.Rule,
		e.Rule,
	)
}

func NewErrEmptyRule(rule Rule) error {
	return &ErrEmptyRule{rule}
}
