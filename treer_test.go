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
	"strings"
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

//

type TestRule string

func (t *TestRule) ID() string                            { return string(*t) }
func (t *TestRule) GetChilds() Treers                     { return nil }
func (t *TestRule) GetParameters() map[string]interface{} { return nil }
func (t *TestRule) String() string                        { return RuleString(t) }
func NewTestRule(id string) *TestRule {
	t := TestRule(id)
	return &t
}

func TestWalkTreerIDChain(t *testing.T) {
	type leveledTreer struct {
		Treer
		Chain []string
	}

	var (
		appendWalker = func(buf *[][]string) func([]string, int, Treer) error {
			return func(chain []string, level int, t Treer) error {
				*buf = append(*buf, chain)
				return nil
			}
		}
		onceWalker = func(buf *[][]string) func([]string, int, Treer) error {
			return func(chain []string, level int, t Treer) error {
				*buf = append(*buf, chain)
				return ErrStopIteration
			}
		}
		skipChildsWhenData = func(data string) func(buf *[][]string) func([]string, int, Treer) error {
			return func(buf *[][]string) func([]string, int, Treer) error {
				return func(chain []string, level int, t Treer) error {
					*buf = append(*buf, chain)
					if strings.EqualFold(
						data,
						t.ID(),
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
		chains [][]string
		walker func(*[][]string) func([]string, int, Treer) error
		walkFn func(Treer, func([]string, int, Treer) error) error
		err    error
	}{
		// BFS
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2"),
						Childs: []*Tree{{Rule: NewTestRule("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2"},
				{"1", "2", "3"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2"),
						Childs: []*Tree{{Rule: NewTestRule("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule:   NewTestRule("2.2"),
						Childs: []*Tree{{Rule: NewTestRule("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.1", "3.1"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2"),
						Childs: []*Tree{{Rule: NewTestRule("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule:   NewTestRule("2.2"),
						Childs: []*Tree{{Rule: NewTestRule("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: skipChildsWhenData("2.1"),
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule: NewTestRule("2.2"),
						Childs: []*Tree{
							{
								Rule: NewTestRule("3.2"),
								Childs: []*Tree{
									{Rule: NewTestRule("3.2.1")},
								},
							},
						},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.1", "3.1"},
				{"1", "2.2", "3.2"},
				{"1", "2.2", "3.2", "3.2.1"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: skipChildsWhenData("3.1"),
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule: NewTestRule("2.2"),
						Childs: []*Tree{
							{
								Rule: NewTestRule("3.2"),
								Childs: []*Tree{
									{Rule: NewTestRule("3.2.1")},
								},
							},
						},
					},
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.2"},
				{"1", "2.1"},
				{"1", "2.2", "3.2"},
				{"1", "2.1", "3.1"},
				{"1", "2.2", "3.2", "3.2.1"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: skipChildsWhenData("3.1"),
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule: NewTestRule("2.2"),
						Childs: []*Tree{
							{
								Rule:   NewTestRule("3.2.1"),
								Childs: []*Tree{{Rule: NewTestRule("4.2")}},
							},
							{Rule: NewTestRule("3.2.2")},
						},
					},
					{
						Rule:   NewTestRule("2.3"),
						Childs: []*Tree{{Rule: NewTestRule("3.3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.3"},
				{"1", "2.1", "3.1"},
				{"1", "2.2", "3.2.1"},
				{"1", "2.2", "3.2.2"},
				{"1", "2.3", "3.3"},
				{"1", "2.2", "3.2.1", "4.2"},
			},
			walkFn: WalkTreerIDChainBFS,
			walker: appendWalker,
		},

		// DFS

		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2"),
						Childs: []*Tree{{Rule: NewTestRule("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2"},
				{"1", "2", "3"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2"),
						Childs: []*Tree{{Rule: NewTestRule("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule:   NewTestRule("2.2"),
						Childs: []*Tree{{Rule: NewTestRule("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.1", "3.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule:   NewTestRule("2.2"),
						Childs: []*Tree{{Rule: NewTestRule("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: skipChildsWhenData("2.1"),
		},
		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule:   NewTestRule("2.1"),
						Childs: []*Tree{{Rule: NewTestRule("3.1")}},
					},
					{
						Rule: NewTestRule("2.2"),
						Childs: []*Tree{
							{
								Rule:   NewTestRule("3.2.1"),
								Childs: []*Tree{{Rule: NewTestRule("4.2")}},
							},
							{Rule: NewTestRule("3.2.2")},
						},
					},
					{
						Rule:   NewTestRule("2.3"),
						Childs: []*Tree{{Rule: NewTestRule("3.3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.1", "3.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2.1"},
				{"1", "2.2", "3.2.1", "4.2"},
				{"1", "2.2", "3.2.2"},
				{"1", "2.3"},
				{"1", "2.3", "3.3"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: appendWalker,
		},

		{
			tree: &Tree{
				Rule: NewTestRule("1"),
				Childs: []*Tree{
					{
						Rule: NewTestRule("2.1"),
						Childs: []*Tree{
							{
								Rule:   NewTestRule("3.1.1"),
								Childs: []*Tree{{Rule: NewTestRule("4.1.1")}},
							},
							{
								Rule:   NewTestRule("3.1.2"),
								Childs: []*Tree{{Rule: NewTestRule("4.1.2")}},
							},
							{
								Rule:   NewTestRule("3.1.3"),
								Childs: []*Tree{{Rule: NewTestRule("4.1.3")}},
							},
						},
					},
					{Rule: NewTestRule("2.2")},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.1", "3.1.1"},
				{"1", "2.1", "3.1.1", "4.1.1"},
				{"1", "2.1", "3.1.2"},
				{"1", "2.1", "3.1.2", "4.1.2"},
				{"1", "2.1", "3.1.3"},
				{"1", "2.1", "3.1.3", "4.1.3"},
				{"1", "2.2"},
			},
			walkFn: WalkTreerIDChainDFS,
			walker: appendWalker,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		buf := &[][]string{}
		err := sample.walkFn(
			sample.tree,
			sample.walker(buf),
		)
		assert.Equal(t, sample.err, err, msg)
		assert.EqualValues(t, sample.chains, *buf, msg)
	}
}
