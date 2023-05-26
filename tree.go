package parse

import (
	"fmt"
)

// Tree represents a single Rule match with corresponding
// information about the input, position and matched Rule.
// It will be recursive in case of nested Rule match.
type Tree struct {
	Rule     Rule
	Location *Location
	Region   *Region
	Depth    int
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


// Hash produces a lication which is believed
// to uniquely identify the node in the tree.
// This is useful for serialization and graphing.
func (t *Tree) Hash() string {
	return fmt.Sprintf("%s:%d", t.Location.String(), t.Depth)
}

// Graph produce Graphviz compatible code
// which could be converted to picture using, for example
// `dot -Tpng > graph.png`.
func (t *Tree) Graph() string {
	s := "digraph G {\n"
	s += "  "
	s += fmt.Sprintf("%q", t.Hash())
	s += fmt.Sprintf(`[label=%q];`, string(t.Data))
	s += "\n"
	s += t.graph()
	s += "}"
	return s
}

// graph is a helper for Graph, and called in non root nodes.
func (t *Tree) graph() string {
	var s string
	for _, child := range t.Childs {
		s += "  "
		s += fmt.Sprintf("%q", child.Hash())
		s += fmt.Sprintf(`[label=%q];`, string(child.Data))
		s += "\n"
		s += "  "
		s += fmt.Sprintf("%q", t.Hash()) + "->" + fmt.Sprintf("%q", child.Hash())
		s += fmt.Sprintf(`[label=%q];`, child.Hash()+": "+child.Name())
		s += "\n"
		s += child.graph()
	}
	return s
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
				"rule: %q\nstart: %d\nend: %d\nlocation: %q\ndata: %q",
				treerIndentCharacter,
				treerIndentSize,
			),
			rule,
			tree.Region.Start,
			tree.Region.End,
			tree.Location.String(),
			string(tree.Data),
		),
		childs,
	)
}
