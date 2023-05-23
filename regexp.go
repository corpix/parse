package parse

import (
	"regexp"
)

var _ Rule = new(Regexp)

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

// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Regexp) Parse(ctx *Context, input []byte) (*Tree, error) {
	buf := r.Regexp.Find(input)
	if buf == nil {
		return nil, NewErrUnexpectedToken(
			r,
			ctx.Location,
			ShowInput(input),
		)
	}
	regexpRegion := r.Regexp.FindIndex(input)

	line, col := ctx.Parser.Locate(ctx.Location.Position)
	return &Tree{
		Rule: r,
		Location: &Location{
			Position: ctx.Location.Position,
			Line:     line,
			Column:   col,
			Depth:    ctx.Location.Depth,
		},
		Region: &Region{
			Start: ctx.Location.Position + regexpRegion[0],
			End:   ctx.Location.Position + regexpRegion[1],
		},
		Data: buf,
	}, nil
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
