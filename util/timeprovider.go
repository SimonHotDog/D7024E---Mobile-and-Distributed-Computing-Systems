package util

import "time"

type ITimeProvider interface {
	// Now returns the current local time.
	Now() time.Time
}

type TimeProvider struct{}

func (m *TimeProvider) Now() time.Time {
	return time.Now()
}

// FakeTimeProvider is a wrapper around the time package.
// This is done to make it possible to mock the time package in tests.
//
// When making a new FakeTimeProvider, the InternalTime must be set else Now()
// will return 0 time.
type FakeTimeProvider struct {
	InternalTime time.Time
}

func (f *FakeTimeProvider) Now() time.Time {
	return f.InternalTime
}
