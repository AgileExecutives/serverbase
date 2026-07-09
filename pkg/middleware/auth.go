package middleware

import (
	internalMiddleware "github.com/AgileExecutives/serverbase/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware returns authentication middleware that can be used by modules
// This is a wrapper around the internal auth middleware to make it accessible
// to modules without exposing internal packages
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return internalMiddleware.AuthMiddleware(db)
}
