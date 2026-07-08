// Package saas provides a core.Module that registers customer and plan routes
// from shared-modules/saas-base.  Newsletter is intentionally excluded because
// the user module already handles newsletter subscriptions via /contacts/newsletter.
//
// Using this module avoids the saas-base CoreModule's AutoMigrate which would
// add a uniqueIndex on newsletters.email, conflicting with the user module's
// internal Newsletter model that intentionally uses a non-unique index.
package saas

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
	saasEntities "github.com/AgileExecutives/shared-modules/saas-base/entities"
	saasHandlers "github.com/AgileExecutives/shared-modules/saas-base/handlers"
	"github.com/gin-gonic/gin"
)

// NewSaaSModule returns a new SaaSModule.
func NewSaaSModule() core.Module {
	return &saaSModule{}
}

type saaSModule struct{}

func (m *saaSModule) Name() string           { return "saas" }
func (m *saaSModule) Version() string        { return "1.0.0" }
func (m *saaSModule) Dependencies() []string { return []string{} }

func (m *saaSModule) Initialize(_ core.ModuleContext) error { return nil }
func (m *saaSModule) Start(_ context.Context) error         { return nil }
func (m *saaSModule) Stop(_ context.Context) error          { return nil }

// Entities returns only Plan and Customer — no Newsletter — so we don't add a
// unique constraint on newsletter.email that conflicts with the user module.
func (m *saaSModule) Entities() []core.Entity {
	return []core.Entity{
		saasEntities.NewPlanEntity(),
		saasEntities.NewCustomerEntity(),
	}
}

func (m *saaSModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{
		&customerRouteProvider{},
		&planRouteProvider{},
	}
}

func (m *saaSModule) EventHandlers() []core.EventHandler    { return nil }
func (m *saaSModule) Services() []core.ServiceProvider      { return nil }
func (m *saaSModule) Middleware() []core.MiddlewareProvider { return nil }
func (m *saaSModule) SwaggerPaths() []string                { return nil }

// ─── Route providers ─────────────────────────────────────────────────────────

type customerRouteProvider struct{}

func (r *customerRouteProvider) GetPrefix() string                { return "/customers" }
func (r *customerRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *customerRouteProvider) GetSwaggerTags() []string         { return []string{"customers"} }
func (r *customerRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	h := saasHandlers.NewCustomerHandlers(ctx.DB)
	auth := router.Group("")
	auth.Use(ctx.Auth.RequireAuth())
	auth.GET("", h.GetCustomers)
	auth.POST("", h.CreateCustomer)
	auth.GET("/:id", h.GetCustomer)
	auth.PUT("/:id", h.UpdateCustomer)
	auth.DELETE("/:id", h.DeleteCustomer)
}

type planRouteProvider struct{}

func (r *planRouteProvider) GetPrefix() string                { return "/plans" }
func (r *planRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *planRouteProvider) GetSwaggerTags() []string         { return []string{"plans"} }
func (r *planRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	h := saasHandlers.NewPlanHandlers(ctx.DB)
	// Public read endpoints
	router.GET("", h.GetPlans)
	router.GET("/:id", h.GetPlan)
	// Admin write endpoints
	admin := router.Group("")
	admin.Use(ctx.Auth.RequireAuth(), ctx.Auth.RequireRole("admin", "super-admin"))
	admin.POST("", h.CreatePlan)
	admin.PUT("/:id", h.UpdatePlan)
	admin.DELETE("/:id", h.DeletePlan)
}
