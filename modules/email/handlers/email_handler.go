package handlers

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/modules/email/services"
	"github.com/AgileExecutives/serverbase/pkg/config"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/middleware"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailHandler struct {
	db           *gorm.DB
	emailService *services.EmailService
}

func NewEmailHandler(db *gorm.DB, emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{db: db, emailService: emailService}
}

func (h *EmailHandler) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	emailRoutes := router.Group("/emails")
	{
		emailRoutes.GET("", h.GetEmails)
		emailRoutes.GET(":id", h.GetEmail)
		emailRoutes.POST("/send", h.SendEmail)
		emailRoutes.GET("/stats", h.GetEmailStats)
	}
}

func (h *EmailHandler) GetPrefix() string { return "" }
func (h *EmailHandler) GetMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{middleware.AuthMiddleware(h.db)}
}
func (h *EmailHandler) GetSwaggerTags() []string { return []string{"emails"} }

func (h *EmailHandler) GetEmails(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := utils.GetOffset(page, limit)

	var emails []models.Email
	var total int64

	query := h.db.Model(&models.Email{})
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to count emails", err.Error()))
		return
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&emails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve emails", err.Error()))
		return
	}

	var responses []models.EmailResponse
	for _, email := range emails {
		responses = append(responses, email.ToResponse())
	}

	response := models.ListResponse{Data: responses, Pagination: models.PaginationResponse{Page: page, Limit: limit, Total: int(total), TotalPages: utils.CalculateTotalPages(int(total), limit)}}
	c.JSON(http.StatusOK, models.SuccessResponse("Emails retrieved successfully", response))
}

func (h *EmailHandler) GetEmail(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid email ID", err.Error()))
		return
	}

	var email models.Email
	if err := h.db.First(&email, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Email not found", "Email with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve email", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Email retrieved successfully", email.ToResponse()))
}

func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req models.EmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try alternative payload shape used by hurl tests (to_email / from_email)
		var alt struct {
			ToEmail   string `json:"to_email"`
			FromEmail string `json:"from_email"`
			Subject   string `json:"subject"`
			Body      string `json:"body"`
			HTMLBody  string `json:"html_body"`
		}
		if err2 := c.ShouldBindJSON(&alt); err2 != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
			return
		}
		req.To = alt.ToEmail
		req.From = alt.FromEmail
		req.Subject = alt.Subject
		req.Body = alt.Body
		req.HTMLBody = alt.HTMLBody
	}

	// Fill defaults if optional fields missing
	if req.From == "" {
		// Use configured default from address
		cfg := config.Load()
		req.From = cfg.Email.FromEmail
	}

	email := models.Email{To: req.To, From: req.From, Subject: req.Subject, Body: req.Body, HTMLBody: req.HTMLBody, Status: "pending"}
	if err := h.db.Create(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create email", err.Error()))
		return
	}

	go func() {
		textBody := req.Body
		if textBody == "" {
			textBody = req.Subject
		}
		err := h.emailService.SendEmail(req.To, req.Subject, req.HTMLBody, textBody)
		status := "sent"
		var errorMessage string
		if err != nil {
			status = "failed"
			errorMessage = err.Error()
		}
		h.db.Model(&email).Updates(models.Email{Status: status, ErrorMessage: errorMessage})
	}()

	c.JSON(http.StatusCreated, models.SuccessResponse("Email queued successfully", email.ToResponse()))
}

func (h *EmailHandler) GetEmailStats(c *gin.Context) {
	type EmailStats struct {
		Total     int64 `json:"total"`
		Pending   int64 `json:"pending"`
		Sent      int64 `json:"sent"`
		Delivered int64 `json:"delivered"`
		Failed    int64 `json:"failed"`
	}
	var stats EmailStats
	h.db.Model(&models.Email{}).Count(&stats.Total)
	h.db.Model(&models.Email{}).Where("status = ?", "pending").Count(&stats.Pending)
	h.db.Model(&models.Email{}).Where("status = ?", "sent").Count(&stats.Sent)
	h.db.Model(&models.Email{}).Where("status = ?", "delivered").Count(&stats.Delivered)
	h.db.Model(&models.Email{}).Where("status = ?", "failed").Count(&stats.Failed)
	c.JSON(http.StatusOK, models.SuccessResponse("Email statistics retrieved successfully", stats))
}
