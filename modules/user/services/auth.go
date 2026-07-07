package services

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"gorm.io/gorm"
)

// AuthService provides authentication related services
type AuthService struct {
	db     *gorm.DB
	logger core.Logger
}

func NewAuthService(db *gorm.DB, logger core.Logger) *AuthService {
	return &AuthService{db: db, logger: logger}
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
