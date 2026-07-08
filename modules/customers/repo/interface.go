package repo

import (
	"context"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
)

// CustomerRepo defines the repository responsibilities for customer entities.
type CustomerRepo interface {
	FindByID(ctx context.Context, id uint) (*models.Customer, error)
	FindByEmail(ctx context.Context, email string) (*models.Customer, error)
	Save(ctx context.Context, c *models.Customer) error
	// FindByTenant returns all customers for a tenant.
	FindByTenant(ctx context.Context, tenantID uint) ([]models.Customer, error)
}
