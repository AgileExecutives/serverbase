package api

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/middleware"
	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Re-export internal response types for external modules and handlers
type APIResponse = models.APIResponse
type PaginationResponse = models.PaginationResponse
type ListResponse = models.ListResponse
type ErrorResponse = models.ErrorResponse

// ErrorResponseFunc returns a JSON error payload used by handlers.
func ErrorResponseFunc(title, detail string) map[string]interface{} {
	return map[string]interface{}{"success": false, "error": title, "detail": detail}
}

// SuccessResponse returns a success envelope with data.
func SuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"success": true, "message": message, "data": data}
}

// SuccessListResponse returns a paginated list envelope.
func SuccessListResponse(items interface{}, page, limit, total int) map[string]interface{} {
	return map[string]interface{}{"success": true, "items": items, "page": page, "limit": limit, "total": total}
}

// GetTenantID extracts tenant id from context or headers.
func GetTenantID(c *gin.Context) (uint, error) {
	return middleware.GetTenantID(c)
}

// SuccessMessageResponse returns a simple success envelope with only a message.
func SuccessMessageResponse(message string) map[string]interface{} {
	return map[string]interface{}{"success": true, "message": message}
}

// GetUser returns a User from context (set by auth middleware) for handlers.
func GetUser(c *gin.Context) (*models.User, error) {
	return middleware.GetUser(c)
}

// AuthMiddleware exposes the internal auth middleware for routes that
// import the api package (convenience wrapper).
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return middleware.AuthMiddleware(db)
}

// GetUserID retrieves the authenticated user's ID from the context.
func GetUserID(c *gin.Context) (uint, error) {
	u, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

// ModuleRouteProvider is a legacy compatibility interface used by older modules
// that expect a simple route provider with prefix and registration methods.
type ModuleRouteProvider interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetPrefix() string
	GetMiddleware() []gin.HandlerFunc
	GetSwaggerTags() []string
}

// Helper to write JSON errors
func JSONError(c *gin.Context, status int, title, detail string) {
	c.JSON(status, ErrorResponseFunc(title, detail))
}

// Convenience wrapper for standard library http handlers if needed
func WriteHTTPError(w http.ResponseWriter, status int, title, detail string) {
	w.WriteHeader(status)
	w.Write([]byte(title + ": " + detail))
}
