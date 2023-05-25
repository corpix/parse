package parse

import (
	"fmt"
	"unicode/utf8"
)

var (
	DefaultParserMaxDepth  = 128
	DefaultParserLineBreak = NewEither(
		"line-break",
		NewTerminal("lf", "\n"),
		NewTerminal("crlf", "\r\n"),
	)
	DefaultParserPath    = "?"
	DefaultParserOptions = []ParserOption{
		ParserOptionMaxDepth(DefaultParserMaxDepth),
		ParserOptionLineBreak(DefaultParserLineBreak),
		ParserOptionPath(DefaultParserPath),
	}

	// DefaultParser is a Parser with default settings.
	DefaultParser = NewParser(DefaultParserOptions...)
)

func NewErrUnmatchedInput(input []byte) error {
	return fmt.Errorf("there are unmatched input left: %q", string(input))
}

// Parser represents a parser which use Rule's
// to parse the input.
type Parser struct {
	MaxDepth  int
	LineBreak Rule
	LineIndex []*Region
	Path      string
}

// ParserOption represents a Parser option
// which mutates Parser in a way which
// is acceptable for this option.
type ParserOption func(*Parser)

// ParserOptionMaxDepth set max depth for rule application recursion.
func ParserOptionMaxDepth(d int) ParserOption {
	return func(p *Parser) { p.MaxDepth = d }
}

// ParserOptionLineBreak set parser line-break Rule.
// Line-breaks used during error reporting,
// they are not consumed and available for otherr rules.
func ParserOptionLineBreak(r Rule) ParserOption {
	return func(p *Parser) { p.LineBreak = r }
}

// ParserOptionPath set parser path meta-information which
// is propagated to each Rule.
func ParserOptionPath(path string) ParserOption {
	return func(p *Parser) { p.Path = path }
}

// LineRegions construct a slice of Region's for given input.
// This regions contains ranges of non line-break symbols from left to right.
// Return value could be used as Parser.LineIndex.
// Also Parser.Parse calls Parser.LineRegions for you automatically.
func (p *Parser) LineRegions(input []byte) []*Region {
	loc := &Location{Path: p.Path}
	ctx := &Context{
		Parser:   p,
		Location: loc,
	}

	var (
		t       *Tree
		err     error
		regions = []*Region{}
		pos     int
		n       int
	)
	for n < len(input) {
		loc.Position = n // NOTE: affects ctx
		t, err = p.LineBreak.Parse(ctx, input[n:])
		if err == nil {
			regions = append(regions, &Region{
				Start: pos,
				End:   t.Region.Start,
			})
			pos = t.Region.End
			n = t.Region.End
		} else {
			switch err.(type) {
			case *ErrUnexpectedEOF:
				break
			}
			n++
		}
	}
	if len(input[pos:]) > 0 { // no line-break at end of file
		regions = append(regions, &Region{
			Start: pos,
			End:   pos + len(input[pos:]) - 1,
		})
	}
	return regions
}

// Locate finds a line & column of the given position.
// It expects Parser.LineIndex to be a sorted slice of Region's
// of non line-break's.
// If there is no LineIndex then it returns 0, 0.
func (p *Parser) Locate(position int) (int, int) {
	var (
		il   = len(p.LineIndex)
		h, t = 0, il - 1
		l    int
		c    int
	)
	if position == 0 { // returning zero if we have no position
		return l, c
	}
	if il == 0 { // returning position if we have no line index (no line-breaks)
		return l, position
	}
	if il == 1 { // returrning position, scoped to region (one line-break in index)
		if position > p.LineIndex[0].End {
			return l, p.LineIndex[0].End
		}
		return l, position
	}

	for h <= t {
		l = (h + t) / 2
		if position >= p.LineIndex[l].Start {
			if position <= p.LineIndex[l].End {
				break
			} else {
				h = l + 1
			}
		} else {
			t = l - 1
		}
	}
	if position > p.LineIndex[l].End {
		// handle case when position is larger than available regions
		c = p.LineIndex[l].End - p.LineIndex[l].Start
	} else {
		c = position - p.LineIndex[l].Start
		if c < 0 {
			// handle case where position points to line-break
			c = 0
		}
	}

	return l, c
}

// Parse parses input with Rule's.
// Calls Parser.LineRegions and store result under Parser.LineIndex.
// Not safe for concurrent use (and not expected to be used concurrently).
func (p *Parser) Parse(r Rule, input []byte) (*Tree, error) {
	if r == nil {
		return nil, NewErrEmptyRule(r, nil)
	}

	p.LineIndex = p.LineRegions(input)
	loc := &Location{Path: p.Path}
	tree, err := r.Parse(&Context{
		Parser:   p,
		Location: loc,
	}, input)
	if err != nil {
		if err == ErrSkipRule {
			return nil, NewErrUnexpectedEOF(r, loc)
		}
		return nil, err
	}

	if tree.Region.End < utf8.RuneCount(input) {
		pos := tree.Region.End
		line, col := p.Locate(pos)
		return nil, NewErrUnexpectedToken(
			r,
			&Location{
				Path:     p.Path,
				Position: pos,
				Line:     line,
				Column:   col,
			},
			ShowInput(input[pos:]),
			NewErrUnmatchedInput(input[tree.Region.End:]),
		)
	}

	return tree, nil
}

// Parse is a shortcut to call the DefaultParser.Parse().
func Parse(rule Rule, input []byte) (*Tree, error) {
	return DefaultParser.Parse(rule, input)
}

// NewParser constructs new *Parser.
func NewParser(op ...ParserOption) *Parser {
	p := &Parser{}
	for _, fn := range DefaultParserOptions {
		fn(p)
	}
	for _, fn := range op {
		fn(p)
	}
	return p
}

//

type Context struct {
	Rule     Rule
	Parser   *Parser
	Location *Location
	Depth    int
}

// Location represents position in input (posirion, line, column).
type Location struct {
	Path     string
	Position int
	Line     int
	Column   int
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Path, l.Line+1, l.Column+1)
}
