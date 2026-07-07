package testutils

// Email is a simple record of a sent email.
type Email struct {
	To      []string
	Subject string
	Body    string
}

// MockEmailSender collects sent emails for assertions in tests.
type MockEmailSender struct {
	Sent []Email
}

// NewMockEmailSender creates a new MockEmailSender.
func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{Sent: make([]Email, 0)}
}

// Send records an email send action.
func (m *MockEmailSender) Send(to []string, subj, body string) error {
	m.Sent = append(m.Sent, Email{To: to, Subject: subj, Body: body})
	return nil
}

// Last returns the last sent email or nil.
func (m *MockEmailSender) Last() *Email {
	if len(m.Sent) == 0 {
		return nil
	}
	return &m.Sent[len(m.Sent)-1]
}
