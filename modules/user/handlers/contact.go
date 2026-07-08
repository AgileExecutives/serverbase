package handlers

import (
	"net/http"
	"time"

	emailServices "github.com/AgileExecutives/serverbase/modules/email/services"
	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContactHandlers struct {
	db           *gorm.DB
	logger       core.Logger
	emailService *emailServices.EmailService
}

func NewContactHandlers(db *gorm.DB, logger core.Logger) *ContactHandlers {
	return &ContactHandlers{
		db:           db,
		logger:       logger,
		emailService: emailServices.NewEmailService(),
	}
}

func (h *ContactHandlers) GetContacts(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := utils.GetOffset(page, limit)
	var contacts []models.Contact
	var total int64
	query := h.db.Model(&models.Contact{})
	if activeStr := c.Query("active"); activeStr != "" {
		if activeStr == "true" {
			query = query.Where("active = ?", true)
		} else if activeStr == "false" {
			query = query.Where("active = ?", false)
		}
	}
	if contactType := c.Query("type"); contactType != "" {
		query = query.Where("type = ?", contactType)
	}
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to count contacts", err.Error()))
		return
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve contacts", err.Error()))
		return
	}
	var responses []models.ContactResponse
	for _, contact := range contacts {
		responses = append(responses, contact.ToResponse())
	}
	response := models.ListResponse{
		Data: responses,
		Pagination: models.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: utils.CalculateTotalPages(int(total), limit),
		},
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Contacts retrieved successfully", response))
}

func (h *ContactHandlers) GetContact(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid contact ID", err.Error()))
		return
	}
	var contact models.Contact
	if err := h.db.First(&contact, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Contact not found", "Contact with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve contact", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Contact retrieved successfully", contact.ToResponse()))
}

func (h *ContactHandlers) CreateContact(c *gin.Context) {
	var req models.ContactCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}
	contactType := req.Type
	if contactType == "" {
		contactType = "contact"
	}
	contact := models.Contact{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Mobile:    req.Mobile,
		Street:    req.Street,
		Zip:       req.Zip,
		City:      req.City,
		Country:   req.Country,
		Type:      contactType,
		Notes:     req.Notes,
		Active:    true,
	}
	if err := h.db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create contact", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.SuccessResponse("Contact created successfully", contact.ToResponse()))
}

func (h *ContactHandlers) UpdateContact(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid contact ID", err.Error()))
		return
	}
	var req models.ContactUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}
	var contact models.Contact
	if err := h.db.First(&contact, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Contact not found", "Contact with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve contact", err.Error()))
		return
	}
	if req.FirstName != "" {
		contact.FirstName = req.FirstName
	}
	if req.LastName != "" {
		contact.LastName = req.LastName
	}
	if req.Email != "" {
		contact.Email = req.Email
	}
	if req.Phone != "" {
		contact.Phone = req.Phone
	}
	if req.Mobile != "" {
		contact.Mobile = req.Mobile
	}
	if req.Street != "" {
		contact.Street = req.Street
	}
	if req.Zip != "" {
		contact.Zip = req.Zip
	}
	if req.City != "" {
		contact.City = req.City
	}
	if req.Country != "" {
		contact.Country = req.Country
	}
	if req.Type != "" {
		contact.Type = req.Type
	}
	if req.Notes != "" {
		contact.Notes = req.Notes
	}
	if req.Active != nil {
		contact.Active = *req.Active
	}
	if err := h.db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update contact", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Contact updated successfully", contact.ToResponse()))
}

func (h *ContactHandlers) DeleteContact(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid contact ID", err.Error()))
		return
	}
	var contact models.Contact
	if err := h.db.First(&contact, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Contact not found", "Contact with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve contact", err.Error()))
		return
	}
	if err := h.db.Delete(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to delete contact", err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Contact deleted successfully", nil))
}

func (h *ContactHandlers) SubmitContactForm(c *gin.Context) {
	var req models.ContactFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format(time.RFC3339)
	}
	err := h.emailService.SendContactFormEmail(req.Name, req.Email, req.Subject, req.Message, req.Timestamp, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send contact form email: " + err.Error()})
		return
	}
	response := models.ContactFormResponse{Message: "Contact form submitted successfully"}
	if req.Newsletter {
		newsletter := basemodels.Newsletter{Name: req.Name, Email: req.Email, Interest: req.Subject, Source: req.Source, LastContact: time.Now()}
		res := h.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "email"}}, DoUpdates: clause.AssignmentColumns([]string{"name", "interest", "source", "last_contact"})}).Create(&newsletter)
		if res.Error != nil {
			h.logger.Warn("newsletter upsert failed", res.Error)
			var existing basemodels.Newsletter
			if r := h.db.Where("email = ?", req.Email).First(&existing); r.Error == nil {
				existing.Name = req.Name
				existing.Interest = req.Subject
				existing.Source = req.Source
				existing.LastContact = time.Now()
				if err2 := h.db.Save(&existing).Error; err2 == nil {
					response.NewsletterAdded = true
					response.NewsletterMessage = "Newsletter subscription updated"
					h.logger.Info("newsletter updated via fallback", "email", existing.Email)
				} else {
					response.NewsletterAdded = false
					response.NewsletterMessage = "Contact form sent, but newsletter update failed"
					h.logger.Error("newsletter fallback update failed", err2)
				}
			} else {
				response.NewsletterAdded = false
				response.NewsletterMessage = "Contact form sent, but newsletter subscription failed"
				h.logger.Warn("newsletter fallback: not found", "email", req.Email, "err", r.Error)
			}
		} else {
			response.NewsletterAdded = true
			response.NewsletterMessage = "Successfully subscribed to newsletter"
			h.logger.Info("newsletter upsert succeeded", "email", newsletter.Email, "rows", res.RowsAffected)
		}
	}
	c.JSON(http.StatusOK, response)
}

func (h *ContactHandlers) GetNewsletterSubscriptions(c *gin.Context) {
	var newsletters []basemodels.Newsletter
	if err := h.db.Find(&newsletters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch newsletter subscriptions"})
		return
	}
	c.JSON(http.StatusOK, newsletters)
}

func (h *ContactHandlers) UnsubscribeFromNewsletter(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}
	// Unscoped so the row is hard-deleted; soft-deleted rows are not re-subscribable
	// via the OnConflict upsert path (no unique index on email in SQLite).
	result := h.db.Unscoped().Where("email = ?", email).Delete(&basemodels.Newsletter{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe from newsletter"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in newsletter subscriptions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed from newsletter"})
}
