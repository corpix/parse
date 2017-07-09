package parse


// Repetition is a Rule which is repeating in the input
// one or more times.
type Repetition struct {
	name     string
	Rule     Rule
	Times    int
	Variadic bool
}

// Name indicates the Name which was given to the rule
// on creation. Name could be not unique.
func (r *Repetition) Name() string {
	return r.name
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (r *Repetition) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Repetition) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Repetition) GetChilds() Treers {
	return Treers{r.Rule}
}

//

// GetParameters returns a KV rule parameters.
func (r *Repetition) GetParameters() RuleParameters {
	return RuleParameters{
		"name":     r.name,
		"times":    r.Times,
		"variadic": r.Variadic,
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Repetition) IsFinite() bool {
	return false
}

//

// NewRepetitionTimes constructs new *Repetition which repeats exactly `times`.
func NewRepetitionTimes(name string, times int, rule Rule) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: false,
	}
}

// NewRepetitionTimesVariadic constructs new variadic *Repetition
// which repeats exactly `times` or more.
func NewRepetitionTimesVariadic(name string, times int, rule Rule) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: true,
	}
}

// NewRepetition constructs new *Repetition which releat one or more times.
func NewRepetition(name string, rule Rule) *Repetition {
	return NewRepetitionTimesVariadic(name, 1, rule)
}
