package parse

var _ Rule = new(Wrapper)

// Wrapper represents a wrapper type for some inner Rule.
// It could be used to wrap a Rule with custom name.
type Wrapper struct {
	name  string
	Rule  Rule
	Hooks []RuleParseHook
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

// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Wrapper) Parse(ctx *Context, input []byte) (*Tree, error) {
	if r.Rule == nil {
		return nil, NewErrEmptyRule(r, r.Rule)
	}

	nextDepth := ctx.Depth + 1
	if nextDepth > ctx.Parser.MaxDepth {
		return nil, NewErrNestingTooDeep(
			ctx.Location,
			nextDepth,
		)
	}

	var (
		subTree   *Tree
		err       error
		line, col = ctx.Parser.Locate(ctx.Location.Position)
	)

	subTree, err = r.Rule.Parse(
		&Context{
			Rule:   r,
			Parser: ctx.Parser,
			Location: &Location{
				Path:     ctx.Location.Path,
				Position: ctx.Location.Position,
				Line:     line,
				Column:   col,
			},
			Depth: nextDepth,
		},
		input,
	)
	if err != nil {

		return nil, err
	}

	region := TreeRegion(subTree)
	tree := &Tree{
		Rule: r,
		Location: &Location{
			Path:     ctx.Location.Path,
			Position: ctx.Location.Position,
			Line:     line,
			Column:   col,
		},
		Region: region,
		Depth:  ctx.Depth,
		Childs: []*Tree{subTree},
		Data:   input[:region.End-region.Start],
	}
	for _, hook := range r.Hooks {
		hook(ctx, tree)
	}
	return tree, nil
}

//

// NewWrapper constructs new Wrapper.
func NewWrapper(name string, r Rule, hooks ...RuleParseHook) *Wrapper {
	return &Wrapper{
		name:  name,
		Rule:  r,
		Hooks: hooks,
	}
}
