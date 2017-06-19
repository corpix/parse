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

// WalkTreerBFS walks the Treer level by level.
// See: https://en.wikipedia.org/wiki/Breadth-first_search
func WalkTreerBFS(tree Treer, fn func(int, Treer) error) error {
	var (
		stack = Treers{}

		current Treer
		childs  Treers

		// starting from the root element
		// only one element left to jump
		// to the next level.
		currentLevelLeft = 1
		level            int
		nextLevelLen     int

		n   int
		err error
	)

	current = tree
	for current != nil {
		if currentLevelLeft == 0 {
			level++
			currentLevelLeft = nextLevelLen
			nextLevelLen = 0
		}

		err = fn(level, current)
		if err != nil {
			switch err {
			case ErrStopIteration:
				return nil
			case ErrSkipBranch:
				err = nil
				goto nextLevel
			default:
				return err
			}
		}

		childs = current.GetChilds()
		if len(childs) > 0 {
			nextLevelLen += len(childs)
			stack = append(
				stack,
				childs...,
			)
		}

	nextLevel:
		if len(stack) == 0 {
			break
		}
		current = stack[0]
		stack = stack[1:]
		n++
		currentLevelLeft--
	}

	return nil
}

// WalkTreerNameChainBFS is a walker which reports nesting as chain of
// Treer node Name's on every iteration and uses WalkerTreerBFS.
func WalkTreerNameChainBFS(tree Treer, fn func([]string, int, Treer) error) error {
	type nodeInfo struct {
		left  int
		chain []string
	}
	var (
		childs     Treers
		childsLen  int
		chain      []string
		chainCopy  []string
		parent     Treer
		parentInfo *nodeInfo

		parents = map[Treer]Treer{}
		info    = map[Treer]*nodeInfo{}
	)

	return WalkTreerBFS(
		tree,
		func(level int, tree Treer) error {
			parent = parents[tree]
			childs = tree.GetChilds()
			for _, v := range childs {
				parents[v] = tree
			}

			if parent != nil {
				parentInfo = info[parent]
				chain = append(
					parentInfo.chain,
					tree.Name(),
				)

				parentInfo.left--
				if parentInfo.left == 0 {
					delete(info, parent)
					delete(parents, parent)
				}
			} else {
				chain = []string{tree.Name()}
			}
			chainCopy = make([]string, len(chain))
			copy(chainCopy, chain)

			childsLen = len(childs)
			if childsLen > 0 {
				info[tree] = &nodeInfo{
					left:  childsLen,
					chain: chainCopy,
				}
			}

			return fn(chainCopy, level, tree)
		},
	)
}

// WalkTreerDFS walks the Treer childs from top to leafs.
// See: https://en.wikipedia.org/wiki/Depth-first_search
func WalkTreerDFS(tree Treer, fn func(int, Treer) error) error {
	var (
		current Treer
		stack   Treers
		level   int
		ok      bool
		err     error
	)
	current = tree
	backlog := map[int]Treers{}

	for current != nil {
		err = fn(level, current)
		if err != nil {
			switch err {
			case ErrStopIteration:
				return nil
			case ErrSkipBranch:
				err = nil
				goto nextLevel
			default:
				return err
			}
		}

		stack = current.GetChilds()
		if len(stack) > 0 {
			level++
			backlog[level] = stack[1:]
			current = stack[0]
			continue
		}

	nextLevel:
		stack, ok = backlog[level]
		if ok && len(stack) > 0 {
			current = stack[0]
			backlog[level] = stack[1:]
			continue
		}

		level--
		if level < 0 {
			break
		}
		goto nextLevel
	}

	return nil
}

// WalkTreerNameChainDFS is a walker which reports nesting as chain of
// Treer node Name's on every iteration and uses WalkerTreerDFS.
func WalkTreerNameChainDFS(tree Treer, fn func([]string, int, Treer) error) error {
	var (
		chain         []string
		chainCopy     []string
		previousLevel int
	)

	return WalkTreerDFS(
		tree,
		func(level int, tree Treer) error {
			if level <= previousLevel {
				chain = chain[:level]
			}
			previousLevel = level

			chain = append(
				chain,
				tree.Name(),
			)
			chainCopy = make([]string, len(chain))
			copy(
				chainCopy,
				chain,
			)

			return fn(chainCopy, level, tree)
		},
	)
}
