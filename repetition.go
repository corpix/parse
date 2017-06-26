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

// Repetition is a Rule which is repeating in the input
// one or more times.
type Repetition struct {
	name     string
	Rule     Rule
	Times    int
	Variadic bool
}

// Name indicates the Name which was given to the rule
// on creation. Name could be not unique.
func (r *Repetition) Name() string {
	return r.name
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (r *Repetition) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().Show(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Repetition) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Repetition) GetChilds() Treers {
	return Treers{r.Rule}
}

//

// GetParameters returns a KV rule parameters.
func (r *Repetition) GetParameters() RuleParameters {
	return RuleParameters{
		"name":     r.name,
		"times":    r.Times,
		"variadic": r.Variadic,
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Repetition) IsFinite() bool {
	return false
}

//

// NewRepetitionTimes constructs new *Repetition which repeats exactly `times`.
func NewRepetitionTimes(name string, times int, rule Rule) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: false,
	}
}

// NewRepetitionTimesVariadic constructs new variadic *Repetition
// which repeats exactly `times` or more.
func NewRepetitionTimesVariadic(name string, times int, rule Rule) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: true,
	}
}

// NewRepetition constructs new *Repetition which releat one or more times.
func NewRepetition(name string, rule Rule) *Repetition {
	return NewRepetitionTimesVariadic(name, 1, rule)
}
