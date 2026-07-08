package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// TenantRepo is a minimal repository interface for tenant operations.
type TenantRepo interface {
	// FindByID returns the tenant by numeric id.
	FindByID(ctx context.Context, id uint) (*models.Tenant, error)

	// FindBySlug returns a tenant by slug.
	FindBySlug(ctx context.Context, slug string) (*models.Tenant, error)

	// Save persists tenant entity.
	Save(ctx context.Context, t *models.Tenant) error
	// List returns all tenants.
	List(ctx context.Context) ([]models.Tenant, error)
}
