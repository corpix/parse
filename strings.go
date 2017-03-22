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
	"strings"
)

// EqualSlicesFold checks that all elements from a
// exists in b with strings.EqualFold.
func EqualSlicesFold(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !strings.EqualFold(a[k], b[k]) {
			return false
		}
	}
	return true
}

// EqualSlicesFoldSome checks that all elements from a
// exists in one of b items with EqualSlicesFold.
func EqualSlicesFoldSome(a []string, b ...[]string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	for k := range b {
		if EqualSlicesFold(a, b[k]) {
			return true
		}
	}
	return false
}
