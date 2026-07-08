package interfaces

import (
	"context"
	"net/http"
	"time"
)

// DB is a narrow abstraction over GORM used by modules for migrations and simple ops.
type DB interface {
	AutoMigrate(dst ...interface{}) error
}

// EventBus is a minimal event bus abstraction.
type EventBus interface {
	Publish(topic string, evt any) error
	Subscribe(topic string, handler func(ctx any, payload any)) (unsubscribe func())
}

// EmailSender sends emails.
type EmailSender interface {
	Send(ctx any, to, subj, body string) error
}

// TenantProvider returns the current tenant id from a request/context.
type TenantProvider interface {
	Current(ctx context.Context) (string, error)
}

// UserRepository is a minimal user repository interface used in services.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (any, error)
	FindByEmail(ctx context.Context, email string) (any, error)
	Save(ctx context.Context, entity any) error
}

// TimeProvider supplies current time for deterministic tests.
type TimeProvider interface {
	Now() time.Time
}

// UUIDProvider produces new UUIDs/IDs.
type UUIDProvider interface {
	New() string
}

// Storage is a small key/value storage abstraction used by modules.
type Storage interface {
	Put(key string, data []byte) error
	Get(key string) ([]byte, error)
}

// HTTPClient is a small abstraction over net/http for modules that call external services.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Logger is the minimal logging API we expect modules to use.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}
