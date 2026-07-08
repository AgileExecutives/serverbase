package handlers_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/modules/base/handlers"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	userservices "github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Handler-level tests for auth endpoints: register, logout, verify-email
func TestAuthEndpoints_Register_Logout_VerifyEmail(t *testing.T) {
	// Ensure JWT secret and disable email verification for register flow
	os.Setenv("JWT_SECRET", "test-secret-for-handlers")
	os.Setenv("FEATURE_EMAIL_VERIFICATION", "false")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("FEATURE_EMAIL_VERIFICATION")

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(db)

	// Migrate models used by these handlers
	testutils.MigrateTestDB(t, db, &models.TokenBlacklist{}, &models.Newsletter{})

	// Ensure auth package uses the test JWT secret
	auth.SetJWTSecret(os.Getenv("JWT_SECRET"))

	logger := testutils.NewMockLogger()
	userRepo := userrepo.NewGormUserRepo(db)
	authSvc := userservices.NewAuthServiceWithRepo(userRepo, nil, userRepo, userRepo, logger)
	ctx := core.ModuleContext{DB: db, Logger: logger}
	authHandlers := handlers.NewAuthHandlers(ctx, authSvc, logger)

	router := testutils.SetupTestRouter()

	// Register endpoint (creates tenant & user)
	router.POST("/auth/register", authHandlers.Register)

	payload := map[string]interface{}{
		"username":          "newuser",
		"email":             "handler.new@example.com",
		"password":          "Pass12345!",
		"accept_terms":      true,
		"company_name":      "HandlerCo",
		"first_name":        "Handler",
		"last_name":         "Test",
		"newsletter_opt_in": true,
	}

	w := testutils.MakeJSONRequest(t, router, "POST", "/auth/register", payload)
	require.Equal(t, 201, w.Code)

	// Logout endpoint (requires auth middleware); create user and token
	tenant := testutils.CreateTestTenant(t, db, "HandlerTenant")
	user := testutils.CreateTestUser(t, db, "logout@example.com", "hashed-password", tenant.ID)
	// ensure user is active & verified
	db.Model(&user).Updates(map[string]interface{}{"active": true, "email_verified": true})

	token, err := auth.GenerateJWT(user.ID, user.TenantID, user.Role)
	require.NoError(t, err)

	// For testing logout we inject the auth context directly to avoid
	// dependency on middleware timing in the test harness.
	router.POST("/auth/logout", func(c *gin.Context) {
		c.Set("token", token)
		c.Set("user", user)
		authHandlers.Logout(c)
	})

	wr := testutils.MakeJSONRequest(t, router, "POST", "/auth/logout", nil)
	require.Equal(t, 200, wr.Code)

	// Verify token was blacklisted
	tokenID, _, err := auth.ParseTokenClaims(token)
	require.NoError(t, err)
	var tb models.TokenBlacklist
	err = db.Where("token_id = ?", tokenID).First(&tb).Error
	require.NoError(t, err)
	require.Equal(t, tokenID, tb.TokenID)

	// Verify-email endpoint: create unverified user and confirm via token
	u2 := testutils.CreateTestUser(t, db, "verifyme@example.com", "h", tenant.ID)
	db.Model(&u2).Updates(map[string]interface{}{"email_verified": false})

	vtoken, err := auth.GenerateVerificationToken(u2.Email, u2.ID)
	require.NoError(t, err)

	// Route for verify email
	router.GET("/auth/verify-email/:token", authHandlers.VerifyEmail)
	path := fmt.Sprintf("/auth/verify-email/%s", vtoken)
	wv := testutils.MakeJSONRequest(t, router, "GET", path, nil)
	require.Equal(t, 200, wv.Code)
	t.Logf("verify response body: %s", wv.Body.String())

	// Handler responded with success; inspect DB record by ID to confirm persistence
	var resp map[string]interface{}
	testutils.ParseJSONResponse(t, wv, &resp)
	require.Equal(t, "Email verified successfully. You can now log in.", resp["message"])

	var after models.User
	err = db.First(&after, u2.ID).Error
	require.NoError(t, err)
	t.Logf("user after verify by id: %+v", after)
	require.True(t, after.EmailVerified, "expected user.EmailVerified to be true after handler VerifyEmail")
}
