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

// Either represents a list of Rule's to match in the data.
// One of the rules in a list must match.
type Either struct {
	id    string
	Rules Rules
}

// ID indicates the ID which was given to the rule
// on creation. ID could be not unique.
func (r *Either) ID() string {
	return r.id
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Either) GetChilds() Treers {
	treers := make(Treers, len(r.Rules))
	for k, v := range r.Rules {
		treers[k] = v
	}
	return treers
}

// GetParameters returns a KV rule parameters.
func (r *Either) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"ID": r.id,
	}
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Either) String() string {
	return RulePrettyString(r)
}

// Add appends a Rule into Either list.
func (r *Either) Add(rule Rule) {
	r.Rules = append(r.Rules, rule)
}

// NewEither constructs *Either Rule.
func NewEither(id string, rules ...Rule) *Either {
	return &Either{
		id,
		Rules(rules),
	}
}
