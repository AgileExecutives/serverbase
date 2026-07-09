package routes

import (
	"github.com/AgileExecutives/serverbase/internal/organizations/handlers"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

// RouteProvider provides routing functionality for organization management
type RouteProvider struct{ handler *handlers.OrganizationHandler }

// NewRouteProvider creates a new route provider
func NewRouteProvider(handler *handlers.OrganizationHandler) *RouteProvider {
	return &RouteProvider{handler: handler}
}

// RegisterRoutes registers the organization management routes
func (rp *RouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	organizations := router.Group("/organizations")
	organizations.Use(ctx.Auth.RequireAuth())
	{
		organizations.GET("/supported-formats", rp.handler.GetSupportedFormats)
		organizations.POST("", rp.handler.CreateOrganization)
		organizations.GET("", rp.handler.GetAllOrganizations)
		organizations.GET("/:id", rp.handler.GetOrganization)
		organizations.PUT("/:id", rp.handler.UpdateOrganization)
		organizations.DELETE("/:id", rp.handler.DeleteOrganization)
	}
}

// GetPrefix returns the route prefix for organization endpoints
func (rp *RouteProvider) GetPrefix() string { return "" }

// GetMiddleware returns middleware to apply to all routes
func (rp *RouteProvider) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }

// GetSwaggerTags returns swagger tags for the routes
func (rp *RouteProvider) GetSwaggerTags() []string { return []string{"organizations"} }
