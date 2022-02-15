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
			v,
			indirectValue(
				reflect.ValueOf(p[v]),
			).Interface(),
		)
	}

	return params
}

// Rule represents a general Rule interface.
type Rule interface {
	Treer

	// Parameters returns a KV rule parameters.
	GetParameters() RuleParameters

	// IsFinite returns true if this rule is
	// not a wrapper for other rules.
	IsFinite() bool
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
