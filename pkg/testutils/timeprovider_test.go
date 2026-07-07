package testutils

import (
	"testing"
	"time"
)

func TestFixedTimeProvider_Now(t *testing.T) {
	tm := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	p := NewFixedTimeProvider(tm)
	if !p.Now().Equal(tm) {
		t.Fatalf("expected %v, got %v", tm, p.Now())
	}
}
