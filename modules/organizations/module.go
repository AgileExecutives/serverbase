// Package organizations provides a core.Module that wires the internal
// organization CRUD handlers into the module system and registers their
// pre-generated swagger docs with the server's DocRegistry.
//
// This module lives in serverbase (not shared-modules) so it can legally
// import serverbase/internal packages.
//
// To regenerate swagger docs run from serverbase/internal/organizations:
//
//go:generate swag init -g doc.go --dir .,handlers,services,../../internal/models,../../pkg/utils --output docs --instanceName organizations
package organizations

import (
	"context"

	orgdocs "github.com/AgileExecutives/serverbase/internal/organizations/docs"
	orghandlers "github.com/AgileExecutives/serverbase/internal/organizations/handlers"
	orgrepo "github.com/AgileExecutives/serverbase/internal/organizations/repo"
	orgroutes "github.com/AgileExecutives/serverbase/internal/organizations/routes"
	orgservices "github.com/AgileExecutives/serverbase/internal/organizations/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
)

// OrganizationsModule wires organization CRUD routes and swagger docs into the
// server module system.
type OrganizationsModule struct {
	routeProvider *orgroutes.RouteProvider
}

// NewOrganizationsModule returns a new OrganizationsModule.
func NewOrganizationsModule() core.Module {
	return &OrganizationsModule{}
}

func (m *OrganizationsModule) Name() string           { return "organizations" }
func (m *OrganizationsModule) Version() string        { return "1.0.0" }
func (m *OrganizationsModule) Dependencies() []string { return []string{} }

func (m *OrganizationsModule) Initialize(ctx core.ModuleContext) error {
	// Use GORM-backed repo for persistence and pass it into the service.
	repo := orgrepo.NewGormOrganizationRepo(ctx.DB)
	svc := orgservices.NewOrganizationServiceWithRepo(repo)
	handler := orghandlers.NewOrganizationHandler(svc)
	m.routeProvider = orgroutes.NewRouteProvider(handler)

	if ctx.DocRegistry != nil {
		ctx.DocRegistry.RegisterDoc(m.Name(), orgdocs.SwaggerInfoorganizations.ReadDoc())
	}
	return nil
}

func (m *OrganizationsModule) Start(_ context.Context) error { return nil }
func (m *OrganizationsModule) Stop(_ context.Context) error  { return nil }

// Entities returns nothing — the Organization entity is owned by the minimal
// shared-modules/organization module which registers it separately.
func (m *OrganizationsModule) Entities() []core.Entity { return nil }

func (m *OrganizationsModule) Routes() []core.RouteProvider {
	if m.routeProvider == nil {
		return nil
	}
	return []core.RouteProvider{m.routeProvider}
}

func (m *OrganizationsModule) EventHandlers() []core.EventHandler    { return nil }
func (m *OrganizationsModule) Services() []core.ServiceProvider      { return nil }
func (m *OrganizationsModule) Middleware() []core.MiddlewareProvider { return nil }
func (m *OrganizationsModule) SwaggerPaths() []string                { return nil }
