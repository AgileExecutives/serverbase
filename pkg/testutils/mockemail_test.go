package testutils

import "testing"

func TestMockEmailSender_RecordSend(t *testing.T) {
	m := NewMockEmailSender()
	if len(m.Sent) != 0 {
		t.Fatalf("expected empty Sent, got %d", len(m.Sent))
	}

	err := m.Send([]string{"a@example.com"}, "subj", "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m.Sent) != 1 {
		t.Fatalf("expected 1 sent email, got %d", len(m.Sent))
	}

	last := m.Last()
	if last == nil {
		t.Fatalf("expected last email, got nil")
	}
	if last.Subject != "subj" {
		t.Fatalf("unexpected subject: %s", last.Subject)
	}
}
