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
	Hooks []RuleParseHook
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

// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Either) Parse(ctx *Context, input []byte) (*Tree, error) {
	if len(r.Rules) == 0 {
		return nil, NewErrEmptyRule(r, ctx.Rule)
	}
	if len(input) == 0 {
		return nil, NewErrUnexpectedEOF(r, ctx.Location)
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
		line, col = ctx.Parser.Locate(ctx.Location.Position)
		err       error
	)
	for _, sr := range r.Rules {
		subTree, err = sr.Parse(
			&Context{
				Rule:   sr,
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
			r,
			ctx.Location,
			input,
			err,
		)
	}
	if err != nil {
		return nil, err
	}

	region := *subTree.Region
	tree := &Tree{
		Rule: r,
		Location: &Location{
			Path:     ctx.Location.Path,
			Position: subTree.Location.Position,
			Line:     line,
			Column:   col,
		},
		Region: &region,
		Depth:  ctx.Depth,
		Childs: []*Tree{subTree},
		Data:   input[:subTree.Region.End-subTree.Region.Start],
	}
	for _, hook := range r.Hooks {
		hook(ctx, tree)
	}
	return tree, nil
}

//

// Add appends a Rule into Either list.
func (r *Either) Add(rule ...Rule) {
	r.Rules = append(r.Rules, rule...)
}

//

// NewEither constructs *Either Rule.
// Valid Either could be constructed with >=2 rules.
func NewEither(name string, rulesOrHooks ...interface{}) *Either {
	rules := []Rule{}
	hooks := []RuleParseHook{}
	for _, ruleOrHook := range rulesOrHooks {
		switch v := ruleOrHook.(type) {
		case Rule:
			rules = append(rules, v)
		case RuleParseHook:
			hooks = append(hooks, v)
		default:
			panic(fmt.Sprintf("unsupported type %T", ruleOrHook))
		}
	}
	return &Either{
		name:  name,
		Rules: rules,
		Hooks: hooks,
	}
}

// NewASCIIRange constructs *Either(Terminal, ...) Rule using specified ASCII range.
func NewASCIIRange(name string, from byte, to byte, hooks ...RuleParseHook) *Either {
	if from > to {
		panic(fmt.Errorf(
			"invalid range, `from` (%d) should be less than `to` (%d)",
			from, to,
		))
	}

	rulesOrHooks := make([]interface{}, int(to-from)+1+len(hooks))
	n := 0
	for chr := from; chr <= to; chr++ {
		rulesOrHooks[n] = NewTerminal(string(chr), string(chr))
		n++
	}
	for _, hook := range hooks {
		rulesOrHooks[n] = hook
		n++
	}

	return NewEither(name, rulesOrHooks...)
}
