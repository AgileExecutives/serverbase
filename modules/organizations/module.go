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
	orgdocs "github.com/AgileExecutives/serverbase/internal/organizations/docs"
	orghandlers "github.com/AgileExecutives/serverbase/internal/organizations/handlers"
	orgrepo "github.com/AgileExecutives/serverbase/internal/organizations/repo"
	orgservices "github.com/AgileExecutives/serverbase/internal/organizations/services"
	"github.com/AgileExecutives/serverbase/module"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

// OrganizationsModule wires organization CRUD routes and swagger docs into the
// server module system.
// NewOrganizationsModule returns a new Organizations module implemented via AdapterModule.
func NewOrganizationsModule() core.Module {
	return module.NewAdapterModule("organizations", "1.0.0", []string{},
		module.WithInit(func(ctx core.ModuleContext) error {
			if ctx.DocRegistry != nil {
				ctx.DocRegistry.RegisterDoc("organizations", orgdocs.SwaggerInfoorganizations.ReadDoc())
			}
			return nil
		}),
		module.WithRoutes(&organizationsRouteProvider{}),
	)
}

// organizationsRouteProvider constructs services/handlers at RegisterRoutes time
// so they can use the ModuleContext (DB, Auth, Logger).
type organizationsRouteProvider struct{}

func (r *organizationsRouteProvider) GetPrefix() string                { return "" }
func (r *organizationsRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *organizationsRouteProvider) GetSwaggerTags() []string         { return []string{"organizations"} }

func (r *organizationsRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Construct service using serverbase org repo and wire to handlers
	repo := orgrepo.NewGormOrganizationRepo(ctx.DB)
	svc := orgservices.NewOrganizationServiceWithRepo(repo)
	h := orghandlers.NewOrganizationHandler(svc)

	organizations := router.Group("/organizations")
	organizations.Use(ctx.Auth.RequireAuth())
	{
		organizations.GET("/supported-formats", h.GetSupportedFormats)
		organizations.POST("", h.CreateOrganization)
		organizations.GET("", h.GetAllOrganizations)
		organizations.GET(":id", h.GetOrganization)
		organizations.PUT(":id", h.UpdateOrganization)
		organizations.DELETE(":id", h.DeleteOrganization)
	}
}
