package parse

// Chain represents a chain of Rule's to match in the data.
type Chain struct {
	name  string
	Rules Rules
}

// Name indicates the name which was given to the rule
// on creation. Name could be not unique.
func (r *Chain) Name() string {
	return r.name
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (r *Chain) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Chain) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Chain) GetChilds() Treers {
	treers := make(Treers, len(r.Rules))
	for k, v := range r.Rules {
		treers[k] = v
	}
	return treers
}

//

// GetParameters returns a KV rule parameters.
func (r *Chain) GetParameters() RuleParameters {
	return RuleParameters{
		"name": r.name,
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Chain) IsFinite() bool {
	return false
}

//

// Add appends a Rule to the Chain.
func (r *Chain) Add(rule Rule) {
	r.Rules = append(r.Rules, rule)
}

//

// NewChain constructs new Chain.
// Valid Chain could be constructed with >=2 rules.
func NewChain(name string, r ...Rule) *Chain {
	return &Chain{
		name,
		r,
	}
}
