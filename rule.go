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
	"strings"
)

const (
	newLine   = "\n"
	space     = " "
	delimiter = ", "
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

// RuleFormatter is a configurable rule pretty
// printer.
type RuleFormatter struct {
	Pretty     bool
	Spaces     int
	MaxLineLen int
}

func (rs *RuleFormatter) indent(depth int, s string) string {
	if rs.Pretty {
		indent := strings.Repeat(space, rs.Spaces*depth)
		lines := strings.Split(s, newLine)
		for k, v := range lines {
			lines[k] = indent + v
		}
		return strings.Join(lines, newLine)
	}
	return s
}

func (rs *RuleFormatter) single(rule Rule, depth int, childs string) string {
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
			params += delimiter
			if rs.Pretty {
				params += newLine
			}
		}
		params += fmt.Sprintf(
			"%s: %v",
			v,
			parameters[v],
		)
	}

	if rs.Pretty && len(params) > rs.MaxLineLen {
		params = newLine + rs.indent(depth, params) + newLine
	} else {
		// XXX: roll back params if line too short.
		params = strings.Join(strings.Split(params, newLine), "")
	}
	if rs.Pretty && len(childs) > 0 {
		childs = newLine + rs.indent(depth, childs) + newLine
	}
	return fmt.Sprintf(
		"%T(%s)(%s)",
		rule,
		params,
		childs,
	)
}

func (rs *RuleFormatter) format(visited map[Rule]bool, depth int, rule Rule) string {
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
					child += delimiter
					if rs.Pretty {
						child += newLine
					}
				}
				child += rs.format(
					visited,
					depth+1,
					v.(Rule),
				)
			}
		}
	}

	return rs.indent(
		depth,
		rs.single(
			rule,
			depth,
			child,
		),
	)
}

// Format the Rule as string.
func (rs *RuleFormatter) Format(rule Rule) string {
	return rs.format(map[Rule]bool{}, 0, rule)
}

// RulesString folds a nested(maybe circular)
// rule into a string representation.
func RuleString(rule Rule) string {
	return (&RuleFormatter{}).
		Format(rule)
}

func RulePrettyString(rule Rule) string {
	return (&RuleFormatter{true, 2, 80}).
		Format(rule)
}
