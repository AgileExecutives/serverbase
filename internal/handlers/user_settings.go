package handlers

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/models"
	baseRepo "github.com/AgileExecutives/serverbase/modules/base/repo"
	baseServices "github.com/AgileExecutives/serverbase/modules/base/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSettingsHandler struct {
	svc *baseServices.UserSettingsService
}

// NewUserSettingsHandler is the legacy DB-based constructor
func NewUserSettingsHandler(db *gorm.DB) *UserSettingsHandler {
	return NewUserSettingsHandlerWithCtx(core.ModuleContext{DB: db})
}

// NewUserSettingsHandlerWithCtx creates a UserSettingsHandler using ModuleContext
func NewUserSettingsHandlerWithCtx(ctx core.ModuleContext) *UserSettingsHandler {
	if svcRaw, ok := ctx.Services.Get("base-user-settings"); ok {
		if svc, ok := svcRaw.(*baseServices.UserSettingsService); ok {
			return &UserSettingsHandler{svc: svc}
		}
	}
	// Fallback: construct directly from DB-backed repo
	return &UserSettingsHandler{svc: baseServices.NewUserSettingsService(baseRepo.NewGormUserSettingsRepo(ctx.DB))}
}

// GetUserSettings retrieves the current user's settings
func (h *UserSettingsHandler) GetUserSettings(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	settings, err := h.svc.GetOrCreate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve user settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings retrieved successfully", settings.ToResponse()))
}

// UpdateUserSettings updates the current user's settings
func (h *UserSettingsHandler) UpdateUserSettings(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	var req models.UserSettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	updated, err := h.svc.Update(user.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update user settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings updated successfully", updated.ToResponse()))
}

// ResetUserSettings resets the current user's settings to defaults
func (h *UserSettingsHandler) ResetUserSettings(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	settings, err := h.svc.Reset(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to reset user settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings reset to defaults", settings.ToResponse()))
}
