package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	baseServices "github.com/AgileExecutives/serverbase/modules/base/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandlers provides basic health endpoints
type HealthHandlers struct {
	db     *gorm.DB
	logger core.Logger
}

func NewHealthHandlers(ctx core.ModuleContext, logger core.Logger) *HealthHandlers {
	return &HealthHandlers{db: ctx.DB, logger: logger}
}

// getEnvAsBool gets environment variable as boolean with fallback to default value
func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultVal
	}
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func (h *HealthHandlers) HealthCheck(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "2.0",
		"database":  "connected",
	}

	// Check database connectivity
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			response["database"] = "disconnected"
			response["status"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}
	}

	// Add environment configuration for test script
	environment := map[string]interface{}{
		"mock_email":         getEnvAsBool("MOCK_EMAIL", false),
		"rate_limit_enabled": getEnvAsBool("RATE_LIMIT_ENABLED", true),
		"email_verification": getEnvAsBool("FEATURE_EMAIL_VERIFICATION", true),
		"gin_mode":           os.Getenv("GIN_MODE"),
	}

	response["environment"] = environment

	c.JSON(http.StatusOK, response)
}
func (h *HealthHandlers) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse("pong", nil))
}

// UserSettingsHandlers minimal stub for user settings endpoints
type UserSettingsHandlers struct {
	svc    *baseServices.UserSettingsService
	logger core.Logger
}

func NewUserSettingsHandlers(ctx core.ModuleContext, logger core.Logger) *UserSettingsHandlers {
	return &UserSettingsHandlers{svc: baseServices.NewUserSettingsService(ctx.DB), logger: logger}
}

func (h *UserSettingsHandlers) GetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	// Delegate to service
	settings, err := h.svc.GetOrCreate(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings retrieved successfully", settings))
}
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

	// Delegate update to service
	updated, err := h.svc.Update(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings updated successfully", updated))
}
func (h *UserSettingsHandlers) ResetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	// Delegate reset to service
	reset, err := h.svc.Reset(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to reset settings", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("User settings reset successfully", reset))
}

// CheckResetToken method on AuthHandlers (minimal)
func (h *AuthHandlers) CheckResetToken(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse("OK", nil))
}
