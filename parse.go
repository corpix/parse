package parse

import (
	"unicode/utf8"
)

var (
	DefaultParserMaxDepth  = 128
	DefaultParserLineBreak = NewEither(
		"line-break",
		NewTerminal("lf", "\n"),
		NewTerminal("crlf", "\r\n"),
	)
	DefaultParserOptions = []ParserOption{
		ParserOptionMaxDepth(DefaultParserMaxDepth),
		ParserOptionLineBreak(DefaultParserLineBreak),
	}

	// DefaultParser is a Parser with default settings.
	DefaultParser = NewParser(DefaultParserOptions...)
)

// Parser represents a parser which use Rule's
// to parse the input.
type Parser struct {
	MaxDepth  int
	LineBreak Rule
	LineIndex []*Region
}

// ParserOption represents a Parser option
// which mutates Parser in a way which
// is acceptable for this option.
type ParserOption func(*Parser)

func ParserOptionMaxDepth(d int) ParserOption {
	return func(p *Parser) { p.MaxDepth = d }
}

func ParserOptionLineBreak(r Rule) ParserOption {
	return func(p *Parser) { p.LineBreak = r }
}

func (p *Parser) LineRegions(input []byte) []*Region {
	loc := &Location{}
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
func (p *Parser) Parse(r Rule, input []byte) (*Tree, error) {
	if r == nil {
		return nil, NewErrEmptyRule(r, nil)
	}

	p.LineIndex = p.LineRegions(input)
	tree, err := r.Parse(&Context{
		Parser:   p,
		Location: &Location{},
	}, input)
	if err != nil {
		if err == ErrSkipRule {
			return nil, NewErrUnexpectedEOF(1, r)
		}
		return nil, err
	}

	if tree.Region.End < utf8.RuneCount(input) {
		return nil, NewErrUnexpectedToken(
			ShowInput(input[tree.Region.End:]),
			tree.Region.End,
			r,
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
