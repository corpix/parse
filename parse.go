package parse

import (
	"bytes"
	"unicode/utf8"
)

var (
	// DefaultParser is a Parser with default settings.
	DefaultParser = NewParser(128)
)

// Parser represents a parser which use Rule's
// to parse the input.
type Parser struct {
	maxDepth int
}

// Parse parses input with Rule's.
func (p *Parser) Parse(rule Rule, input []byte) (*Tree, error) {
	tree, err := p.parse(rule, input, nil, 0, 0)
	if err != nil {
		if err == ErrSkipRule {
			return nil, NewErrUnexpectedEOF(1, rule)
		}
		return nil, err
	}

	if tree.End < utf8.RuneCount(input) {
		return nil, NewErrUnexpectedToken(
			ShowInput(input[tree.End:]),
			p.humanizePosition(tree.End),
			rule,
		)
	}

	return tree, nil
}

func (p *Parser) parse(rule Rule, input []byte, parent Rule, position int, depth int) (*Tree, error) {
	var (
		tree    *Tree
		subTree *Tree
		buf     []byte
		bounds  *Bounds
		err     error
	)

	if depth > p.maxDepth {
		return nil, NewErrNestingTooDeep(
			depth,
			p.humanizePosition(position),
		)
	}

	tree = &Tree{}

	if rule == nil {
		return nil, NewErrEmptyRule(rule, parent)
	}

	switch v := rule.(type) {
	case *Terminal:
		length := utf8.RuneCount(v.Value)
		if length == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		if utf8.RuneCount(input) < length {
			return nil, NewErrUnexpectedEOF(
				p.humanizePosition(position),
				v,
			)
		}
		buf = input[:length]

		if !bytes.EqualFold(buf, v.Value) {
			return nil, NewErrUnexpectedToken(
				ShowInput(input),
				p.humanizePosition(position),
				v,
			)
		}

		tree.Start = position
		tree.End = position + length
		tree.Rule = v
		tree.Data = buf
	case *Either:
		if len(v.Rules) == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		for _, r := range v.Rules {
			subTree, err = p.parse(
				r,
				input,
				v,
				position,
				depth+1,
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
		if err != nil {
			return nil, err
		}

		tree.Childs = make([]*Tree, 1)
		tree.Childs[0] = subTree
		tree.Start = subTree.Start
		tree.End = subTree.End
		tree.Rule = v
		tree.Data = input[:subTree.End-subTree.Start]
	case *Chain:
		if len(v.Rules) == 0 {
			return nil, NewErrEmptyRule(v, parent)
		}

		buf = input
		for _, r := range v.Rules {
			subTree, err = p.parse(
				r,
				buf,
				v,
				position,
				depth+1,
			)
			if err != nil {
				if err == ErrSkipRule {
					continue
				}
				return nil, err
			}
			position = subTree.End
			buf = buf[subTree.End-subTree.Start:]

			tree.Childs = append(
				tree.Childs,
				subTree,
			)
		}

		bounds = GetTreesBounds(tree.Childs)
		tree.Start = bounds.Starting
		tree.End = bounds.Closing
		tree.Rule = v
		tree.Data = input[:bounds.Closing-bounds.Starting]
	case *Repetition:
		seen := 0
		buf = input

	repetitionLoop:
		for {
			subTree, err = p.parse(
				v.Rule,
				buf,
				v,
				position,
				depth+1,
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
					if seen == 0 && v.Times == 0 && v.Variadic {
						return nil, ErrSkipRule
					}
					break repetitionLoop
				default:
					return nil, err
				}
			}
			seen++

			movePos := subTree.End - subTree.Start
			if !v.Variadic && seen > v.Times {
				return nil, NewErrUnexpectedToken(
					ShowInput(input[position:]),
					position+movePos,
					v,
				)
			}

			position += movePos
			buf = buf[movePos:]

			tree.Childs = append(
				tree.Childs,
				subTree,
			)
		}
		if err != nil && len(tree.Childs) == 0 {
			return nil, err
		}
		if seen < v.Times {
			if err != nil {
				return nil, err
			}
			return nil, NewErrUnexpectedToken(
				input,
				position,
				v,
			)
		}

		bounds = GetTreesBounds(tree.Childs)
		tree.Start = bounds.Starting
		tree.End = bounds.Closing
		tree.Rule = v
		tree.Data = input[:bounds.Closing-bounds.Starting]
	case *Wrapper:
		if v.Rule == nil {
			return nil, NewErrEmptyRule(v, parent)
		}

		subTree, err = p.parse(
			v.Rule,
			input,
			v,
			position,
			depth+1,
		)
		if err != nil {
			return nil, err
		}

		tree.Childs = make([]*Tree, 1)
		tree.Childs[0] = subTree
		tree.Start = subTree.Start
		tree.End = subTree.End
		tree.Rule = v
		tree.Data = input[:subTree.End-subTree.Start]
	default:
		return nil, NewErrUnsupportedRule(rule)
	}

	return tree, nil
}

// humanizePosition is just a little helper which wraps
// common position operations before it will be showed to
// human.
func (p *Parser) humanizePosition(position int) int {
	return position + 1
}

// Parse is a shortcut to call the DefaultParser.Parse().
func Parse(rule Rule, input []byte) (*Tree, error) {
	return DefaultParser.Parse(rule, input)
}

// NewParser constructs new *Parser.
func NewParser(maxDepth int) *Parser {
	return &Parser{maxDepth}
}
