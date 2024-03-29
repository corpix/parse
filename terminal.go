package parse

import (
	"bytes"
	"unicode/utf8"
)

var _ Rule = new(Terminal)

// Terminal is a Rule which is literal in input.
type Terminal struct {
	name  string
	Value []byte
	Hooks []RuleParseHook
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

// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Terminal) Parse(ctx *Context, input []byte) (*Tree, error) {
	length := utf8.RuneCount(r.Value)
	if length == 0 {
		return nil, NewErrEmptyRule(r, ctx.Rule)
	}

	if utf8.RuneCount(input) < length {
		return nil, NewErrUnexpectedEOF(r, ctx.Location)
	}

	buf := input[:length]
	if !bytes.Equal(buf, r.Value) {
		return nil, NewErrUnexpectedToken(r, ctx.Location, ShowInput(input))
	}

	line, col := ctx.Parser.Locate(ctx.Location.Position)
	tree := &Tree{
		Rule: r,
		Location: &Location{
			Path:     ctx.Location.Path,
			Position: ctx.Location.Position,
			Line:     line,
			Column:   col,
		},
		Region: &Region{
			Start: ctx.Location.Position,
			End:   ctx.Location.Position + length,
		},
		Depth: ctx.Depth,
		Data:  buf,
	}
	for _, hook := range r.Hooks {
		hook(ctx, tree)
	}
	return tree, nil
}

//

// NewTerminal constructs a new *Terminal.
func NewTerminal(name string, v string, hooks ...RuleParseHook) *Terminal {
	return &Terminal{
		name:  name,
		Value: []byte(v),
		Hooks: hooks,
	}
}
