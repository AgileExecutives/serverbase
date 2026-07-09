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
	"github.com/AgileExecutives/serverbase/pkg/repos"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactHandler struct {
	svc          *services.ContactService
	emailService *emailServices.EmailService
}

// NewContactHandler creates a new contact handler (legacy DB-based constructor)
func NewContactHandler(db *gorm.DB) *ContactHandler {
	return NewContactHandlerWithCtx(core.ModuleContext{DB: db})
}

// NewContactHandlerWithCtx creates a new contact handler using ModuleContext
func NewContactHandlerWithCtx(ctx core.ModuleContext) *ContactHandler {
	// Try to get ContactService from service registry
	if svcRaw, ok := ctx.Services.Get("contact"); ok {
		if svc, ok := svcRaw.(*services.ContactService); ok {
			return &ContactHandler{svc: svc, emailService: emailServices.NewEmailService()}
		}
	}
	// Fallback: construct service from repo factory
	rf := repos.NewGormRepoFactory(ctx.DB)
	contactRepo := rf.ContactRepo()
	svc := services.NewContactServiceWithRepo(contactRepo, ctx.Logger)
	return &ContactHandler{svc: svc, emailService: emailServices.NewEmailService()}
}

// GetContacts retrieves all contacts with pagination
// DISABLED-SWAGGER: @Summary Get all contacts
// DISABLED-SWAGGER: @Description Get a paginated list of all contacts
// DISABLED-SWAGGER: @Tags contacts
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param page query int false "Page number" default(1)
// DISABLED-SWAGGER: @Param limit query int false "Items per page" default(10)
// DISABLED-SWAGGER: @Param active query bool false "Filter by active status"
// DISABLED-SWAGGER: @Param type query string false "Filter by contact type"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.ListResponse}
// DISABLED-SWAGGER: @Failure 500 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contacts [get]
func (h *ContactHandler) GetContacts(c *gin.Context) {
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

// GetContact retrieves a specific contact by ID
// DISABLED-SWAGGER: @Summary Get contact by ID
// DISABLED-SWAGGER: @Description Get a specific contact by its ID
// DISABLED-SWAGGER: @Tags contacts
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Contact ID"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.ContactResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contacts/{id} [get]
func (h *ContactHandler) GetContact(c *gin.Context) {
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

// CreateContact creates a new contact
// DISABLED-SWAGGER: @Summary Create a new contact
// DISABLED-SWAGGER: @Description Create a new contact
// DISABLED-SWAGGER: @Tags contacts
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param request body models.ContactCreateRequest true "Contact creation data"
// DISABLED-SWAGGER: @Success 201 {object} models.APIResponse{data=models.ContactResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contacts [post]
func (h *ContactHandler) CreateContact(c *gin.Context) {
	var req models.ContactCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Set default type if not provided
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

// UpdateContact updates an existing contact
// DISABLED-SWAGGER: @Summary Update a contact
// DISABLED-SWAGGER: @Description Update an existing contact by ID
// DISABLED-SWAGGER: @Tags contacts
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Contact ID"
// DISABLED-SWAGGER: @Param request body models.ContactUpdateRequest true "Contact update data"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.ContactResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contacts/{id} [put]
func (h *ContactHandler) UpdateContact(c *gin.Context) {
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

	// Update fields if provided
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

// DeleteContact deletes a contact (soft delete)
// DISABLED-SWAGGER: @Summary Delete a contact
// DISABLED-SWAGGER: @Description Soft delete a contact by ID
// DISABLED-SWAGGER: @Tags contacts
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Contact ID"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contacts/{id} [delete]
func (h *ContactHandler) DeleteContact(c *gin.Context) {
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

// SubmitContactForm handles contact form submissions
// DISABLED-SWAGGER: @Summary Submit contact form
// DISABLED-SWAGGER: @Description Submit a contact form and optionally subscribe to newsletter
// DISABLED-SWAGGER: @Tags contact
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Param contactForm body models.ContactFormRequest true "Contact form data"
// DISABLED-SWAGGER: @Success 200 {object} models.ContactFormResponse
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 500 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contact/form [post]
func (h *ContactHandler) SubmitContactForm(c *gin.Context) {
	var req models.ContactFormRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Set timestamp if not provided
	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Send email to support
	err := h.emailService.SendContactFormEmail(
		req.Name,
		req.Email,
		req.Subject,
		req.Message,
		req.Timestamp,
		req.Source,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Internal server error", "Failed to send contact form email: "+err.Error()))
		return
	}

	response := models.ContactFormResponse{
		Message: "Contact form submitted successfully",
	}

	// Handle newsletter subscription if requested
	if req.Newsletter {
		newsletter := basemodels.Newsletter{Name: req.Name, Email: req.Email, Interest: req.Subject, Source: req.Source, LastContact: time.Now()}
		ok, nerr := h.svc.UpsertNewsletter(c.Request.Context(), &newsletter)
		if nerr != nil {
			response.NewsletterAdded = false
			response.NewsletterMessage = "Contact form sent, but newsletter subscription failed"
		} else if ok {
			response.NewsletterAdded = true
			response.NewsletterMessage = "Successfully subscribed to newsletter"
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetNewsletterSubscriptions gets all newsletter subscriptions (admin only)
// DISABLED-SWAGGER: @Summary Get newsletter subscriptions
// DISABLED-SWAGGER: @Description Get all newsletter subscriptions for admin users
// DISABLED-SWAGGER: @Tags contact
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Success 200 {array} models.Newsletter
// DISABLED-SWAGGER: @Failure 500 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Router /contact/newsletter [get]
func (h *ContactHandler) GetNewsletterSubscriptions(c *gin.Context) {
	list, err := h.svc.ListNewsletters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Internal server error", "Failed to fetch newsletter subscriptions"))
		return
	}
	c.JSON(http.StatusOK, list)
}

// UnsubscribeFromNewsletter handles newsletter unsubscription
// DISABLED-SWAGGER: @Summary Unsubscribe from newsletter
// DISABLED-SWAGGER: @Description Unsubscribe an email from the newsletter
// DISABLED-SWAGGER: @Tags contact
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Param email query string true "Email to unsubscribe"
// DISABLED-SWAGGER: @Success 200 {object} map[string]string
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 500 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /contact/newsletter/unsubscribe [delete]
func (h *ContactHandler) UnsubscribeFromNewsletter(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", "Email parameter is required"))
		return
	}
	rows, err := h.svc.DeleteNewsletterByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Internal server error", "Failed to unsubscribe from newsletter"))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Not found", "Email not found in newsletter subscriptions"))
		return
	}
	c.JSON(http.StatusOK, models.SuccessMessageResponse("Successfully unsubscribed from newsletter"))
}
