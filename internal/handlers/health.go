package handlers

import (
	basehandlers "github.com/AgileExecutives/serverbase/modules/base/handlers"
	"github.com/AgileExecutives/serverbase/pkg/config"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler delegates health endpoints to the base module's handlers.
type HealthHandler struct{ delegate *basehandlers.HealthHandlers }

// NewHealthHandler creates a new health handler (legacy DB-based constructor)
func NewHealthHandler(db *gorm.DB, cfg config.Config) *HealthHandler {
	return NewHealthHandlerWithCtx(core.ModuleContext{DB: db}, cfg)
}

// NewHealthHandlerWithCtx creates a HealthHandler using ModuleContext
func NewHealthHandlerWithCtx(ctx core.ModuleContext, cfg config.Config) *HealthHandler {
	// Delegate to base module's health handlers
	return &HealthHandler{delegate: basehandlers.NewHealthHandlers(ctx, ctx.Logger)}
}

func (h *HealthHandler) Health(c *gin.Context) { h.delegate.HealthCheck(c) }
func (h *HealthHandler) Ping(c *gin.Context)   { h.delegate.Ping(c) }
