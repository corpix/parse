package parse

import (
	"fmt"
	"reflect"
	"sort"
)

// RuleParameters is a key-value mapping of
// the Rule custom settings.
type RuleParameters map[string]interface{}

// String encodes a RuleParameters as string.
func (p RuleParameters) String() string {
	var (
		keys   = []string{}
		params = ""
	)
	for k := range p {
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
			v, indirectValue(reflect.ValueOf(p[v])).Interface(),
		)
	}

	return params
}

//

type RuleParseHook func(ctx *Context, t *Tree) error

// Rule represents a general Rule interface.
type Rule interface {
	Treer

	// Parameters returns a KV rule parameters.
	GetParameters() RuleParameters

	// IsFinite returns true if this rule is
	// not a wrapper for other rules.
	IsFinite() bool

	// Parse consumes some bytes from input & emits a Tree
	// using settings defined during creation of the concrete Rule type.
	// May return an error if something goes wrong, should provide some
	// location information to the user which points to position in input.
	Parse(ctx *Context, input []byte, hooks ...RuleParseHook) (*Tree, error)
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
