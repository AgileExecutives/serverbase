package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestAuthService_RegisterAndNewsletter(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	user := &models.User{Email: "new@example.com", FirstName: "New", LastName: "User"}
	reqErr := svc.SaveUser(ctx, user)
	require.NoError(t, reqErr)
	retrieved, err := svc.FindByEmail(ctx, "new@example.com")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	require.Equal(t, "new@example.com", retrieved.Email)

	// Newsletter
	n := &models.Newsletter{Name: "New User", Email: "new@example.com", Source: "test"}
	err = svc.SaveNewsletter(ctx, n)
	require.NoError(t, err)

	list := repo.ListNewsletters()
	require.Len(t, list, 1)
	require.Equal(t, "new@example.com", list[0].Email)
}

func TestAuthService_Logout_Blacklist(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	tb := &models.TokenBlacklist{TokenID: "jti-1", UserID: 1}
	err := svc.BlacklistToken(ctx, tb)
	require.NoError(t, err)

	bl := repo.ListBlacklists()
	require.Len(t, bl, 1)
	require.Equal(t, "jti-1", bl[0].TokenID)
}

func TestAuthService_VerifyEmail_SaveUser(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	user := &models.User{Email: "verify@example.com", FirstName: "V", LastName: "E"}
	require.NoError(t, svc.SaveUser(ctx, user))

	// Simulate verification change
	user.EmailVerified = true
	require.NoError(t, svc.SaveUser(ctx, user))

	got, err := svc.FindByEmail(ctx, "verify@example.com")
	require.NoError(t, err)
	require.True(t, got.EmailVerified)
}
