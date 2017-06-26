package parse

// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
