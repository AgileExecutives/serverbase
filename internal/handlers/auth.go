package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	emailServices "github.com/AgileExecutives/serverbase/modules/email/services"
	usersvc "github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/config"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService *usersvc.AuthService
	cfg         config.Config
}

// NewAuthHandler creates a new auth handler wired to AuthService
func NewAuthHandler(authSvc *usersvc.AuthService) *AuthHandler {
	return &AuthHandler{authService: authSvc, cfg: config.Load()}
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	userPtr, err := h.authService.FindByEmail(c.Request.Context(), req.Email)
	if err != nil || userPtr == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid credentials", "User not found"))
		return
	}
	user := *userPtr

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Account disabled", "User account is not active"))
		return
	}

	// Check if email is verified (only if email verification is required)
	if h.cfg.Email.RequireEmailVerification && !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Email not verified", "Please verify your email address before logging in"))
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid credentials", "Password mismatch"))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.TenantID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate token", err.Error()))
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", response))
}

// Logout blacklists the current JWT token
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from context (set by auth middleware)
	tokenString, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("No token found", "Token not provided"))
		return
	}

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Parse token to get JTI and expiration
	tokenID, expiresAt, err := auth.ParseTokenClaims(tokenString.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid token", err.Error()))
		return
	}

	// Add token to blacklist via AuthService
	blacklistEntry := models.TokenBlacklist{TokenID: tokenID, UserID: user.ID, ExpiresAt: expiresAt, Reason: "User logout"}
	if err := h.authService.BlacklistToken(c.Request.Context(), &blacklistEntry); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to logout", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Logout successful", nil))
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Validate required terms acceptance
	if !req.AcceptTerms {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Terms not accepted", "You must accept the terms and conditions to register"))
		return
	}

	// Validate password against requirements
	if err := utils.ValidatePassword(req.Password); err != nil {
		log.Printf("❌ Register: Password validation failed: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid password", err.Error()))
		return
	}

	// Check if username or email already exists
	existingUserPtr, _ := h.authService.FindByEmail(c.Request.Context(), req.Email)
	if existingUserPtr != nil {
		c.JSON(http.StatusConflict, models.ErrorResponseFunc("User already exists", "Username or email already taken"))
		return
	}

	var tenantID uint

	// Check if there's a user-signup token in the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		signupToken := authHeader[7:]

		// Try to validate as user-signup token
		if signupTenantID, err := auth.ValidateUserSignupToken(signupToken); err == nil {
			// Valid user-signup token found, join the specified tenant
			tenant, terr := h.authService.FindTenantByID(c.Request.Context(), signupTenantID)
			if terr != nil || tenant == nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Tenant not found", "Signup token references invalid tenant"))
				return
			}

			tenantID = tenant.ID

			// Set default role for joining existing tenant via invitation
			if req.Role == "" {
				req.Role = "user"
			}

			log.Printf("✅ User signup via invitation token to tenant: %s (ID: %d)", tenant.Name, tenant.ID)
		} else {
			// Not a valid user-signup token, proceed with normal registration
			if req.CompanyName == "" {
				c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Company name required", "Company name is required to create a new tenant"))
				return
			}

			// Create new tenant - logic will be handled below
		}
	} else {
		// No token provided, company_name is required to create new tenant
		if req.CompanyName == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Company name required", "Company name is required to create a new tenant"))
			return
		}
	}

	// Create new tenant if tenantID is still 0 (no valid signup token was found)
	if tenantID == 0 {
		if existingTenant, eterr := h.authService.FindTenantByName(c.Request.Context(), req.CompanyName); eterr == nil && existingTenant != nil {
			c.JSON(http.StatusConflict, models.ErrorResponseFunc("Company already exists", "A company with this name already exists"))
			return
		}

		// Generate slug and ensure uniqueness
		slug := utils.GenerateSlug(req.CompanyName)
		var existingSlugs []string
		if slugs, serr := h.authService.ListTenantSlugs(c.Request.Context()); serr == nil {
			existingSlugs = slugs
		}
		slug = utils.EnsureUniqueSlug(slug, existingSlugs)

		tenant := models.Tenant{Name: req.CompanyName, Slug: slug}
		if err := h.authService.CreateTenant(c.Request.Context(), &tenant); err != nil {
			log.Printf("❌ Register: Failed to create tenant: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create tenant", err.Error()))
			return
		}

		tenantID = tenant.ID
		// User is admin of their own tenant
		req.Role = "admin"
		log.Printf("✅ Created new tenant with MinIO bucket: %s (ID: %d, Slug: %s)", tenant.Name, tenant.ID, tenant.Slug)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to hash password", err.Error()))
		return
	}

	// Create user
	emailVerified := !h.cfg.Email.RequireEmailVerification
	var emailVerifiedAt *time.Time
	if emailVerified {
		now := time.Now()
		emailVerifiedAt = &now
	}

	user := models.User{
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Role:            req.Role,
		TenantID:        tenantID,
		Active:          true,
		EmailVerified:   emailVerified,
		EmailVerifiedAt: emailVerifiedAt,
	}

	if err := h.authService.SaveUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create user", err.Error()))
		return
	}

	// Event is published by AuthService.SaveUser for new users; no-op here.

	// Handle newsletter subscription if opted in
	if req.NewsletterOptIn {
		newsletter := models.Newsletter{Name: req.FirstName + " " + req.LastName, Email: req.Email, Interest: "general", Source: "registration"}
		if err := h.authService.SaveNewsletter(c.Request.Context(), &newsletter); err != nil {
			log.Printf("⚠️ Register: Failed to create newsletter subscription: %v (continuing anyway)", err)
		} else {
			log.Printf("✅ Newsletter subscription created for %s", req.Email)
		}
	}

	// Send verification email only if email verification is required
	if h.cfg.Email.RequireEmailVerification {
		// Generate email verification token
		verificationToken, err := auth.GenerateVerificationToken(user.Email, user.ID)
		if err != nil {
			log.Printf("❌ Register: Failed to generate verification token: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate verification token", err.Error()))
			return
		}

		// Send verification email
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173" // Default for development
		}
		verificationURL := frontendURL + "/verify-email?token=" + verificationToken

		emailService := emailServices.NewEmailService()
		err = emailService.SendVerificationEmail(user.Email, user.FirstName, verificationURL)
		if err != nil {
			log.Printf("⚠️ Register: Failed to send verification email: %v (continuing anyway)", err)
			// Don't fail registration if email fails, user can request resend
		}
	}

	// Generate onboarding token for limited access
	onboardingToken, err := auth.GenerateOnboardingToken(user.ID, user.TenantID, user.Role)
	if err != nil {
		log.Printf("❌ Register: Failed to generate onboarding token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate onboarding token", err.Error()))
		return
	}

	// Return user data with onboarding token
	response := models.LoginResponse{
		Token: onboardingToken,
		User:  user.ToResponse(),
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("User created successfully. Please check your email to verify your account.", response))
}

