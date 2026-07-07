package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handlers *AuthHandlers
}

func NewAuthRoutes(handlers *AuthHandlers) core.RouteProvider {
	return &AuthRoutes{handlers: handlers}
}

func (r *AuthRoutes) GetPrefix() string                { return "/auth" }
func (r *AuthRoutes) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *AuthRoutes) GetSwaggerTags() []string         { return []string{"authentication"} }

func (r *AuthRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	router.POST("/login", r.handlers.Login)
	router.POST("/register", r.handlers.Register)
	router.POST("/logout", ctx.Auth.RequireAuth(), r.handlers.Logout)
	router.POST("/refresh", ctx.Auth.RequireAuth(), r.handlers.RefreshToken)
	router.GET("/me", ctx.Auth.RequireAuth(), r.handlers.Me)
	router.POST("/change-password", ctx.Auth.RequireAuth(), r.handlers.ChangePassword)
	router.GET("/verify-email/:token", r.handlers.VerifyEmail)
	router.GET("/check-verification-token/:token", r.handlers.CheckVerificationToken)
	router.POST("/forgot-password", r.handlers.ForgotPassword)
	router.GET("/check-reset-token/:token", r.handlers.CheckResetToken)
	router.POST("/new-password/:token", r.handlers.ResetPassword)
	router.GET("/password-security", r.handlers.GetPasswordSecurity)
}

type ContactRoutes struct{ handlers *ContactHandlers }

func NewContactRoutes(handlers *ContactHandlers) core.RouteProvider {
	return &ContactRoutes{handlers: handlers}
}
func (r *ContactRoutes) GetPrefix() string                { return "/contacts" }
func (r *ContactRoutes) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *ContactRoutes) GetSwaggerTags() []string         { return []string{"contact-form", "newsletter"} }
func (r *ContactRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	router.POST("/form", r.handlers.SubmitContactForm)
	router.GET("", ctx.Auth.RequireAuth(), r.handlers.GetContacts)
	router.GET("/:id", ctx.Auth.RequireAuth(), r.handlers.GetContact)
	router.POST("", ctx.Auth.RequireAuth(), r.handlers.CreateContact)
	router.PUT("/:id", ctx.Auth.RequireAuth(), r.handlers.UpdateContact)
	router.DELETE("/:id", ctx.Auth.RequireAuth(), r.handlers.DeleteContact)
	router.GET("/newsletter", ctx.Auth.RequireAuth(), r.handlers.GetNewsletterSubscriptions)
	router.DELETE("/newsletter/unsubscribe", ctx.Auth.RequireAuth(), r.handlers.UnsubscribeFromNewsletter)
}

type UserSettingsRoutes struct{ handlers *UserSettingsHandlers }

func NewUserSettingsRoutes(handlers *UserSettingsHandlers) core.RouteProvider {
	return &UserSettingsRoutes{handlers: handlers}
}
func (r *UserSettingsRoutes) GetPrefix() string                { return "/user-settings" }
func (r *UserSettingsRoutes) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *UserSettingsRoutes) GetSwaggerTags() []string         { return []string{"user-settings"} }
func (r *UserSettingsRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	router.GET("", ctx.Auth.RequireAuth(), r.handlers.GetUserSettings)
	router.PUT("", ctx.Auth.RequireAuth(), r.handlers.UpdateUserSettings)
	router.POST("/reset", ctx.Auth.RequireAuth(), r.handlers.ResetUserSettings)
}

// HealthRoutes provides health check routes
type HealthRoutes struct{ handlers *HealthHandlers }

func NewHealthRoutes(handlers *HealthHandlers) core.RouteProvider {
	return &HealthRoutes{handlers: handlers}
}

func (r *HealthRoutes) GetPrefix() string                { return "/health" }
func (r *HealthRoutes) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *HealthRoutes) GetSwaggerTags() []string         { return []string{"health"} }
func (r *HealthRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	router.GET("", r.handlers.HealthCheck)
	router.GET("/ping", r.handlers.Ping)
}
