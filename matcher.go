package parse


// Matcher represents an interface of the object
// which could be checked to match to some slice of strings.
type Matcher interface {
	Match(chain []string) bool
}

//

// AllMatcher represents a Matcher which wraps a slice of Matcher
// where each Matcher should match for the wrapper to be true.
type AllMatcher []Matcher

// Match checks that every Matcher.Match(...) returns true
// and returns true too, otherwise it returns false.
// If no Matcher presented in a slice then it will
// return false.
func (m AllMatcher) Match(chain []string) bool {
	if len(m) == 0 {
		return false
	}
	for _, v := range m {
		if !v.Match(chain) {
			return false
		}
	}
	return true
}

// NewAllMatcher creates new AllMatcher.
func NewAllMatcher(matchers []Matcher) AllMatcher {
	return AllMatcher(matchers)
}

//

// SomeMatcher represents a Matcher which wraps a slice of Matcher
// where some(one of) Matcher should match for the wrapper to be true.
type SomeMatcher []Matcher

// Match checks that some(one of) Matcher.Match(...) returns true
// and returns true too, otherwise it returns false.
// If no Matcher presented in a slice then it will
// return false.
func (m SomeMatcher) Match(chain []string) bool {
	for _, v := range m {
		if v.Match(chain) {
			return true
		}
	}
	return false
}

// NewSomeMatcher creates new SomeMatcher.
func NewSomeMatcher(matchers []Matcher) SomeMatcher {
	return SomeMatcher(matchers)
}

//

// PrefixMatcher represents a Matcher which is true
// when chain(which should be matched) has the specified
// prefix.
type PrefixMatcher []string

// Match checks chain has prefix m.
func (m PrefixMatcher) Match(chain []string) bool {
	return EqualSlicesFoldPrefix(chain, m)
}

// NewPrefixMatcher creates new PrefixMatcher.
func NewPrefixMatcher(prefix []string) PrefixMatcher {
	return PrefixMatcher(prefix)
}

//

// SuffixMatcher represents a Matcher which is true
// when chain(which should be matched) has the specified
// suffix.
type SuffixMatcher []string

// Match checks chain has suffix m.
func (m SuffixMatcher) Match(chain []string) bool {
	return EqualSlicesFoldSuffix(chain, m)
}

// NewSuffixMatcher creates new SuffixMatcher.
func NewSuffixMatcher(suffix []string) SuffixMatcher {
	return SuffixMatcher(suffix)
}

//

// EqualMatcher represents a Matcher which is true
// when chain(which should be matched) has same elements
// as matcher.
type EqualMatcher []string

// Match checks chain has same elements as m.
func (m EqualMatcher) Match(chain []string) bool {
	return EqualSlicesFold(m, chain)
}

// NewEqualMatcher creates new EqualMatcher.
func NewEqualMatcher(sample []string) EqualMatcher {
	return EqualMatcher(sample)
}

//

// LengthMatcher represents a Matcher which is true
// when chain(which should be matched) has same length
// as int stored in matcher.
type LengthMatcher int

// Match checks chain length is same as int(m).
func (m LengthMatcher) Match(chain []string) bool {
	return int(m) == len(chain)
}

// NewLengthMatcher creates new LengthMatcher.
func NewLengthMatcher(length int) LengthMatcher {
	return LengthMatcher(length)
}
