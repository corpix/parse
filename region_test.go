package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestTreeRegion(t *testing.T) {
	samples := []struct {
		trees  []*Tree
		region *Region
	}{
		{
			[]*Tree{},
			&Region{0, 0},
		},
		{
			[]*Tree{
				{
					Region: &Region{
						Start: 0,
						End:   0,
					},
				},
			},
			&Region{0, 0},
		},
		{
			[]*Tree{
				{
					Region: &Region{
						Start: 0,
						End:   5,
					},
				},
			},
			&Region{0, 5},
		},
		{
			[]*Tree{
				{
					Region: &Region{
						Start: 0,
						End:   3,
					},
				},
				{
					Region: &Region{
						Start: 3,
						End:   9,
					},
				},
			},
			&Region{0, 9},
		},
		{
			[]*Tree{
				{
					Region: &Region{
						Start: 0,
						End:   3,
					},
				},
				{
					Region: &Region{
						Start: -1,
						End:   999,
					},
				},
				{
					Region: &Region{
						Start: 3,
						End:   9,
					},
				},
			},
			&Region{0, 9},
		},
	}

	for k, sample := range samples {
		region := TreeRegion(sample.trees...)
		assert.Equal(
			t,
			sample.region,
			region,
			spew.Sdump(k, sample),
		)
	}
}
