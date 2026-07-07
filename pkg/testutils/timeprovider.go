package testutils

import "time"

// TimeProvider is a simple interface to obtain the current time.
type TimeProvider interface {
	Now() time.Time
}

// FixedTimeProvider always returns the configured time.
type FixedTimeProvider struct {
	t time.Time
}

// NewFixedTimeProvider creates a FixedTimeProvider.
func NewFixedTimeProvider(t time.Time) *FixedTimeProvider {
	return &FixedTimeProvider{t: t}
}

// Now returns the configured time.
func (f *FixedTimeProvider) Now() time.Time { return f.t }
