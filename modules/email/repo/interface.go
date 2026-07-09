package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// EmailRepo defines persistence operations for emails
type EmailRepo interface {
	List(ctx context.Context, offset, limit int, status string) ([]models.Email, int64, error)
	FindByID(ctx context.Context, id uint) (*models.Email, error)
	Create(ctx context.Context, e *models.Email) error
	UpdateStatus(ctx context.Context, id uint, status, errorMessage string) error
	Stats(ctx context.Context) (map[string]int64, error)
}
