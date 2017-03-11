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
