package parse

import (
	"fmt"
)

var _ Rule = new(Either)

// Either represents a list of Rule's to match in the data.
// One of the rules in a list must match.
type Either struct {
	name  string
	Rules Rules
}

// Name indicates the name which was given to the rule
// on creation. Name could be not unique.
func (r *Either) Name() string {
	return r.name
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (r *Either) Show(childs string) string {
	return RuleShow(
		r,
		r.GetParameters().String(),
		childs,
	)
}

// String returns rule as a string,
// resolving recursion with `<circular>` placeholder.
func (r *Either) String() string {
	return TreerString(r)
}

// GetChilds returns a slice of Rule which is
// children for current Rule.
func (r *Either) GetChilds() Treers {
	treers := make(Treers, len(r.Rules))
	for k, v := range r.Rules {
		treers[k] = v
	}
	return treers
}

//

// GetParameters returns a KV rule parameters.
func (r *Either) GetParameters() RuleParameters {
	return RuleParameters{
		"name": r.name,
	}
}

// IsFinite returns true if this rule is
// not a wrapper for other rules.
func (r *Either) IsFinite() bool {
	return false
}

func (r *Either) Parse(ctx *Context, input []byte) (*Tree, error) {
	if len(r.Rules) == 0 {
		return nil, NewErrEmptyRule(r, ctx.Rule)
	}
	if len(input) == 0 {
		return nil, NewErrUnexpectedEOF(ctx.Location.Position, r)
	}

	nextDepth := ctx.Location.Depth + 1
	if nextDepth > ctx.Parser.MaxDepth {
		return nil, NewErrNestingTooDeep(
			nextDepth,
			ctx.Location.Position,
		)
	}

	var (
		subTree *Tree
		err     error
	)
	for _, sr := range r.Rules {
		subTree, err = sr.Parse(
			&Context{
				Rule:   sr,
				Parser: ctx.Parser,
				Location: &Location{
					Position: ctx.Location.Position,
					Line:     ctx.Location.Line,   // FIXME
					Column:   ctx.Location.Column, // FIXME
					Depth:    nextDepth,
				},
			},
			input,
		)
		if err != nil {
			if err == ErrSkipRule {
				continue
			}
			switch err.(type) {
			case *ErrUnexpectedToken, *ErrUnexpectedEOF:
				continue
			default:
				return nil, err
			}
		}
		break
	}
	if subTree == nil {
		return nil, NewErrUnexpectedToken(
			input,
			ctx.Location.Position,
			r, // FIXME: pass sub-error to reason more precisely?
		)
	}
	if err != nil {
		return nil, err
	}

	region := *subTree.Region
	return &Tree{
		Rule: r,
		Location: &Location{
			Position: subTree.Location.Position,
			Line:     ctx.Location.Line,   // FIXME
			Column:   ctx.Location.Column, // FIXME
			Depth:    ctx.Location.Depth,
		},
		Region: &region,
		Childs: []*Tree{subTree},
		Data:   input[:subTree.Region.End-subTree.Region.Start],
	}, nil
}

//

// Add appends a Rule into Either list.
func (r *Either) Add(rule Rule) {
	r.Rules = append(r.Rules, rule)
}

//

// NewEither constructs *Either Rule.
// Valid Either could be constructed with >=2 rules.
func NewEither(name string, r ...Rule) *Either {
	return &Either{
		name,
		r,
	}
}

// NewASCIIRange constructs *Either(Terminal, ...) Rule using specified ASCII range.
func NewASCIIRange(name string, from byte, to byte) *Either {
	if from > to {
		panic(fmt.Errorf(
			"invalid range, `from` (%d) should be less than `to` (%d)",
			from, to,
		))
	}

	amount := to - from
	terms := make([]Rule, amount+1)
	for chr := to; chr >= from; chr-- {
		terms[amount] = NewTerminal(string(chr), string(chr))
		amount--
	}

	return NewEither(name, terms...)
}
