package parse

import (
	"bytes"
)

// Bounds represents a starting and ending points of something.
// It could be a position of '(' and ')' in 'foo(bar)' for example.
type Bounds struct {
	Starting int
	Closing  int
}

// GetBounds duty is to find a bounds of starting and closing
// inside data
func GetBounds(position int, data []byte, starting []byte, closing []byte) (*Bounds, error) {
	return getBounds(
		position,
		0,
		data,
		starting,
		closing,
	)
}

func getBounds(position, depth int, data []byte, starting []byte, closing []byte) (*Bounds, error) {
	var (
		startingWasFound bool
		err              error
		bounds           *Bounds
		tail             []byte
	)

	res := &Bounds{}
	dataLength := len(data)

	for k := 0; k < dataLength; {
		tail = data[k:]
		switch {
		case bytes.HasPrefix(tail, closing):
			if startingWasFound {
				res.Closing = k
				return res, nil
			}
		case bytes.HasPrefix(tail, starting):
			if startingWasFound {
				bounds, err = getBounds(
					position+k,
					depth+1,
					tail,
					starting,
					closing,
				)
				if err != nil {
					return nil, err
				}
				k += bounds.Closing
			} else {
				res.Starting = k + len(starting)
				startingWasFound = true
			}
		}

		k++
	}

	if startingWasFound {
		return nil, NewErrBoundIncomplete(
			position+(dataLength-1),
			starting,
			closing,
		)
	}
	return res, nil
}

// GetTreeBounds constructs a *Bounds from the *Tree.
func GetTreeBounds(tree *Tree) *Bounds {
	return &Bounds{
		tree.Start,
		tree.End,
	}
}

// GetTreesBounds constructs a *Bounds from a slice of *Tree.
// It watches only first and last element.
func GetTreesBounds(trees []*Tree) *Bounds {
	if len(trees) == 0 {
		return &Bounds{0, 0}
	}

	return &Bounds{
		trees[0].Start,
		trees[len(trees)-1].End,
	}
}
