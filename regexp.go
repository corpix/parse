package parse

import (
	"regexp"
)

// Regexp is a Rule which should match Go regexp on input.
type Regexp struct {
	name   string
	Regexp *regexp.Regexp
	Expr   string
}

// Name indicates the name which was given to the rule
// on creation. Name could be not unique.
func (r *Regexp) Name() string {
	return r.name
}

func (r *Regexp) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Regexp) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Regexp) GetChilds() Treers {
	return nil
}

//

// GetParameters returns a KV rule parameters.
func (r *Regexp) GetParameters() RuleParameters {
	return RuleParameters{
		"name": r.name,
		"expr": string(r.Expr),
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Regexp) IsFinite() bool {
	return true
}

//

// NewRegexp constructs a new *Regexp.
func NewRegexp(name string, expr string) *Regexp {
	return &Regexp{
		name:   name,
		Regexp: regexp.MustCompile(expr),
		Expr:   expr,
	}
}
