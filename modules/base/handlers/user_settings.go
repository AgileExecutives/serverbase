package handlers

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/models"
	baseServices "github.com/AgileExecutives/serverbase/modules/base/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

type UserSettingsHandlers struct {
	svc    *baseServices.UserSettingsService
	logger core.Logger
}

func NewUserSettingsHandlers(svc *baseServices.UserSettingsService, logger core.Logger) *UserSettingsHandlers {
	return &UserSettingsHandlers{
		svc:    svc,
		logger: logger,
	}
}

// @Summary Get user settings
// @ID getUserSettings
// @Description Get settings for the authenticated user (creates default if not found)
// @Tags user-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserSettings}
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user-settings [get]
func (h *UserSettingsHandlers) GetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	s, err := h.svc.GetOrCreate(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve settings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User settings retrieved successfully", s))
}

// @Summary Update user settings
// @ID updateUserSettings
// @Description Update settings for the authenticated user
// @Tags user-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body models.UserSettingsUpdateRequest true "Updated settings"
// @Success 200 {object} models.APIResponse{data=models.UserSettings}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user-settings [put]
func (h *UserSettingsHandlers) UpdateUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	var req models.UserSettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	updated, err := h.svc.Update(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update settings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User settings updated successfully", updated))
}

// @Summary Reset user settings
// @ID resetUserSettings
// @Description Reset user settings to default values
// @Tags user-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserSettings}
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user-settings/reset [post]
func (h *UserSettingsHandlers) ResetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	reset, err := h.svc.Reset(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to reset settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings reset successfully", reset))
}
