package middleware

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware is a lightweight placeholder for the user module
type AuthMiddleware struct {
	db     *gorm.DB
	logger core.Logger
}

func NewAuthMiddleware(ctx core.ModuleContext, logger core.Logger) *AuthMiddleware {
	return &AuthMiddleware{db: ctx.DB, logger: logger}
}

type AuthMiddlewareProvider struct{ mw *AuthMiddleware }

func NewAuthMiddlewareProvider(mw *AuthMiddleware) core.MiddlewareProvider {
	return &AuthMiddlewareProvider{mw: mw}
}

func (p *AuthMiddlewareProvider) Name() string { return "auth" }
func (p *AuthMiddlewareProvider) Handler() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
func (p *AuthMiddlewareProvider) Priority() int     { return 0 }
func (p *AuthMiddlewareProvider) ApplyTo() []string { return []string{} }
