package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	saas_handlers "github.com/AgileExecutives/shared-modules/saas-base/handlers"
	saasrepo "github.com/AgileExecutives/shared-modules/saas-base/repo"
	saas_services "github.com/AgileExecutives/shared-modules/saas-base/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlanHandler delegates to shared saas-base plan handlers (uses service layer)
type PlanHandler struct{ delegate *saas_handlers.PlanHandlers }

// NewPlanHandler is the legacy DB-based constructor
func NewPlanHandler(db *gorm.DB) *PlanHandler {
	return NewPlanHandlerWithCtx(core.ModuleContext{DB: db})
}

// NewPlanHandlerWithCtx creates a PlanHandler using ModuleContext
func NewPlanHandlerWithCtx(ctx core.ModuleContext) *PlanHandler {
	// Try to lookup registered saas plan service
	if svcRaw, ok := ctx.Services.Get("saas-base-plan"); ok {
		if svc, ok := svcRaw.(*saas_services.PlanService); ok {
			return &PlanHandler{delegate: saas_handlers.NewPlanHandlers(svc)}
		}
	}
	// Fallback: create plan service from saas repo
	planRepo := saasrepo.NewGormPlanRepo(ctx.DB)
	svc := saas_services.NewPlanService(planRepo)
	return &PlanHandler{delegate: saas_handlers.NewPlanHandlers(svc)}
}

func (h *PlanHandler) GetPlans(c *gin.Context)   { h.delegate.GetPlans(c) }
func (h *PlanHandler) GetPlan(c *gin.Context)    { h.delegate.GetPlan(c) }
func (h *PlanHandler) CreatePlan(c *gin.Context) { h.delegate.CreatePlan(c) }
func (h *PlanHandler) UpdatePlan(c *gin.Context) { h.delegate.UpdatePlan(c) }
func (h *PlanHandler) DeletePlan(c *gin.Context) { h.delegate.DeletePlan(c) }
