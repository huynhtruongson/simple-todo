package common

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
type expectedActualAssertion func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
type expectedWithinDuration func(t assert.TestingT, expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func AssertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var e asserter
	a(&e, expected, actual, msgAndArgs...)
	return e.err
}

// assertActual is a helper function to allow the step function to call
// assertion functions where you want to compare an actual value to a
// predined state like nil, empty or true/false.
func AssertActual(a expectedActualAssertion, object interface{}, msgAndArgs ...interface{}) error {
	var e asserter
	a(&e, object, msgAndArgs...)
	return e.err
}

func AssertWithinDuration(a expectedWithinDuration, expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) error {
	var e asserter
	a(&e, expected, actual, delta, msgAndArgs...)
	return e.err
}
