package handlers

import (
	"errors"
	"net/http"
	"time"

	emailServices "github.com/AgileExecutives/serverbase/modules/email/services"
	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactHandlers struct {
	svc          *services.ContactService
	logger       core.Logger
	emailService *emailServices.EmailService
}

func NewContactHandlers(ctx core.ModuleContext, svc *services.ContactService, logger core.Logger) *ContactHandlers {
	return &ContactHandlers{
		svc:          svc,
		logger:       logger,
		emailService: emailServices.NewEmailService(),
	}
}

func (h *ContactHandlers) GetContacts(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := utils.GetOffset(page, limit)
	var activePtr *bool
	if activeStr := c.Query("active"); activeStr != "" {
		v := activeStr == "true"
		activePtr = &v
	}
	contactType := c.Query("type")
	contacts, total, err := h.svc.ListContacts(c.Request.Context(), offset, limit, activePtr, contactType)
	if err != nil {
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
	contact, err := h.svc.GetContact(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err := h.svc.CreateContact(c.Request.Context(), &contact); err != nil {
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
	contact, err := h.svc.GetContact(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err := h.svc.UpdateContact(c.Request.Context(), contact); err != nil {
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
	contact, err := h.svc.GetContact(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Contact not found", "Contact with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve contact", err.Error()))
		return
	}
	if err := h.svc.DeleteContact(c.Request.Context(), contact); err != nil {
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
		ok, nerr := h.svc.UpsertNewsletter(c.Request.Context(), &newsletter)
		if nerr != nil {
			h.logger.Warn("newsletter upsert failed", nerr)
			response.NewsletterAdded = false
			response.NewsletterMessage = "Contact form sent, but newsletter subscription failed"
		} else if ok {
			response.NewsletterAdded = true
			response.NewsletterMessage = "Successfully subscribed to newsletter"
			h.logger.Info("newsletter upsert succeeded", "email", newsletter.Email)
		}
	}
	c.JSON(http.StatusOK, response)
}

func (h *ContactHandlers) GetNewsletterSubscriptions(c *gin.Context) {
	list, err := h.svc.ListNewsletters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch newsletter subscriptions"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ContactHandlers) UnsubscribeFromNewsletter(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}
	// Unscoped so the row is hard-deleted; soft-deleted rows are not re-subscribable
	// via the OnConflict upsert path (no unique index on email in SQLite).
	rows, err := h.svc.DeleteNewsletterByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe from newsletter"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in newsletter subscriptions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed from newsletter"})
}
