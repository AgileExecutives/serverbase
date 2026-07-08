package services

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
)

// AuthService provides authentication related services
type AuthService struct {
	userRepo repo.UserRepo
	logger   core.Logger
}

// NewAuthServiceWithRepo creates an AuthService backed by a module-local UserRepo.
func NewAuthServiceWithRepo(r repo.UserRepo, logger core.Logger) *AuthService {
	return &AuthService{userRepo: r, logger: logger}
}

// Backwards compatible constructor - kept for compatibility (wraps repo not provided).
func NewAuthService(db interface{}, logger core.Logger) *AuthService {
	// If a DB was provided, a repository adapter will be created by the module.
	return &AuthService{userRepo: nil, logger: logger}
}

// FindByEmail looks up a user by email.
func (s *AuthService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if s.userRepo == nil {
		return nil, nil
	}
	return s.userRepo.FindByEmail(ctx, email)
}

// SaveUser saves the user entity via repository.
func (s *AuthService) SaveUser(ctx context.Context, u *models.User) error {
	if s.userRepo == nil {
		return nil
	}
	return s.userRepo.Save(ctx, u)
}

type AuthServiceProvider struct{ service *AuthService }

func NewAuthServiceProvider(service *AuthService) core.ServiceProvider {
	return &AuthServiceProvider{service: service}
}

func (p *AuthServiceProvider) ServiceName() string           { return "auth" }
func (p *AuthServiceProvider) ServiceInterface() interface{} { return (*AuthService)(nil) }
func (p *AuthServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	return p.service, nil
}
