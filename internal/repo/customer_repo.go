package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// CustomerRepo defines persistence operations for customers used by internal handlers.
type CustomerRepo interface {
	ListByTenant(ctx context.Context, tenantID uint, offset, limit int, active *bool) ([]models.Customer, int64, error)
	GetByID(ctx context.Context, id, tenantID uint) (*models.Customer, error)
	Create(ctx context.Context, c *models.Customer) error
	Update(ctx context.Context, c *models.Customer) error
	Delete(ctx context.Context, c *models.Customer) error
	PlanExists(ctx context.Context, planID uint) (bool, error)
}
