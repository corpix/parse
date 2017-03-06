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
	starting []byte
	closing  []byte
	position int
}

func (e *ErrBoundIncomplete) Error() string {
	return fmt.Sprintf(
		"Bound start token '%s' found but close token '%s' is not, bound incomplete at position %d",
		string(e.starting),
		string(e.closing),
		e.position,
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
	position int
}

func (e *ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"Unexpected EOF at position '%d'",
		e.position,
	)
}

func NewErrUnexpectedEOF(position int) error {
	return &ErrUnexpectedEOF{position}
}

//

type ErrUnexpectedToken struct {
	token    []byte
	position int
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf(
		"Unexpected token '%s' at position '%d'",
		e.token,
		e.position,
	)
}

func NewErrUnexpectedToken(token []byte, position int) error {
	return &ErrUnexpectedToken{token, position}
}

//

type ErrNestringTooDeep struct {
	nesting  int
	position int
}

func (e *ErrNestringTooDeep) Error() string {
	return fmt.Sprintf(
		"Nesting too deep, counted to '%d' levels at position %d",
		e.nesting,
		e.position,
	)
}

func NewErrNestingTooDeep(nesting int, position int) error {
	return &ErrNestringTooDeep{nesting, position}
}
