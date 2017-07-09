package parse


import (
	"reflect"
)

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

func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
