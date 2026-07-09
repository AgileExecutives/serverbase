package handlers_test

import (
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	orghandlers "github.com/AgileExecutives/serverbase/internal/organizations/handlers"
	orgrepo "github.com/AgileExecutives/serverbase/internal/organizations/repo"
	orgservices "github.com/AgileExecutives/serverbase/internal/organizations/services"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOrganizationHandler_CreateAndList_WithInMemoryRepo(t *testing.T) {
	r := orgrepo.NewInMemoryOrganizationRepo()
	svc := orgservices.NewOrganizationServiceWithRepo(r)
	h := orghandlers.NewOrganizationHandler(svc)

	router := testutils.SetupTestRouter()

	// POST /organizations with injected user context
	router.POST("/organizations", func(c *gin.Context) {
		// inject authenticated user for tenant isolation
		c.Set("user", &models.User{ID: 1, TenantID: 1, Active: true})
		h.CreateOrganization(c)
	})

	payload := map[string]interface{}{
		"name":       "Test Org",
		"owner_name": "Owner",
		"email":      "org@example.com",
	}

	w := testutils.MakeJSONRequest(t, router, "POST", "/organizations", payload)
	require.Equal(t, 201, w.Code)

	// Now list organizations
	router.GET("/organizations", func(c *gin.Context) {
		c.Set("user", &models.User{ID: 1, TenantID: 1, Active: true})
		h.GetAllOrganizations(c)
	})

	wl := testutils.MakeJSONRequest(t, router, "GET", "/organizations", nil)
	require.Equal(t, 200, wl.Code)
}
