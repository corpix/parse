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
	"fmt"
	"reflect"
	"sort"
)

// RuleParameters is a key-value mapping of
// the Rule custom settings.
type RuleParameters map[string]interface{}

// Rule represents a general Rule interface.
type Rule interface {
	Treer

	// Parameters returns a KV rule parameters.
	GetParameters() RuleParameters

	// IsFinite returns true if this rule is
	// not a wrapper for other rules.
	IsFinite() bool
}

// FIXME: Probably this function should be a RuleParameters.Show()?
// In this case we could make RuleShow signature more simple, but should we?
// RuleParametersShow returns a Rule.GetParameters() encoded as string.
func RuleParametersShow(parameters RuleParameters) string {
	var (
		keys   = []string{}
		params = ""
	)
	for k, _ := range parameters {
		keys = append(
			keys,
			k,
		)
	}
	sort.Strings(keys)
	for k, v := range keys {
		if k > 0 {
			params += treerDelimiter
		}
		params += fmt.Sprintf(
			"%s: %v",
			v,
			indirectValue(
				reflect.ValueOf(parameters[v]),
			).Interface(),
		)
	}

	return params
}

// RuleShow returns a Rule encoded as a string.
// It requires some parts to be prepared(encoded into a string).
func RuleShow(rule Rule, parameters string, childs string) string {
	return fmt.Sprintf(
		"%T(%s)(%s)",
		rule,
		parameters,
		childs,
	)
}
