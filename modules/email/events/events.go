package events

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"gorm.io/gorm"
)

type EmailSentHandler struct{ db *gorm.DB }

func NewEmailSentHandler(db *gorm.DB) *EmailSentHandler    { return &EmailSentHandler{db: db} }
func (h *EmailSentHandler) EventType() string              { return "email.sent" }
func (h *EmailSentHandler) Handle(event interface{}) error { return nil }
func (h *EmailSentHandler) Priority() int                  { return 100 }

type EmailFailedHandler struct {
	db     *gorm.DB
	logger core.Logger
}

func NewEmailFailedHandler(db *gorm.DB, logger core.Logger) *EmailFailedHandler {
	return &EmailFailedHandler{db: db, logger: logger}
}
func (h *EmailFailedHandler) EventType() string { return "email.failed" }
func (h *EmailFailedHandler) Handle(event interface{}) error {
	h.logger.Error("Email sending failed", "event", event)
	return nil
}
func (h *EmailFailedHandler) Priority() int { return 100 }

type EmailQueuedEvent struct {
	EmailID uint   `json:"email_id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}
type EmailSentEvent struct {
	EmailID uint   `json:"email_id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	SentAt  string `json:"sent_at"`
}
type EmailFailedEvent struct {
	EmailID uint   `json:"email_id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}
