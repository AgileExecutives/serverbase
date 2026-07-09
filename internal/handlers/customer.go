package handlers

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/internal/repo"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	repo repo.CustomerRepo
}

// NewCustomerHandler creates a new customer handler using the provided repo
func NewCustomerHandler(r repo.CustomerRepo) *CustomerHandler {
	return &CustomerHandler{repo: r}
}

// GetCustomers retrieves all customers with pagination and tenant isolation
// DISABLED-SWAGGER: @Summary Get all customers
// DISABLED-SWAGGER: @Description Get a paginated list of customers for the authenticated tenant
// DISABLED-SWAGGER: @Tags customers
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param page query int false "Page number" default(1)
// DISABLED-SWAGGER: @Param limit query int false "Items per page" default(10)
// DISABLED-SWAGGER: @Param active query bool false "Filter by active status"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.ListResponse}
// DISABLED-SWAGGER: @Failure 500 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /customers [get]
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	// Get user from context for tenant isolation
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	page, limit := utils.GetPaginationParams(c)
	offset := utils.GetOffset(page, limit)

	var customers []models.Customer
	var total int64

	ctx := c.Request.Context()
	// use repo to fetch customers
	var activeFilter *bool
	if activeStr := c.Query("active"); activeStr != "" {
		if activeStr == "true" {
			t := true
			activeFilter = &t
		} else if activeStr == "false" {
			f := false
			activeFilter = &f
		}
	}

	customers, total, err := h.repo.ListByTenant(ctx, user.TenantID, offset, limit, activeFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve customers", err.Error()))
		return
	}

	// active filtering handled by repo

	// Convert to response format
	var responses []models.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, customer.ToResponse())
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

	c.JSON(http.StatusOK, models.SuccessResponse("Customers retrieved successfully", response))
}

// GetCustomer retrieves a specific customer by ID with tenant isolation
// DISABLED-SWAGGER: @Summary Get customer by ID
// DISABLED-SWAGGER: @Description Get a specific customer by its ID within the authenticated tenant
// DISABLED-SWAGGER: @Tags customers
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Customer ID"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	// Get user from context for tenant isolation
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	ctx := c.Request.Context()
	customer, err := h.repo.GetByID(ctx, uint(id), user.TenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Customer retrieved successfully", customer.ToResponse()))
}

// CreateCustomer creates a new customer
// DISABLED-SWAGGER: @Summary Create a new customer
// DISABLED-SWAGGER: @Description Create a new customer within the authenticated tenant
// DISABLED-SWAGGER: @Tags customers
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param request body models.CustomerCreateRequest true "Customer creation data"
// DISABLED-SWAGGER: @Success 201 {object} models.APIResponse{data=models.CustomerResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	// Get user from context for tenant isolation
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	var req models.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Ensure customer is created within user's organization
	req.TenantID = user.TenantID

	// Verify the plan exists
	ctx := c.Request.Context()
	ok, err := h.repo.PlanExists(ctx, req.PlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to verify plan", err.Error()))
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Plan not found", "Invalid plan ID"))
		return
	}

	customer := models.Customer{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Street:        req.Street,
		Zip:           req.Zip,
		City:          req.City,
		Country:       req.Country,
		TaxID:         req.TaxID,
		VAT:           req.VAT,
		PlanID:        req.PlanID,
		TenantID:      req.TenantID,
		Status:        "active",
		PaymentMethod: req.PaymentMethod,
		Active:        true,
	}

	if err := h.repo.Create(c.Request.Context(), &customer); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create customer", err.Error()))
		return
	}

	// Note: Plan and Tenant relations temporarily disabled due to GORM relation issues
	// h.db.Preload("Plan").Preload("Tenant").First(&customer, customer.ID)

	c.JSON(http.StatusCreated, models.SuccessResponse("Customer created successfully", customer.ToResponse()))
}

// UpdateCustomer updates an existing customer
// DISABLED-SWAGGER: @Summary Update a customer
// DISABLED-SWAGGER: @Description Update an existing customer by ID within the authenticated tenant
// DISABLED-SWAGGER: @Tags customers
// DISABLED-SWAGGER: @Accept json
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Customer ID"
// DISABLED-SWAGGER: @Param request body models.CustomerUpdateRequest true "Customer update data"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	// Get user from context for tenant isolation
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	var req models.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	ctx := c.Request.Context()
	customer, err := h.repo.GetByID(ctx, uint(id), user.TenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
		return
	}

	// Update fields if provided
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Street != "" {
		customer.Street = req.Street
	}
	if req.Zip != "" {
		customer.Zip = req.Zip
	}
	if req.City != "" {
		customer.City = req.City
	}
	if req.Country != "" {
		customer.Country = req.Country
	}
	if req.TaxID != "" {
		customer.TaxID = req.TaxID
	}
	if req.VAT != "" {
		customer.VAT = req.VAT
	}
	if req.PlanID != nil {
		ok, err := h.repo.PlanExists(ctx, *req.PlanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to verify plan", err.Error()))
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Plan not found", "Invalid plan ID"))
			return
		}
		customer.PlanID = *req.PlanID
	}
	if req.Status != "" {
		customer.Status = req.Status
	}
	if req.PaymentMethod != "" {
		customer.PaymentMethod = req.PaymentMethod
	}
	if req.Active != nil {
		customer.Active = *req.Active
	}

	if err := h.repo.Update(c.Request.Context(), customer); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update customer", err.Error()))
		return
	}

	// Note: Plan and Tenant relations temporarily disabled due to GORM relation issues
	// h.db.Preload("Plan").Preload("Tenant").First(&customer, customer.ID)

	c.JSON(http.StatusOK, models.SuccessResponse("Customer updated successfully", customer.ToResponse()))
}

// DeleteCustomer deletes a customer (soft delete)
// DISABLED-SWAGGER: @Summary Delete a customer
// DISABLED-SWAGGER: @Description Soft delete a customer by ID within the authenticated tenant
// DISABLED-SWAGGER: @Tags customers
// DISABLED-SWAGGER: @Produce json
// DISABLED-SWAGGER: @Security BearerAuth
// DISABLED-SWAGGER: @Param id path int true "Customer ID"
// DISABLED-SWAGGER: @Success 200 {object} models.APIResponse
// DISABLED-SWAGGER: @Failure 400 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Failure 404 {object} models.ErrorResponse
// DISABLED-SWAGGER: @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	// Get user from context for tenant isolation
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}
	user := userInterface.(*models.User)

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	customer, err := h.repo.GetByID(c.Request.Context(), uint(id), user.TenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
		return
	}

	if err := h.repo.Delete(c.Request.Context(), customer); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to delete customer", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Customer deleted successfully", nil))
}
