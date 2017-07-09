package parse


// Terminal is a Rule which is literal in input.
type Terminal struct {
	name  string
	Value []byte
}

// Name indicates the name which was given to the rule
// on creation. Name could be not unique.
func (r *Terminal) Name() string {
	return r.name
}

func (r *Terminal) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Terminal) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Terminal) GetChilds() Treers {
	return nil
}

//

// GetParameters returns a KV rule parameters.
func (r *Terminal) GetParameters() RuleParameters {
	return RuleParameters{
		"name":  r.name,
		"value": string(r.Value),
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Terminal) IsFinite() bool {
	return true
}

//

// NewTerminal constructs a new *Terminal.
func NewTerminal(name string, v string) *Terminal {
	return &Terminal{name, []byte(v)}
}