// ChangePassword changes the user's password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ ChangePassword: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	log.Printf("🔍 ChangePassword: Received request - has_current: %v, has_new: %v", req.CurrentPassword != "", req.NewPassword != "")

	// Validate new password against requirements
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		log.Printf("❌ ChangePassword: Password validation failed: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid password", err.Error()))
		return
	}

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid current password", "Current password is incorrect"))
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to hash password", err.Error()))
		return
	}

	// Update password via AuthService
	user.PasswordHash = string(hashedPassword)
	if err := h.authService.SaveUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update password", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password changed successfully", nil))
}

// Me returns the current user information
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Note: Tenant relation temporarily disabled due to GORM relation issues
	// h.db.Preload("Tenant").First(user, user.ID)

	c.JSON(http.StatusOK, models.SuccessResponse("User retrieved successfully", user.ToResponse()))
}

// ForgotPassword sends a password reset link to the user's email
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Check if user exists
	userPtr, _ := h.authService.FindByEmail(c.Request.Context(), req.Email)
	if userPtr == nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, models.SuccessResponse("If the email exists, a password reset link has been sent", nil))
		return
	}

	// Check if user is active
	if !userPtr.Active {
		c.JSON(http.StatusOK, models.SuccessResponse("If the email exists, a password reset link has been sent", nil))
		return
	}

	// Get reset token expiry from env (default 2 hours)
	expiryStr := os.Getenv("RESET_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "2h"
	}
	expiryDuration, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiryDuration = 2 * time.Hour
	}

	// Generate reset token
	token, err := auth.GenerateResetToken(userPtr.Email, expiryDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate reset token", err.Error()))
		return
	}

	// Build reset URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Get reset route slug from env (default: /auth/new-password/)
	resetRouteSlug := os.Getenv("RESET_PASSWORD_ROUTE")
	if resetRouteSlug == "" {
		resetRouteSlug = "/auth/new-password/"
	}

	resetURL := frontendURL + resetRouteSlug + token

	// Send email
	emailService := emailServices.NewEmailService()
	userName := userPtr.FirstName
	if userName == "" {
		userName = userPtr.Username
	}

	if err := emailService.SendPasswordResetEmail(userPtr.Email, userName, resetURL); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to send reset email", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password reset link has been sent to your email", nil))
}

