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

	"github.com/AgileExecutives/serverbase/module"
	custrepo "github.com/AgileExecutives/serverbase/modules/customers/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
	saasEntities "github.com/AgileExecutives/shared-modules/saas-base/entities"
	saasHandlers "github.com/AgileExecutives/shared-modules/saas-base/handlers"
	saasrepo "github.com/AgileExecutives/shared-modules/saas-base/repo"
	saassvc "github.com/AgileExecutives/shared-modules/saas-base/services"
	"github.com/gin-gonic/gin"
)

// NewSaaSModule returns a new SaaSModule implemented via AdapterModule.
func NewSaaSModule() core.Module {
	return module.NewAdapterModule("saas", "1.0.0", []string{},
		module.WithEntities(saasEntities.NewPlanEntity(), saasEntities.NewCustomerEntity()),
		module.WithRoutes(&customerRouteProvider{}, &planRouteProvider{}),
	)
}

type saaSModule struct{}

func (m *saaSModule) Initialize(_ core.ModuleContext) error { return nil }
func (m *saaSModule) Start(_ context.Context) error         { return nil }
func (m *saaSModule) Stop(_ context.Context) error          { return nil }

// Note: the route provider implementations remain unchanged and construct
// services/handlers at RegisterRoutes time so they work with the adapter.

// ─── Route providers ─────────────────────────────────────────────────────────

type customerRouteProvider struct{}

func (r *customerRouteProvider) GetPrefix() string                { return "/customers" }
func (r *customerRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *customerRouteProvider) GetSwaggerTags() []string         { return []string{"customers"} }
func (r *customerRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Construct customer service using serverbase customer repo and wire to handlers
	custRepo := custrepo.NewGormCustomerRepo(ctx.DB)
	custSvc := saassvc.NewCustomerServiceWithDB(custRepo, ctx.DB, ctx.Logger)

	h := saasHandlers.NewCustomerHandlers(custSvc)
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
	// Construct plan service using saas-base plan repo
	planRepo := saasrepo.NewGormPlanRepo(ctx.DB)
	planSvc := saassvc.NewPlanService(planRepo)
	h := saasHandlers.NewPlanHandlers(planSvc)
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
