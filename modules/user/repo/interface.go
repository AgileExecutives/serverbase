package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// UserRepo defines the narrow data-access methods services should depend on.
// Keep this small and focused: services should use only the operations they need.
type UserRepo interface {
	// FindByID returns the user for the given numeric id.
	FindByID(ctx context.Context, id uint) (*models.User, error)

	// FindByEmail looks up a user by email address.
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// Save persists the user entity (create or update).
	Save(ctx context.Context, u *models.User) error
}

// NewsletterRepo persists newsletter subscriptions.
type NewsletterRepo interface {
	SaveNewsletter(ctx context.Context, n *models.Newsletter) error
}

// TokenBlacklistRepo persists token blacklist entries used for logout.
type TokenBlacklistRepo interface {
	SaveTokenBlacklist(ctx context.Context, tb *models.TokenBlacklist) error
}
