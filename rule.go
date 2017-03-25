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
	"sort"
)

// Rule represents a general Rule interface.
type Rule interface {
	Treer
	// Parameters returns a KV rule parameters.
	GetParameters() map[string]interface{}

	// String returns rule as a string,
	// resolving recursion with `<circular>` placeholder.
	String() string
}

const (
	circularLabel = "<circular>"
)

func ruleShow(rule Rule, childs string) string {
	var (
		parameters = rule.GetParameters()
		paramKeys  = []string{}
		params     = ""
	)
	for k, _ := range parameters {
		paramKeys = append(
			paramKeys,
			k,
		)
	}
	sort.Strings(paramKeys)
	for k, v := range paramKeys {
		if k > 0 {
			params += ", "
		}
		params += fmt.Sprintf(
			"%s: %v",
			v,
			parameters[v],
		)
	}

	return fmt.Sprintf(
		"%T(%s)(%s)",
		rule,
		params,
		childs,
	)
}

func ruleString(visited map[Rule]bool, rule Rule) string {
	var (
		child string
	)

	if _, ok := visited[rule]; ok {
		child = circularLabel
	} else {
		visited[rule] = true

		childs := rule.GetChilds()
		if len(childs) > 0 {
			for k, v := range childs {
				if k > 0 {
					child += ", "
				}
				child += ruleString(
					visited,
					v.(Rule),
				)
			}
		}
	}

	return ruleShow(
		rule,
		child,
	)
}

// RulesString folds a nested(maybe circular)
// rule into a string representation.
func RuleString(rule Rule) string {
	return ruleString(map[Rule]bool{}, rule)
}
