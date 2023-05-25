package parse

import (
	e "errors"
	"fmt"
)

var (
	ErrStopIteration = e.New("Stop iteration")
	ErrSkipBranch    = e.New("Skip branch")
	ErrSkipRule      = e.New("Skip rule")
)

// ErrBoundIncomplete is an error which mean
// that a closing token was not
// found in the input which is making a requested
// logical «bound» to be incomplete.
type ErrBoundIncomplete struct {
	Starting []byte
	Closing  []byte
	Location *Location
}

func (e *ErrBoundIncomplete) Error() string {
	return fmt.Sprintf(
		"Bound start token '%s' found but close token '%s' is not, bound incomplete at %q",
		string(e.Starting),
		string(e.Closing),
		e.Location,
	)
}

// NewErrBoundIncomplete constructs new ErrBoundIncomplete.
func NewErrBoundIncomplete(starting, closing []byte, l *Location) error {
	return &ErrBoundIncomplete{starting, closing, l}
}

//

// ErrUnsupportedRule is an error which mean
// that parser support for specifier Rule is not implemented.
type ErrUnsupportedRule struct {
	Rule Rule
}

func (e *ErrUnsupportedRule) Error() string {
	return fmt.Sprintf(
		"Unsupported rule type '%T'",
		e.Rule,
	)
}

// NewErrUnsupportedRule constructs new ErrUnsupportedRule.
func NewErrUnsupportedRule(rule Rule) error {
	return &ErrUnsupportedRule{rule}
}

//

// ErrUnexpectedEOF is an error which mean
// that EOF was meat while parser wanted more
// input.
type ErrUnexpectedEOF struct {
	Rule     Rule
	Location *Location
}

func (e *ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"Unexpected EOF at %q while applying '%s' rule",
		e.Location,
		e.Rule.Name(),
	)
}

// NewErrUnexpectedEOF constructs new ErrUnexpectedEOF.
func NewErrUnexpectedEOF(r Rule, l *Location) error {
	return &ErrUnexpectedEOF{r, l}
}

//

// ErrUnexpectedToken is an error which mean
// that token read from current position in input
// is not expected by the current Rule.
type ErrUnexpectedToken struct {
	Rule     Rule
	Location *Location
	Token    []byte
	Inner    []error
}

func (e *ErrUnexpectedToken) Error() string {
	innerErrs := ""
	for _, innerErr := range e.Inner {
		innerErrs = ": " + innerErr.Error()
	}

	return fmt.Sprintf(
		"Unexpected token '%s' at %q while applying '%s' rule%s",
		e.Token,
		e.Location,
		e.Rule.Name(),
		innerErrs,
	)
}

// NewErrUnexpectedToken constructs new ErrUnexpectedToken.
func NewErrUnexpectedToken(r Rule, l *Location, token []byte, inner ...error) error {
	return &ErrUnexpectedToken{
		Rule:     r,
		Location: l,
		Token:    token,
		Inner:    inner,
	}
}

//

// ErrNestingTooDeep is an error which mean
// the Rule nesting is too deep.
type ErrNestingTooDeep struct {
	Location *Location
	Depth    int
}

func (e *ErrNestingTooDeep) Error() string {
	return fmt.Sprintf(
		"Nesting too deep, counted '%d' levels at %q",
		e.Depth,
		e.Location,
	)
}

// NewErrNestingTooDeep constructs new ErrNestingTooDeep.
func NewErrNestingTooDeep(l *Location, depth int) error {
	return &ErrNestingTooDeep{l, depth}
}

//

// ErrEmptyRule is an error which mean
// a Rule with empty content was passed to the parser.
type ErrEmptyRule struct {
	Rule   Rule
	Inside Rule
}

func (e *ErrEmptyRule) Error() string {
	var inside string
	if e.Inside == nil {
		inside = "<root>"
	} else {
		inside = e.Inside.Name()
	}
	return fmt.Sprintf(
		"Empty rule of type '%T' = '%s' inside '%s' rule",
		e.Rule,
		e.Rule,
		inside,
	)
}

// NewErrEmptyRule constructs new ErrEmptyRule.
func NewErrEmptyRule(rule Rule, inside Rule) error {
	return &ErrEmptyRule{rule, inside}
}
