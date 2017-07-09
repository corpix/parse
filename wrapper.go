package parse


// Wrapper represents a wrapper type for some inner Rule.
// It could be used to wrap a Rule with custom name.
type Wrapper struct {
	name string
	Rule Rule
}

// Name indicates the name which was given to the rule
// on creation. Name could be not unique.
func (r *Wrapper) Name() string {
	return r.name
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (r *Wrapper) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Wrapper) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Wrapper) GetChilds() Treers {
	return Treers{r.Rule}
}

//

// GetParameters returns a KV rule parameters.
func (r *Wrapper) GetParameters() RuleParameters {
	return RuleParameters{
		"name": r.name,
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Wrapper) IsFinite() bool {
	return false
}

//

// NewWrapper constructs new Wrapper.
func NewWrapper(name string, r Rule) *Wrapper {
	return &Wrapper{name, r}
}
