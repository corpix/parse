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
