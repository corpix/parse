package parse


const (
	newLine              = "\n"
	treerIndentCharacter = " "
	treerDelimiter       = ", "
)

const (
	treerIndentSize = 2
)

const (
	circularLabel = "<circular>"
	nilLabel      = "<nil>"
)

// treerString folds a tree into string representation(with `Show`) while
// resolving the pointer loops.
func treerString(t Treer, visited map[interface{}]bool) string {
	var (
		childs         string
		alreadyVisited bool
	)

	_, alreadyVisited = visited[t]
	if alreadyVisited {
		return t.Show(circularLabel)
	}

	visited[t] = true

	for k, v := range t.GetChilds() {
		if k > 0 {
			childs += treerDelimiter
			childs += newLine
		}
		if v == nil {
			childs += nilLabel
			continue
		}
		childs += treerString(
			v,
			visited,
		)
	}

	if len(childs) > 0 {
		childs = newLine + indent(
			childs,
			treerIndentCharacter,
			treerIndentSize,
		) + newLine
	}

	return t.Show(childs)
}

// TreerString prints a Treer as human-readable string.
func TreerString(t Treer) string {
	return treerString(t, map[interface{}]bool{})
}
