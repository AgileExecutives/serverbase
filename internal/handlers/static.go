package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	statichandlers "github.com/AgileExecutives/shared-modules/static/handlers"
	"github.com/gin-gonic/gin"
)

// StaticHandler delegates static file serving to shared static handlers.
type StaticHandler struct {
	delegate *statichandlers.StaticHandlers
}

// NewStaticHandler creates a StaticHandler using defaults (fallback logger and FS repo)
func NewStaticHandler() *StaticHandler { return NewStaticHandlerWithCtx(core.ModuleContext{}) }

// NewStaticHandlerWithCtx creates a StaticHandler using ModuleContext (preferred)
func NewStaticHandlerWithCtx(ctx core.ModuleContext) *StaticHandler {
	// Use provided logger if available, otherwise fallback to StdLogger adapter
	var logger core.Logger = ctx.Logger
	if logger == nil {
		logger = core.NewLogger()
	}
	repo := statichandlers.NewFSStaticRepo("./statics/json")
	return &StaticHandler{delegate: statichandlers.NewStaticHandlers(logger, repo)}
}

// ServeStaticJSON delegates to the shared module handler
func (h *StaticHandler) ServeStaticJSON(c *gin.Context) { h.delegate.ServeStaticJSON(c) }

// ListStaticJSON delegates to the shared module handler
func (h *StaticHandler) ListStaticJSON(c *gin.Context) { h.delegate.ListStaticJSON(c) }

// Legacy exported functions kept for compatibility — use module-based handlers when possible
var legacyStatic = NewStaticHandler()

func ServeStaticJSON(c *gin.Context) { legacyStatic.ServeStaticJSON(c) }
func ListStaticJSON(c *gin.Context)  { legacyStatic.ListStaticJSON(c) }
