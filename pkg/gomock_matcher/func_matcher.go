package gomock_matcher

import (
	"fmt"
)

type Matcher interface {
	// Matches returns whether x is a match.
	Matches(x interface{}) bool

	// String describes what the matcher matches.
	String() string
}

type MatcherWithAssertFunc struct {
	AssertFunc func(got interface{}) error
	err        error
}

var _ Matcher = (*MatcherWithAssertFunc)(nil)

func (m *MatcherWithAssertFunc) Matches(x interface{}) bool {
	err := m.AssertFunc(x)
	if err != nil {
		m.err = err
		return false
	}
	return true
}

func (m *MatcherWithAssertFunc) String() string {
	if m.err != nil {
		return fmt.Sprintf("%v", m.err)
	}
	return fmt.Sprintf("matches with AssertFunc")
}

func MatchesFunc(assertFunc func(got interface{}) error) *MatcherWithAssertFunc {
	return &MatcherWithAssertFunc{
		AssertFunc: assertFunc,
	}
}
