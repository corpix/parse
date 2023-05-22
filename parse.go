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

// Parse parses input with Rule's.
func (p *Parser) Parse(r Rule, input []byte) (*Tree, error) {
	if r == nil {
		return nil, NewErrEmptyRule(r, nil)
	}

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
