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
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestWalkTreer(t *testing.T) {
	type leveledTreer struct {
		Treer
		Level int
	}

	var (
		appendWalker = func(buf *Treers) func(int, Treer) error {
			return func(level int, t Treer) error {
				*buf = append(*buf, &leveledTreer{t, level})
				return nil
			}
		}
		onceWalker = func(buf *Treers) func(int, Treer) error {
			return func(level int, t Treer) error {
				*buf = append(*buf, &leveledTreer{t, level})
				return ErrStopIteration
			}
		}
		skipChildsWhenData = func(data string) func(buf *Treers) func(int, Treer) error {
			return func(buf *Treers) func(int, Treer) error {
				return func(level int, t Treer) error {
					*buf = append(*buf, &leveledTreer{t, level})
					if bytes.EqualFold(
						[]byte(data),
						t.(*Tree).Data,
					) {
						return ErrSkipBranch
					}
					return nil
				}
			}
		}
	)

	samples := []struct {
		tree   Treer
		trees  Treers
		walker func(*Treers) func(int, Treer) error
		walkFn func(Treer, func(int, Treer) error) error
		err    error
	}{
		// BFS

		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data: []byte("1"),
					},
					2,
				},
			},
			walkFn: WalkTreerBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
						},
					},
					0,
				},
			},
			walkFn: WalkTreerBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 2},
				&leveledTreer{&Tree{Data: []byte("2")}, 2},
			},
			walkFn: WalkTreerBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
						},
					},
					0,
				},
			},
			walkFn: WalkTreerBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number1"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data:   []byte("number2"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number1"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number2"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number1"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number2"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("2")}, 2},
			},
			walkFn: WalkTreerBFS,
			walker: skipChildsWhenData("number1"),
		},

		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data: []byte("number"),
						Childs: []*Tree{
							{
								Data:   []byte("2.1"),
								Childs: []*Tree{{Data: []byte("2.1.1")}},
							},
							{Data: []byte("2.2")},
						},
					},
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data: []byte("number"),
								Childs: []*Tree{
									{
										Data:   []byte("2.1"),
										Childs: []*Tree{{Data: []byte("2.1.1")}},
									},
									{Data: []byte("2.2")},
								},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("3")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data: []byte("number"),
						Childs: []*Tree{
							{
								Data:   []byte("2.1"),
								Childs: []*Tree{{Data: []byte("2.1.1")}},
							},
							{Data: []byte("2.2")},
						},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 2},
				&leveledTreer{
					&Tree{
						Data:   []byte("2.1"),
						Childs: []*Tree{{Data: []byte("2.1.1")}},
					},
					2,
				},
				&leveledTreer{&Tree{Data: []byte("2.2")}, 2},
				&leveledTreer{&Tree{Data: []byte("3")}, 2},
				&leveledTreer{&Tree{Data: []byte("2.1.1")}, 3},
			},
			walkFn: WalkTreerBFS,
			walker: appendWalker,
		},

		// DFS

		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 2},
			},
			walkFn: WalkTreerDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
						},
					},
					0,
				},
			},
			walkFn: WalkTreerDFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 2},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("2")}, 2},
			},
			walkFn: WalkTreerDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number1"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data:   []byte("number2"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number1"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number2"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number1"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number2"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("2")}, 2},
			},
			walkFn: WalkTreerDFS,
			walker: skipChildsWhenData("number1"),
		},
		{
			tree: &Tree{
				Data: []byte("numbers"),
				Childs: []*Tree{
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					{
						Data: []byte("number"),
						Childs: []*Tree{
							{
								Data:   []byte("2.1"),
								Childs: []*Tree{{Data: []byte("2.1.1")}},
							},
							{Data: []byte("2.2")},
						},
					},
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data: []byte("number"),
								Childs: []*Tree{
									{
										Data:   []byte("2.1"),
										Childs: []*Tree{{Data: []byte("2.1.1")}},
									},
									{Data: []byte("2.2")},
								},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("3")}},
							},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 2},
				&leveledTreer{
					&Tree{
						Data: []byte("number"),
						Childs: []*Tree{
							{
								Data:   []byte("2.1"),
								Childs: []*Tree{{Data: []byte("2.1.1")}},
							},
							{Data: []byte("2.2")},
						},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("2.1"),
						Childs: []*Tree{{Data: []byte("2.1.1")}},
					},
					2,
				},
				&leveledTreer{&Tree{Data: []byte("2.1.1")}, 3},
				&leveledTreer{&Tree{Data: []byte("2.2")}, 2},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
					1,
				},
				&leveledTreer{&Tree{Data: []byte("3")}, 2},
			},
			walkFn: WalkTreerDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Data: []byte("data"),
				Childs: []*Tree{
					{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("3")}},
							},
						},
					},
					{Data: []byte("colon")},
				},
			},
			trees: Treers{
				&leveledTreer{
					&Tree{
						Data: []byte("data"),
						Childs: []*Tree{
							{
								Data: []byte("numbers"),
								Childs: []*Tree{
									{
										Data:   []byte("number"),
										Childs: []*Tree{{Data: []byte("1")}},
									},
									{
										Data:   []byte("number"),
										Childs: []*Tree{{Data: []byte("2")}},
									},
									{
										Data:   []byte("number"),
										Childs: []*Tree{{Data: []byte("3")}},
									},
								},
							},
							{Data: []byte("colon")},
						},
					},
					0,
				},
				&leveledTreer{
					&Tree{
						Data: []byte("numbers"),
						Childs: []*Tree{
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("1")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("2")}},
							},
							{
								Data:   []byte("number"),
								Childs: []*Tree{{Data: []byte("3")}},
							},
						},
					},
					1,
				},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("1")}},
					},
					2,
				},
				&leveledTreer{&Tree{Data: []byte("1")}, 3},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
					2,
				},
				&leveledTreer{&Tree{Data: []byte("2")}, 3},
				&leveledTreer{
					&Tree{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
					2,
				},
				&leveledTreer{&Tree{Data: []byte("3")}, 3},
				&leveledTreer{&Tree{Data: []byte("colon")}, 1},
			},
			walkFn: WalkTreerDFS,
			walker: appendWalker,
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		buf := &Treers{}
		err := sample.walkFn(
			sample.tree,
			sample.walker(buf),
		)
		assert.Equal(t, sample.err, err, msg)
		assert.EqualValues(t, sample.trees, *buf, msg)
	}
}
