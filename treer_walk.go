package parse


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
