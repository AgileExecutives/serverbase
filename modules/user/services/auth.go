package services

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	tenantrepo "github.com/AgileExecutives/serverbase/modules/tenant/repo"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
)

// AuthService provides authentication related services
// TenantCreator defines the subset of TenantService required by AuthService.
type TenantCreator interface {
	CreateTenant(ctx context.Context, req models.TenantCreateRequest) (*models.Tenant, error)
}

type AuthService struct {
	userRepo       userrepo.UserRepo
	tenantRepo     tenantrepo.TenantRepo
	newsletterRepo userrepo.NewsletterRepo
	tokenRepo      userrepo.TokenBlacklistRepo
	logger         core.Logger
	tenantService  TenantCreator
}

// NewAuthServiceWithRepo creates an AuthService backed by module-local repos.
func NewAuthServiceWithRepo(userR userrepo.UserRepo, tenantR tenantrepo.TenantRepo, newsletterR userrepo.NewsletterRepo, tokenR userrepo.TokenBlacklistRepo, logger core.Logger) *AuthService {
	return &AuthService{userRepo: userR, tenantRepo: tenantR, newsletterRepo: newsletterR, tokenRepo: tokenR, logger: logger}
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

// SaveNewsletter persists a newsletter subscription via the user repo.
func (s *AuthService) SaveNewsletter(ctx context.Context, n *models.Newsletter) error {
	// Prefer explicit newsletter repo when available
	if s.newsletterRepo != nil {
		return s.newsletterRepo.SaveNewsletter(ctx, n)
	}
	// Fallback to underlying userRepo if it still implements the method
	if s.userRepo == nil {
		return nil
	}
	if nr, ok := s.userRepo.(interface {
		SaveNewsletter(ctx context.Context, n *models.Newsletter) error
	}); ok {
		return nr.SaveNewsletter(ctx, n)
	}
	return nil
}

// BlacklistToken persists a token blacklist entry via the user repo.
func (s *AuthService) BlacklistToken(ctx context.Context, tb *models.TokenBlacklist) error {
	if s.tokenRepo != nil {
		return s.tokenRepo.SaveTokenBlacklist(ctx, tb)
	}
	if s.userRepo == nil {
		return nil
	}
	if tr, ok := s.userRepo.(interface {
		SaveTokenBlacklist(ctx context.Context, tb *models.TokenBlacklist) error
	}); ok {
		return tr.SaveTokenBlacklist(ctx, tb)
	}
	return nil
}

// FindTenantByID looks up a tenant by numeric id using the injected TenantRepo.
func (s *AuthService) FindTenantByID(ctx context.Context, id uint) (*models.Tenant, error) {
	if s.tenantRepo == nil {
		return nil, nil
	}
	return s.tenantRepo.FindByID(ctx, id)
}

// FindTenantByName performs a name-based lookup (tenant names are not guaranteed unique,
// but used here to detect conflicts during signup flow).
func (s *AuthService) FindTenantByName(ctx context.Context, name string) (*models.Tenant, error) {
	if s.tenantRepo == nil {
		return nil, nil
	}
	tenants, err := s.tenantRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, t := range tenants {
		if t.Name == name {
			return &t, nil
		}
	}
	return nil, nil
}

// ListTenantSlugs returns all tenant slugs to help ensure unique slug generation.
func (s *AuthService) ListTenantSlugs(ctx context.Context) ([]string, error) {
	if s.tenantRepo == nil {
		return nil, nil
	}
	tenants, err := s.tenantRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	slugs := make([]string, 0, len(tenants))
	for _, t := range tenants {
		slugs = append(slugs, t.Slug)
	}
	return slugs, nil
}

// CreateTenant delegates tenant creation to the tenant repo.
func (s *AuthService) CreateTenant(ctx context.Context, t *models.Tenant) error {
	// If a higher-level TenantService is available, prefer it so tenant
	// creation logic (buckets, side-effects) is centralized.
	if s.tenantService != nil {
		req := models.TenantCreateRequest{CustomerID: t.CustomerID, Name: t.Name, Slug: t.Slug}
		created, err := s.tenantService.CreateTenant(ctx, req)
		if err != nil {
			return err
		}
		// propagate created ID back to the provided tenant pointer
		t.ID = created.ID
		return nil
	}
	if s.tenantRepo == nil {
		return nil
	}
	return s.tenantRepo.Save(ctx, t)
}

// SetTenantService allows wiring an internal TenantService so AuthService
// delegates tenant creation to the centralized implementation.
func (s *AuthService) SetTenantService(ts TenantCreator) {
	s.tenantService = ts
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
