package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
)

// Test verification token generation and applying the verification change via AuthService
func TestAuthService_VerificationTokenFlow(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()

	user := &models.User{Username: "verifyuser", Email: "vt@example.com", FirstName: "V", LastName: "T", PasswordHash: "h"}
	require.NoError(t, svc.SaveUser(ctx, user))

	token, err := auth.GenerateVerificationToken(user.Email, user.ID)
	require.NoError(t, err)

	uid, email, err := auth.ValidateVerificationToken(token)
	require.NoError(t, err)
	require.Equal(t, user.ID, uid)
	require.Equal(t, user.Email, email)

	// Mark verified and persist via service
	user.EmailVerified = true
	require.NoError(t, svc.SaveUser(ctx, user))

	got, err := svc.FindByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.EmailVerified)
}
