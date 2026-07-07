package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/AgileExecutives/serverbase/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *EmailHandler) GetLatestEmails(c *gin.Context) {
	mockEmail := os.Getenv("MOCK_EMAIL")
	if strings.ToLower(mockEmail) != "true" {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponseFunc("Mock email not enabled", "This endpoint is only available when MOCK_EMAIL=true"))
		return
	}

	emails, err := h.emailService.GetLatestEmails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve mock emails", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Mock emails retrieved successfully", emails))
}
