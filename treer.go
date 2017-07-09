package parse


// Treer is a tree node interface.
type Treer interface {
	// Name returns current node name.
	Name() string

	// Show this node as a string.
	// You should provide childs as string
	// to this function, it does not care
	// about nesting in a tree, it only shows
	// string representation of itself.
	Show(childs string) string

	// String returns a tree string representation
	// from current node to the leafs.
	// You should use `TreerString(Treer) string` func
	// inside your implementation because it will resolve
	// the loops counting visited nodes.
	String() string

	// GetChilds returns a slice of the Treer with children nodes.
	GetChilds() Treers
}
