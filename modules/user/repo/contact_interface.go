package repo

import (
	"context"

	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/pkg/models"
)

// ContactRepo defines data access for contacts and newsletter subscriptions.
type ContactRepo interface {
	// ListContacts returns a slice of contacts and the total count matching filters.
	ListContacts(ctx context.Context, offset, limit int, active *bool, contactType string) ([]models.Contact, int64, error)

	// FindByID returns a contact by numeric id.
	FindByID(ctx context.Context, id uint) (*models.Contact, error)

	// CreateContact persists a new contact.
	CreateContact(ctx context.Context, c *models.Contact) error

	// UpdateContact updates an existing contact.
	UpdateContact(ctx context.Context, c *models.Contact) error

	// DeleteContact deletes the provided contact.
	DeleteContact(ctx context.Context, c *models.Contact) error

	// UpsertNewsletter inserts or updates a newsletter subscription. Returns true if subscription present/updated.
	UpsertNewsletter(ctx context.Context, n *basemodels.Newsletter) (bool, error)

	// ListNewsletters returns all newsletter subscriptions.
	ListNewsletters(ctx context.Context) ([]basemodels.Newsletter, error)

	// DeleteNewsletterByEmail deletes a newsletter subscription by email and returns rows affected.
	DeleteNewsletterByEmail(ctx context.Context, email string) (int64, error)
}
