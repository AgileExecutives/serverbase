package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

type PublicEmailRoutes struct{ emailHandler *EmailHandler }

func NewPublicEmailRoutes(emailHandler *EmailHandler) *PublicEmailRoutes {
	return &PublicEmailRoutes{emailHandler: emailHandler}
}
func (h *PublicEmailRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	router.GET("/emails/latest-emails", h.emailHandler.GetLatestEmails)
}
func (h *PublicEmailRoutes) GetPrefix() string                { return "" }
func (h *PublicEmailRoutes) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (h *PublicEmailRoutes) GetSwaggerTags() []string         { return []string{"emails"} }
