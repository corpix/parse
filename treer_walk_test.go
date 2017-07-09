package parse


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

func TestWalkTreerNameChain(t *testing.T) {
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
						t.Name(),
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
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2"},
				{"1", "2", "3"},
			},
			walkFn: WalkTreerNameChainBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerNameChainBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule:   newTestRuleFinite("2.2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.2")}},
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
			walkFn: WalkTreerNameChainBFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerNameChainBFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule:   newTestRuleFinite("2.2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerNameChainBFS,
			walker: skipChildsWhenData("2.1"),
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule: newTestRuleFinite("2.2"),
						Childs: []*Tree{
							{
								Rule: newTestRuleFinite("3.2"),
								Childs: []*Tree{
									{Rule: newTestRuleFinite("3.2.1")},
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
			walkFn: WalkTreerNameChainBFS,
			walker: skipChildsWhenData("3.1"),
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule: newTestRuleFinite("2.2"),
						Childs: []*Tree{
							{
								Rule: newTestRuleFinite("3.2"),
								Childs: []*Tree{
									{Rule: newTestRuleFinite("3.2.1")},
								},
							},
						},
					},
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
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
			walkFn: WalkTreerNameChainBFS,
			walker: skipChildsWhenData("3.1"),
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule: newTestRuleFinite("2.2"),
						Childs: []*Tree{
							{
								Rule:   newTestRuleFinite("3.2.1"),
								Childs: []*Tree{{Rule: newTestRuleFinite("4.2")}},
							},
							{Rule: newTestRuleFinite("3.2.2")},
						},
					},
					{
						Rule:   newTestRuleFinite("2.3"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.3")}},
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
			walkFn: WalkTreerNameChainBFS,
			walker: appendWalker,
		},

		// DFS

		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2"},
				{"1", "2", "3"},
			},
			walkFn: WalkTreerNameChainDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
			},
			walkFn: WalkTreerNameChainDFS,
			walker: onceWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule:   newTestRuleFinite("2.2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.2")}},
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
			walkFn: WalkTreerNameChainDFS,
			walker: appendWalker,
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule:   newTestRuleFinite("2.2"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.2")}},
					},
				},
			},
			chains: [][]string{
				{"1"},
				{"1", "2.1"},
				{"1", "2.2"},
				{"1", "2.2", "3.2"},
			},
			walkFn: WalkTreerNameChainDFS,
			walker: skipChildsWhenData("2.1"),
		},
		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule:   newTestRuleFinite("2.1"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.1")}},
					},
					{
						Rule: newTestRuleFinite("2.2"),
						Childs: []*Tree{
							{
								Rule:   newTestRuleFinite("3.2.1"),
								Childs: []*Tree{{Rule: newTestRuleFinite("4.2")}},
							},
							{Rule: newTestRuleFinite("3.2.2")},
						},
					},
					{
						Rule:   newTestRuleFinite("2.3"),
						Childs: []*Tree{{Rule: newTestRuleFinite("3.3")}},
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
			walkFn: WalkTreerNameChainDFS,
			walker: appendWalker,
		},

		{
			tree: &Tree{
				Rule: newTestRuleFinite("1"),
				Childs: []*Tree{
					{
						Rule: newTestRuleFinite("2.1"),
						Childs: []*Tree{
							{
								Rule:   newTestRuleFinite("3.1.1"),
								Childs: []*Tree{{Rule: newTestRuleFinite("4.1.1")}},
							},
							{
								Rule:   newTestRuleFinite("3.1.2"),
								Childs: []*Tree{{Rule: newTestRuleFinite("4.1.2")}},
							},
							{
								Rule:   newTestRuleFinite("3.1.3"),
								Childs: []*Tree{{Rule: newTestRuleFinite("4.1.3")}},
							},
						},
					},
					{Rule: newTestRuleFinite("2.2")},
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
			walkFn: WalkTreerNameChainDFS,
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
