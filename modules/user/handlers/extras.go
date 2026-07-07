package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandlers provides basic health endpoints
type HealthHandlers struct {
	db     *gorm.DB
	logger core.Logger
}

func NewHealthHandlers(db *gorm.DB, logger core.Logger) *HealthHandlers {
	return &HealthHandlers{db: db, logger: logger}
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
	db     *gorm.DB
	logger core.Logger
}

func NewUserSettingsHandlers(db *gorm.DB, logger core.Logger) *UserSettingsHandlers {
	return &UserSettingsHandlers{db: db, logger: logger}
}

func (h *UserSettingsHandlers) GetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	var settings models.UserSettings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.UserSettings{
				UserID:   userID.(uint),
				Language: "en",
				Timezone: "UTC",
				Theme:    "light",
				Settings: "{}",
			}
			if err := h.db.Create(&settings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create default settings", err.Error()))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve settings", err.Error()))
			return
		}
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

	var settings models.UserSettings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.UserSettings{UserID: userID.(uint)}
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve settings", err.Error()))
			return
		}
	}

	if req.Language != "" {
		settings.Language = req.Language
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}
	if req.Theme != "" {
		settings.Theme = req.Theme
	}
	if req.Settings != "" {
		settings.Settings = req.Settings
	}

	if err := h.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update settings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User settings updated successfully", settings))
}
func (h *UserSettingsHandlers) ResetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Unauthorized", "User ID not found in context"))
		return
	}

	var settings models.UserSettings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Settings not found", "No settings found for user"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve settings", err.Error()))
		return
	}

	settings.Language = "en"
	settings.Timezone = "UTC"
	settings.Theme = "light"
	settings.Settings = "{}"

	if err := h.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to reset settings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User settings reset successfully", settings))
}

// CheckResetToken method on AuthHandlers (minimal)
func (h *AuthHandlers) CheckResetToken(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse("OK", nil))
}
