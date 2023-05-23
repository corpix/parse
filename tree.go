package parse

import (
	"fmt"
)

// Location represents position in input (posirion, line, column) & ast (depth).
type Location struct {
	Position int
	Line     int
	Column   int
	Depth    int
}

func (l *Location) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Column)
}

// Tree represents a single Rule match with corresponding
// information about the input, position and matched Rule.
// It will be recursive in case of nested Rule match.
type Tree struct {
	Rule     Rule
	Location *Location
	Region   *Region
	Childs   []*Tree
	Data     []byte
}

// Name returns current node name.
func (t *Tree) Name() string {
	if t.Rule != nil {
		return t.Rule.Name()
	}
	return ""
}

// GetChilds returns a slice of the Treer with children nodes.
func (t *Tree) GetChilds() Treers {
	treer := make(Treers, len(t.Childs))
	for k, v := range t.Childs {
		treer[k] = v
	}
	return treer
}

// Show this node as a string.
// You should provide childs as string
// to this function, it does not care
// about nesting in a tree, it only shows
// string representation of itself.
func (t *Tree) Show(childs string) string {
	var r string

	if len(t.Rule.GetChilds()) == 0 {
		r = t.Rule.Show("")
	} else {
		r = t.Rule.Show("...")
	}

	return TreeShow(t, r, childs)
}

func (t *Tree) String() string {
	return TreerString(t)
}

//

// TreeShow returns a Tree encoded as a string.
// It requires some parts to be prepared(encoded into a string).
func TreeShow(tree *Tree, rule string, childs string) string {
	return fmt.Sprintf(
		"%s{\n%s\n}(%s)",
		tree.Name(),
		fmt.Sprintf(
			indent(
				"rule: %s\nstart: %d\nend: %d\ndata: %s",
				treerIndentCharacter,
				treerIndentSize,
			),
			rule,
			tree.Region.Start,
			tree.Region.End,
			string(tree.Data),
		),
		childs,
	)
}
