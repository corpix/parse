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
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalk(t *testing.T) {
	samples := []struct {
		tree   *Tree
		data   []string
		walker func(*Tree, func(*Tree))
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
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			data: []string{
				"numbers",
				"number",
				"number",
				"1",
				"2",
			},
			walker: WalkBFS,
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
			data: []string{
				"numbers",
				"number",
				"1",
			},
			walker: WalkBFS,
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
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
				},
			},
			data: []string{
				"numbers",
				"number",
				"number",
				"number",
				"1",
				"2",
				"3",
			},
			walker: WalkBFS,
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
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("2")}},
					},
				},
			},
			data: []string{
				"numbers",
				"number",
				"1",
				"number",
				"2",
			},
			walker: WalkDFS,
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
			data: []string{
				"numbers",
				"number",
				"1",
			},
			walker: WalkDFS,
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
					{
						Data:   []byte("number"),
						Childs: []*Tree{{Data: []byte("3")}},
					},
				},
			},
			data: []string{
				"numbers",
				"number",
				"1",
				"number",
				"2",
				"number",
				"3",
			},
			walker: WalkDFS,
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
			data: []string{
				"data",
				"numbers",
				"number",
				"1",
				"number",
				"2",
				"number",
				"3",
				"colon",
			},
			walker: WalkDFS,
		},
	}
	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		buf := []string{}
		sample.walker(
			sample.tree,
			func(tree *Tree) { buf = append(buf, string(tree.Data)) },
		)
		assert.EqualValues(t, sample.data, buf, msg)
	}
}
