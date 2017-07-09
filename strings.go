package parse


import (
	"strings"
)

func indent(s string, character string, size int) string {
	indent := strings.Repeat(character, size)
	lines := strings.Split(s, newLine)
	for k, v := range lines {
		if len(v) > 0 {
			lines[k] = indent + v
		}
	}
	return strings.Join(lines, newLine)
}

// EqualSlicesFoldPrefix checks that a is prefixed or equal with prefix
// via EqualSlicesFold.
func EqualSlicesFoldPrefix(a, prefix []string) bool {
	l, lprefix := len(a), len(prefix)
	switch {
	case lprefix == 0:
		return true
	case l < lprefix:
		return false
	case l > lprefix:
		return EqualSlicesFold(a[:lprefix], prefix)
	default:
		return EqualSlicesFold(a, prefix)
	}
}

// EqualSlicesFoldSuffix checks that a is suffixed or equal with suffix
// via EqualSlicesFold.
func EqualSlicesFoldSuffix(a, suffix []string) bool {
	l, lsuffix := len(a), len(suffix)
	switch {
	case lsuffix == 0:
		return true
	case l < lsuffix:
		return false
	case l > lsuffix:
		return EqualSlicesFold(a[l-lsuffix:], suffix)
	default:
		return EqualSlicesFold(a, suffix)
	}
}

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
