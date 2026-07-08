package handlers

import (
	"os"
	"testing"

	"github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_WithInMemoryRepo(t *testing.T) {
	// Prepare fake repo and service
	fake := repo.NewInMemoryUserRepo()
	logger := testutils.NewMockLogger()
	authSvc := services.NewAuthServiceWithRepo(fake, nil, fake, fake, logger)

	// create user in fake repo
	password := "Secret123!"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{Email: "fake@example.com", PasswordHash: string(hashed), Active: true, EmailVerified: true}
	_ = fake.Save(nil, u)

	// create handler with nil DB (login does not require DB access)
	ctx := core.ModuleContext{DB: nil, Logger: logger}
	h := NewAuthHandlers(ctx, authSvc, logger)
	router := testutils.SetupTestRouter()
	router.POST("/auth/login", h.Login)

	payload := map[string]interface{}{"email": u.Email, "password": password}
	w := testutils.MakeJSONRequest(t, router, "POST", "/auth/login", payload)
	require.Equal(t, 200, w.Code)
}

func TestRegister_ExistingUser_Path_UsesAuthService(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	fake := repo.NewInMemoryUserRepo()
	logger := testutils.NewMockLogger()
	authSvc := services.NewAuthServiceWithRepo(fake, nil, fake, fake, logger)

	// existing user
	oldPass := "OldPass123!"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(oldPass), bcrypt.DefaultCost)
	u := &models.User{Email: "exist@example.com", PasswordHash: string(hashed), Active: true, EmailVerified: true}
	_ = fake.Save(nil, u)

	ctx := core.ModuleContext{DB: nil, Logger: logger}
	h := NewAuthHandlers(ctx, authSvc, logger)
	router := testutils.SetupTestRouter()
	router.POST("/auth/register", h.Register)

	// register request for same email should update password and return token
	newPass := "NewSecret123!"
	payload := map[string]interface{}{"email": u.Email, "password": newPass, "accept_terms": true, "username": "existuser", "first_name": "Exist", "last_name": "User"}
	w := testutils.MakeJSONRequest(t, router, "POST", "/auth/register", payload)
	require.Equal(t, 200, w.Code)
}
