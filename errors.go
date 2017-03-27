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
	e "errors"
	"fmt"
)

var (
	ErrStopIteration = e.New("Stop iteration")
	ErrSkipBranch    = e.New("Skip branch")
	ErrSkipRule      = e.New("Skip rule")
)

// ErrBoundIncomplete is an error which mean
// that a closing token was not
// found in the input which is making a requested
// logical «bound» to be incomplete.
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

// NewErrBoundIncomplete constructs new ErrBoundIncomplete.
func NewErrBoundIncomplete(position int, starting, closing []byte) error {
	return &ErrBoundIncomplete{starting, closing, position}
}

//

// ErrUnsupportedRule is an error which mean
// that parser support for specifier Rule is not implemented.
type ErrUnsupportedRule struct {
	Rule
}

func (e *ErrUnsupportedRule) Error() string {
	return fmt.Sprintf(
		"Unsupported rule type '%T'",
		e.Rule,
	)
}

// NewErrUnsupportedRule constructs new ErrUnsupportedRule.
func NewErrUnsupportedRule(rule Rule) error {
	return &ErrUnsupportedRule{rule}
}

//

// ErrUnexpectedEOF is an error which mean
// that EOF was meat while parser wanted more
// input.
type ErrUnexpectedEOF struct {
	Position int
	Rule
}

func (e *ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"Unexpected EOF at position '%d' while applying '%s'",
		e.Position,
		e.Rule,
	)
}

// NewErrUnexpectedEOF constructs new ErrUnexpectedEOF.
func NewErrUnexpectedEOF(position int, rule Rule) error {
	return &ErrUnexpectedEOF{position, rule}
}

//

// ErrUnexpectedToken is an error which mean
// that token read from current position in input
// is not expected by the current Rule.
type ErrUnexpectedToken struct {
	Token    []byte
	Position int
	Rule
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf(
		"Unexpected token '%s' at position '%d' while applying '%s'",
		e.Token,
		e.Position,
		e.Rule,
	)
}

// NewErrUnexpectedToken constructs new ErrUnexpectedToken.
func NewErrUnexpectedToken(token []byte, position int, rule Rule) error {
	return &ErrUnexpectedToken{token, position, rule}
}

//

// ErrNestingTooDeep is an error which mean
// the Rule nesting is too deep.
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

// NewErrNestingTooDeep constructs new ErrNestingTooDeep.
func NewErrNestingTooDeep(nesting int, position int) error {
	return &ErrNestingTooDeep{nesting, position}
}

//

// ErrEmptyRule is an error which mean
// a Rule with empty content was passed to the parser.
type ErrEmptyRule struct {
	Rule
	Inside Rule
}

func (e *ErrEmptyRule) Error() string {
	return fmt.Sprintf(
		"Empty rule of type '%T' = '%s' inside '%s' rule",
		e.Rule,
		e.Rule,
		e.Inside,
	)
}

// NewErrEmptyRule constructs new ErrEmptyRule.
func NewErrEmptyRule(rule Rule, inside Rule) error {
	return &ErrEmptyRule{rule, inside}
}
