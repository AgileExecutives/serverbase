package testutils_test

import (
	"net/http"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	basehandlers "github.com/AgileExecutives/serverbase/modules/base/handlers"
	emailhandlers "github.com/AgileExecutives/serverbase/modules/email/handlers"
	emailservices "github.com/AgileExecutives/serverbase/modules/email/services"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	userservices "github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
)

// Unit example: service logic with in-memory repo
func TestAuthService_WithInMemoryRepo(t *testing.T) {
	repo := userrepo.NewInMemoryUserRepo()
	logger := testutils.NewMockLogger()
	svc := userservices.NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	u := &models.User{Username: "u1", Email: "u1@example.com", PasswordHash: "h"}
	require.NoError(t, svc.SaveUser(nil, u))

	found, err := svc.FindByEmail(nil, u.Email)
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, u.Email, found.Email)
}

// Email handler example: use in-memory email repo and service
func TestEmailHandler_WithInMemoryRepo(t *testing.T) {
	repo := testutils.NewInMemoryEmailRepo()
	emailSvc := emailservices.NewEmailService()
	h := emailhandlers.NewEmailHandler(repo, emailSvc)

	r := testutils.SetupTestRouter()
	r.POST("/emails/send", h.SendEmail)

	payload := map[string]interface{}{
		"to":      "a@example.com",
		"from":    "b@example.com",
		"subject": "hello",
		"body":    "body",
	}
	w := testutils.MakeJSONRequest(t, r, "POST", "/emails/send", payload)
	require.Equal(t, http.StatusCreated, w.Code)
}

// Service integration example: sqlite in-memory DB
func TestService_WithSQLiteDB(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(db)
	testutils.MigrateTestDB(t, db, &models.User{})

	userRepo := userrepo.NewGormUserRepo(db)
	logger := testutils.NewMockLogger()
	svc := userservices.NewAuthServiceWithRepo(userRepo, nil, userRepo, userRepo, logger)

	u := &models.User{Username: "u2", Email: "u2@example.com", PasswordHash: "h"}
	require.NoError(t, svc.SaveUser(nil, u))

	got, err := svc.FindByEmail(nil, u.Email)
	require.NoError(t, err)
	require.Equal(t, u.Email, got.Email)
}

// Handler example: basic httptest using SetupTestRouter
func TestHandler_RegisterEndpointSkeleton(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(db)
	testutils.MigrateTestDB(t, db, &models.User{}, &models.Tenant{})

	logger := testutils.NewMockLogger()
	userRepo := userrepo.NewGormUserRepo(db)
	authSvc := userservices.NewAuthServiceWithRepo(userRepo, nil, userRepo, userRepo, logger)
	ctx := core.ModuleContext{DB: db, Logger: logger}
	h := basehandlers.NewAuthHandlers(ctx, authSvc, logger)

	r := testutils.SetupTestRouter()
	r.POST("/auth/register", h.Register)

	payload := map[string]interface{}{
		"username":     "ruser",
		"email":        "r@example.com",
		"password":     "Pass123!",
		"first_name":   "R",
		"last_name":    "User",
		"company_name": "RCo",
		"accept_terms": true,
	}
	w := testutils.MakeJSONRequest(t, r, "POST", "/auth/register", payload)
	require.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK)
}
