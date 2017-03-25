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

import ()

// Repetition is a Rule which is repeating in the input
// one or more times.
type Repetition struct {
	id       string
	Rule     Rule
	Times    int
	Variadic bool
}

// ID indicates the ID which was given to the rule
// on creation. ID could be not unique.
func (r *Repetition) ID() string {
	return r.id
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Repetition) GetChilds() Treers {
	return Treers{r.Rule}
}

// GetParameters returns a KV rule parameters.
func (r *Repetition) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"ID":       r.id,
		"Times":    r.Times,
		"Variadic": r.Variadic,
	}
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Repetition) String() string {
	return RulePrettyString(r)
}

// NewRepetition constructs new *Repetition which repeats exactly `times`.
func NewRepetitionTimes(id string, times int, rule Rule) *Repetition {
	return &Repetition{
		id:       id,
		Rule:     rule,
		Times:    times,
		Variadic: false,
	}
}

// NewRepetitionTimesVariadic constructs new variadic *Repetition
// which repeats exactly `times` or more.
func NewRepetitionTimesVariadic(id string, times int, rule Rule) *Repetition {
	return &Repetition{
		id:       id,
		Rule:     rule,
		Times:    times,
		Variadic: true,
	}
}

// NewRepetition constructs new *Repetition which releat one or more times.
func NewRepetition(id string, rule Rule) *Repetition {
	return NewRepetitionTimesVariadic(id, 1, rule)
}
