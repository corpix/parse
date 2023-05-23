package parse

// Region represents a starting and ending points of something.
// It could be a position of '(' and ')' in 'foo(bar)' for example.
type Region struct {
	Start int
	End   int
}

// TreeRegion constructs a *Region from the one or more *Tree.
func TreeRegion(tree ...*Tree) *Region {
	if len(tree) == 0 {
		return &Region{0, 0}
	}

	return &Region{
		tree[0].Region.Start,
		tree[len(tree)-1].Region.End,
	}
}
