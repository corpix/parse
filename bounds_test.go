package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestBoundsGetBounds(t *testing.T) {
	samples := []struct {
		pos    int
		ts     []byte
		start  []byte
		end    []byte
		result *Bounds
		err    error
	}{
		{
			0,
			[]byte{},
			[]byte{'('},
			[]byte{')'},
			&Bounds{0, 0},
			nil,
		},
		{
			0,
			[]byte("Array(1)"),
			[]byte{'('},
			[]byte{')'},
			&Bounds{6, 7},
			nil,
		},
		{
			0,
			[]byte("Array()"),
			[]byte{'('},
			[]byte{')'},
			&Bounds{6, 6},
			nil,
		},
		{
			0,
			[]byte("Array(Array(1))"),
			[]byte{'('},
			[]byte{')'},
			&Bounds{6, 14},
			nil,
		},
		{
			0,
			[]byte("Array("),
			[]byte{'('},
			[]byte{')'},
			nil,
			NewErrBoundIncomplete(5, []byte{'('}, []byte{')'}),
		},
		{
			0,
			[]byte("Array(1"),
			[]byte{'('},
			[]byte{')'},
			nil,
			NewErrBoundIncomplete(6, []byte{'('}, []byte{')'}),
		},

		//

		{
			0,
			[]byte{},
			[]byte("begin"),
			[]byte("end"),
			&Bounds{0, 0},
			nil,
		},
		{
			0,
			[]byte("func begin 1 end"),
			[]byte("begin"),
			[]byte("end"),
			&Bounds{10, 13},
			nil,
		},
		{
			0,
			[]byte("func begin end"),
			[]byte("begin"),
			[]byte("end"),
			&Bounds{10, 11},
			nil,
		},
		{
			0,
			[]byte("func beginend"),
			[]byte("begin"),
			[]byte("end"),
			&Bounds{10, 10},
			nil,
		},
		{
			0,
			[]byte("func begin func begin 1 end end"),
			[]byte("begin"),
			[]byte("end"),
			&Bounds{10, 28},
			nil,
		},
		{
			0,
			[]byte("func begin"),
			[]byte("begin"),
			[]byte("end"),
			nil,
			NewErrBoundIncomplete(9, []byte("begin"), []byte("end")),
		},
		{
			0,
			[]byte("func begin 1"),
			[]byte("begin"),
			[]byte("end"),
			nil,
			NewErrBoundIncomplete(11, []byte("begin"), []byte("end")),
		},
	}

	var (
		bounds *Bounds
		err    error
	)

	for k, sample := range samples {
		bounds, err = GetBounds(
			sample.pos,
			sample.ts,
			sample.start,
			sample.end,
		)
		assert.Equal(
			t,
			sample.err,
			err,
			spew.Sdump(k, sample),
		)
		if sample.err == nil {
			assert.Equal(
				t,
				sample.result,
				bounds,
				spew.Sdump(k, sample),
			)
		} else {
			assert.Equal(
				t,
				(*Bounds)(nil),
				bounds,
				spew.Sdump(k, sample),
			)
		}
	}
}

func TestGetTreeBounds(t *testing.T) {
	samples := []struct {
		tree   *Tree
		bounds *Bounds
	}{
		{
			&Tree{},
			&Bounds{0, 0},
		},
		{
			&Tree{
				Start: 0,
				End:   0,
			},
			&Bounds{0, 0},
		},
		{
			&Tree{
				Start: 1,
				End:   25,
			},
			&Bounds{1, 25},
		},
	}

	for k, sample := range samples {
		bounds := GetTreeBounds(sample.tree)
		assert.Equal(
			t,
			sample.bounds,
			bounds,
			spew.Sdump(k, sample),
		)
	}
}

func TestGetTreesBounds(t *testing.T) {
	samples := []struct {
		trees  []*Tree
		bounds *Bounds
	}{
		{
			[]*Tree{},
			&Bounds{0, 0},
		},
		{
			[]*Tree{
				{
					Start: 0,
					End:   0,
				},
			},
			&Bounds{0, 0},
		},
		{
			[]*Tree{
				{
					Start: 0,
					End:   5,
				},
			},
			&Bounds{0, 5},
		},
		{
			[]*Tree{
				{
					Start: 0,
					End:   3,
				},
				{
					Start: 3,
					End:   9,
				},
			},
			&Bounds{0, 9},
		},
		{
			[]*Tree{
				{
					Start: 0,
					End:   3,
				},
				{
					Start: -1,
					End:   999,
				},
				{
					Start: 3,
					End:   9,
				},
			},
			&Bounds{0, 9},
		},
	}

	for k, sample := range samples {
		bounds := GetTreesBounds(sample.trees)
		assert.Equal(
			t,
			sample.bounds,
			bounds,
			spew.Sdump(k, sample),
		)
	}
}
