package testutils

import (
	"testing"
)

func TestNewTestContext_Basic(t *testing.T) {
	_, ctx, w := NewTestContext("GET", "/ping", []byte("{}"))
	if ctx.Request == nil {
		t.Fatalf("expected request in context")
	}
	if w == nil {
		t.Fatalf("expected ResponseRecorder")
	}
}
