package parse

import (
	"fmt"
)

// Tree represents a single Rule match with corresponding
// information about the input, position and matched Rule.
// It will be recursive in case of nested Rule match.
type Tree struct {
	Rule   Rule
	Start  int
	End    int
	Childs []*Tree
	Data   []byte
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
	var (
		c string
		r string
	)
	if len(childs) > 0 {
		c = indent(
			childs,
			treerIndentCharacter,
			treerIndentSize,
		)
	}

	if len(t.Rule.GetChilds()) == 0 {
		r = t.Rule.Show("")
	} else {
		r = t.Rule.Show("...")
	}

	return TreeShow(t, r, c)
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
			tree.Start,
			tree.End,
			string(tree.Data),
		),
		childs,
	)
}
