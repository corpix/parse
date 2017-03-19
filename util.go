package parse

import ()

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

// ShowInput return a byte slice of input and formats it with ellipsis,
// preparing it to be printed to human.
func ShowInput(buf []byte) []byte {
	var (
		bufCopy            []byte
		n                  int
		hardLimit          = 32
		spaceBreakTreshold = hardLimit / 2
		ellipsis           = []byte("...")
	)

	for k, v := range buf {
		switch {
		case v == '\n':
		case v == ' ' && k > spaceBreakTreshold:
		default:
			continue
		}
		n = k + len(ellipsis)
		break
	}
	if n == 0 {
		bufCopy = make([]byte, len(buf))
		copy(bufCopy, buf)
		return buf
	}

	if n > hardLimit {
		n = hardLimit
	}

	bufCopy = make([]byte, n)
	copy(bufCopy, buf)

	for k, v := range ellipsis {
		bufCopy[n-k-1] = v
	}

	return bufCopy
}
