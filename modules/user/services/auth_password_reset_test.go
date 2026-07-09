package services

import (
	"context"
	"testing"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
)

// Test password reset flow: generate reset token, validate, apply new password via AuthService
func TestAuthService_PasswordResetFlow(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	user := &models.User{Email: "reset@example.com", FirstName: "R", LastName: "User", PasswordHash: "oldhash"}
	require.NoError(t, svc.SaveUser(ctx, user))

	// generate a reset token valid for a short duration
	token, err := auth.GenerateResetToken(user.Email, 1*time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// validate token and extract email
	gotEmail, err := auth.ValidateResetToken(token)
	require.NoError(t, err)
	require.Equal(t, user.Email, gotEmail)

	// simulate user posting new password: update password hash and persist via service
	user.PasswordHash = "newhash"
	require.NoError(t, svc.SaveUser(ctx, user))

	// retrieve and verify password changed
	got, err := svc.FindByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "newhash", got.PasswordHash)
}
