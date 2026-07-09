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

func TestValidateResetToken_Expired(t *testing.T) {
	// create a token that is already expired
	token, err := auth.GenerateResetToken("a@b.com", -1*time.Second)
	require.NoError(t, err)

	_, err = auth.ValidateResetToken(token)
	require.Error(t, err)
}

func TestValidateResetToken_Invalid(t *testing.T) {
	_, err := auth.ValidateResetToken("not-a-token")
	require.Error(t, err)
}

func TestAuthService_SaveNewsletter_NoRepo(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	// newsletterRepo is nil -> expect error
	svc := NewAuthServiceWithRepo(repo, nil, nil, repo, logger)

	ctx := context.Background()
	n := &models.Newsletter{Name: "N", Email: "x@x.com", Source: "t"}
	err := svc.SaveNewsletter(ctx, n)
	require.Error(t, err)
	require.Contains(t, err.Error(), "newsletter repo not provided")
}

func TestAuthService_BlacklistToken_NoRepo(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	// tokenRepo is nil -> expect error
	svc := NewAuthServiceWithRepo(repo, nil, repo, nil, logger)

	ctx := context.Background()
	tb := &models.TokenBlacklist{TokenID: "t1", UserID: 1}
	err := svc.BlacklistToken(ctx, tb)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token blacklist repo not provided")
}

func TestAuthService_FindByEmail_NoUserRepo_NoPanic(t *testing.T) {
	logger := testutils.NewMockLogger()
	// construct via NewAuthService (db-based constructor) which leaves userRepo nil
	svc := NewAuthService(nil, logger)

	ctx := context.Background()
	u, err := svc.FindByEmail(ctx, "noone@example.com")
	require.NoError(t, err)
	require.Nil(t, u)

	// SaveUser should be a no-op and not error
	err = svc.SaveUser(ctx, &models.User{Email: "x@x.com"})
	require.NoError(t, err)
}
