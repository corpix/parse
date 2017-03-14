package parse

import ()

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

// WalkBFS walks the Tree level by level.
// See: https://en.wikipedia.org/wiki/Breadth-first_search
func WalkBFS(tree *Tree, fn func(*Tree)) {
	var (
		current *Tree
		stack   = []*Tree{}
	)

	if tree == nil {
		return
	}

	current = tree
	for {
		fn(current)
		if current.Childs != nil {
			for _, v := range current.Childs {
				stack = append(
					stack,
					v,
				)
			}
		}

		if len(stack) == 0 {
			break
		}

		current = stack[0]
		stack = stack[1:]
	}
}

// WalkDFS walks the Tree node by node to leafs.
// See: https://en.wikipedia.org/wiki/Depth-first_search
func WalkDFS(tree *Tree, fn func(*Tree)) {
	var (
		current *Tree
		level   int
		childs  []*Tree
		ok      bool
	)
	current = tree
	backlog := map[int][]*Tree{}

	for current != nil {
		fn(current)
		if current.Childs != nil && len(current.Childs) > 0 {
			level++
			backlog[level] = current.Childs[1:]
			current = current.Childs[0]
			continue
		}

	nextLevel:
		childs, ok = backlog[level]
		if ok && len(childs) > 0 {
			current = childs[0]
			backlog[level] = childs[1:]
			continue
		}

		level--
		if level < 0 {
			break
		}
		goto nextLevel
	}
}
