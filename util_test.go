package parse

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestShowInput(t *testing.T) {
	samples := []struct {
		input  string
		output string
	}{
		{"hello", "hello"},
		{"hello\ngoodbye", "hello..."},
		{"hello\ngoodbye\ngoodbye", "hello..."},
		{"Array(FixedString(16))", "Array(FixedString(16))"},
		{
			"hello so long so long so long so long so long",
			"hello so long so long...",
		},
		{
			"hello so long so long so long",
			"hello so long so long...",
		},
		{
			"hello lo lo lo lo lo",
			"hello lo lo lo lo...",
		},
		{
			"hello lo lo lo lo",
			"hello lo lo lo lo",
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		input := make([]byte, len(sample.input))
		copy(input, []byte(sample.input))

		output := ShowInput(input)
		assert.Equal(t, []byte(sample.input), input, msg)
		assert.Equal(t, []byte(sample.output), output, msg)
	}
}
