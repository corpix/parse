package parse

import (
	"fmt"
)

var _ Rule = new(Repetition)

func NewErrRepetitionTooMuchOccurrences(want, got int) error {
	return fmt.Errorf(
		"too much occurences of expression: want %d, got %d",
		want, got,
	)
}

func NewErrRepetitionNotEnoughOccurrences(want, got int) error {
	return fmt.Errorf(
		"not enough occurences of expression: want %d, got %d",
		want, got,
	)
}

var ErrRepetitionNothingMatched = fmt.Errorf("nothing matched")

// Repetition is a Rule which is repeating in the input
// one or more times.
type Repetition struct {
	name     string
	Rule     Rule
	Times    int
	Variadic bool
	Hooks    []RuleParseHook
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

// Parse consumes some bytes from input & emits a Tree
// using settings defined during creation of the concrete Rule type.
// May return an error if something goes wrong, should provide some
// location information to the user which points to position in input.
func (r *Repetition) Parse(ctx *Context, input []byte) (*Tree, error) {
	nextDepth := ctx.Depth + 1
	if nextDepth > ctx.Parser.MaxDepth {
		return nil, NewErrNestingTooDeep(ctx.Location, nextDepth)
	}

	var (
		occurrences = 0
		subInput    = input
		subTree     *Tree
		subChilds   = []*Tree{}
		pos         = ctx.Location.Position
		line, col   int
		loc         *Location
		err         error
	)
repeat:
	for {
		if len(subInput) == 0 {
			break
		}

		line, col = ctx.Parser.Locate(pos)
		loc = &Location{
			Path:     ctx.Location.Path,
			Position: pos,
			Line:     line,
			Column:   col,
		}
		subTree, err = r.Rule.Parse(
			&Context{
				Rule:     r,
				Parser:   ctx.Parser,
				Location: loc,
				Depth:    nextDepth,
			},
			subInput,
		)
		if err != nil {
			if err == ErrSkipRule {
				break
			}
			switch err.(type) {
			case *ErrUnexpectedToken, *ErrUnexpectedEOF:
				// XXX: We need to skip current rule
				// if it has no matches and rule is variadic,
				// should repeat 0 or more times.
				// In this case we have seen nothing and it is
				// ok to skip.
				if occurrences == 0 && r.Times == 0 && r.Variadic {
					return nil, ErrSkipRule
				}
				break repeat
			default:
				return nil, err
			}
		}
		occurrences++

		movePos := subTree.Region.End - subTree.Region.Start
		if !r.Variadic && occurrences > r.Times {
			return nil, NewErrUnexpectedToken(
				r,
				subTree.Location,
				ShowInput(input[pos:]),
				NewErrRepetitionTooMuchOccurrences(r.Times, occurrences),
			)
		}

		subInput = subInput[movePos:]
		subChilds = append(subChilds, subTree)
		pos += movePos
	}
	if err != nil && len(subChilds) == 0 { // nothing matched
		return nil, NewErrUnexpectedToken(
			r,
			loc,
			input,
			ErrRepetitionNothingMatched,
		)
	}
	if occurrences < r.Times {
		if err != nil {
			return nil, err
		}
		return nil, NewErrUnexpectedToken(
			r,
			loc,
			input,
			NewErrRepetitionNotEnoughOccurrences(r.Times, occurrences),
		)
	}

	region := TreeRegion(subChilds...)
	line, col = ctx.Parser.Locate(ctx.Location.Position)
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
		Childs: subChilds,
		Data:   input[:region.End-region.Start],
	}
	for _, hook := range r.Hooks {
		hook(ctx, tree)
	}
	return tree, nil
}

//

// NewRepetitionTimes constructs new *Repetition which repeats exactly `times`.
func NewRepetitionTimes(name string, times int, rule Rule, hooks ...RuleParseHook) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: false,
		Hooks:    hooks,
	}
}

// NewRepetitionTimesVariadic constructs new variadic *Repetition
// which repeats exactly `times` or more.
func NewRepetitionTimesVariadic(name string, times int, rule Rule, hooks ...RuleParseHook) *Repetition {
	return &Repetition{
		name:     name,
		Rule:     rule,
		Times:    times,
		Variadic: true,
		Hooks:    hooks,
	}
}

// NewRepetition constructs new *Repetition which releat one or more times.
func NewRepetition(name string, rule Rule, hooks ...RuleParseHook) *Repetition {
	return NewRepetitionTimesVariadic(name, 1, rule, hooks...)
}
