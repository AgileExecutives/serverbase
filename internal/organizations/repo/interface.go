package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// OrganizationRepo defines persistence operations used by OrganizationService
type OrganizationRepo interface {
	Create(ctx context.Context, o *models.Organization) error
	GetByID(ctx context.Context, id, tenantID uint) (*models.Organization, error)
	ListByTenant(ctx context.Context, offset, limit int, tenantID uint) ([]models.Organization, int64, error)
	Update(ctx context.Context, o *models.Organization) error
	Delete(ctx context.Context, id, tenantID uint) error
}