// ResetPassword resets the user's password using a reset token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", "Reset token is required"))
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Validate password against requirements
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Password validation failed", err.Error()))
		return
	}

	// Validate reset token and get email
	email, err := auth.ValidateResetToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid or expired reset token", err.Error()))
		return
	}

	// Find user by email
	userPtr, err := h.authService.FindByEmail(c.Request.Context(), email)
	if err != nil || userPtr == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("User not found", "Invalid reset token"))
		return
	}

	// Check if user is active
	if !userPtr.Active {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Account disabled", "User account is not active"))
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to hash password", err.Error()))
		return
	}

	// Update password via AuthService
	userPtr.PasswordHash = string(hashedPassword)
	if err := h.authService.SaveUser(c.Request.Context(), userPtr); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update password", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password has been reset successfully", nil))
}

// GetPasswordSecurity godoc
func (h *AuthHandler) GetPasswordSecurity(c *gin.Context) {
	requirements := utils.GetPasswordRequirements()
	c.JSON(http.StatusOK, requirements)
}

// VerifyEmail verifies a user's email address using the verification token
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Missing token", "Verification token is required"))
		return
	}

	// Validate verification token
	userID, email, err := auth.ValidateVerificationToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid or expired token", err.Error()))
		return
	}

	// Find user by email and verify ID matches
	userPtr, err := h.authService.FindByEmail(c.Request.Context(), email)
	if err != nil || userPtr == nil || userPtr.ID != userID {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("User not found", "Invalid verification token"))
		return
	}

	// Check if already verified
	if userPtr.EmailVerified {
		c.JSON(http.StatusOK, models.SuccessResponse("Email already verified", nil))
		return
	}

	// Update user's email verification status
	now := time.Now()
	userPtr.EmailVerified = true
	userPtr.EmailVerifiedAt = &now

	if err := h.authService.SaveUser(c.Request.Context(), userPtr); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to verify email", err.Error()))
		return
	}

	log.Printf("✅ Email verified for user %d (%s)", userPtr.ID, userPtr.Email)
	c.JSON(http.StatusOK, models.SuccessResponse("Email verified successfully. You can now log in.", nil))
}

// GenerateSignupToken generates a signup token for inviting users to a tenant
func (h *AuthHandler) GenerateSignupToken(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Check if user has permission to invite others (admin, owner, or super-admin)
	if user.Role != "admin" && user.Role != "owner" && user.Role != "super-admin" {
		c.JSON(http.StatusForbidden, models.ErrorResponseFunc("Permission denied", "Only admins and owners can generate signup tokens"))
		return
	}

	// Generate signup token for the user's tenant
	token, err := auth.GenerateUserSignupToken(user.TenantID, user.Email)
	if err != nil {
		log.Printf("❌ GenerateSignupToken: Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate token", err.Error()))
		return
	}

	response := map[string]interface{}{
		"token":      token,
		"expires_in": "7 days",
		"tenant_id":  user.TenantID,
		"usage":      "Include this token in the Authorization header (Bearer token) when registering new users to add them to your tenant",
	}

	log.Printf("✅ Generated signup token for tenant %d by user %d (%s)", user.TenantID, user.ID, user.Email)
	c.JSON(http.StatusOK, models.SuccessResponse("Signup token generated successfully", response))
}
