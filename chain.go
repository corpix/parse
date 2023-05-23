package parse

var _ Rule = new(Chain)

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


// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Chain) Parse(ctx *Context, input []byte) (*Tree, error) {
	if len(r.Rules) == 0 {
		return nil, NewErrEmptyRule(r, ctx.Rule)
	}

	nextDepth := ctx.Location.Depth + 1
	if nextDepth > ctx.Parser.MaxDepth {
		return nil, NewErrNestingTooDeep(
			nextDepth,
			ctx.Location.Position,
		)
	}

	var (
		subInput  = input
		subTrees  = make([]*Tree, len(r.Rules))
		n         int
		pos       = ctx.Location.Position
		movPos    int
		line, col int
		err       error
	)
	for _, sr := range r.Rules {
		line, col = ctx.Parser.Locate(pos)
		subTrees[n], err = sr.Parse(
			&Context{
				Rule:   r,
				Parser: ctx.Parser,
				Location: &Location{
					Position: pos,
					Line:     line,
					Column:   col,
					Depth:    nextDepth,
				},
			},
			subInput,
		)
		if err != nil {
			if err == ErrSkipRule {
				continue
			}
			return nil, err
		}
		movPos = subTrees[n].Region.End - subTrees[n].Region.Start
		pos += movPos
		subInput = subInput[movPos:]
		n++
	}
	if err != nil {
		return nil, err
	}
	subTrees = subTrees[:n] // NOTE: because some Rule's could be skipped

	region := TreeRegion(subTrees...)
	line, col = ctx.Parser.Locate(ctx.Location.Position)
	return &Tree{
		Rule: r,
		Location: &Location{
			Position: ctx.Location.Position,
			Line:     line,
			Column:   col,
			Depth:    ctx.Location.Depth,
		},
		Region: region,
		Childs: subTrees,
		Data:   input[:region.End-region.Start],
	}, nil
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
