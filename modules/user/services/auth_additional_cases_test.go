package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	tenantrepo "github.com/AgileExecutives/serverbase/modules/tenant/repo"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestValidateVerificationToken_Expired(t *testing.T) {
	// ensure deterministic secret
	auth.SetJWTSecret("test-secret-verify")

	// build expired verification token
	claims := auth.JWTClaims{
		UserID:    123,
		TokenType: "verification",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "e@x.com",
			ID:        "verify_expired",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tok.SignedString([]byte("test-secret-verify"))
	require.NoError(t, err)

	_, _, err = auth.ValidateVerificationToken(tokenStr)
	require.Error(t, err)
}

type failingTenantService struct {
	err error
}

func (f *failingTenantService) CreateTenant(ctx context.Context, req models.TenantCreateRequest) (*models.Tenant, error) {
	return nil, f.err
}

func TestAuthService_CreateTenant_PropagatesError(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	// wire a tenant service that fails
	svc.SetTenantService(&failingTenantService{err: errors.New("boom")})

	tnt := &models.Tenant{CustomerID: 1, Name: "Acme", Slug: "acme"}
	err := svc.CreateTenant(context.Background(), tnt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "boom")
}

func TestAuthService_DuplicateNewsletter_Allowed(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	n := &models.Newsletter{Name: "Dup", Email: "dup@example.com", Source: "t"}
	require.NoError(t, svc.SaveNewsletter(ctx, n))
	require.NoError(t, svc.SaveNewsletter(ctx, n))

	list := repo.ListNewsletters()
	require.Len(t, list, 2)
}

func TestAuthService_DuplicateTenantSlugs_Listed(t *testing.T) {
	tr := tenantrepo.NewInMemoryTenantRepo()
	// save two tenants with same slug
	t1 := &models.Tenant{CustomerID: 1, Name: "A", Slug: "same"}
	t2 := &models.Tenant{CustomerID: 2, Name: "B", Slug: "same"}
	require.NoError(t, tr.Save(context.Background(), t1))
	require.NoError(t, tr.Save(context.Background(), t2))

	svc := NewAuthServiceWithRepo(nil, tr, nil, nil, testutils.NewMockLogger())
	slugs, err := svc.ListTenantSlugs(context.Background())
	require.NoError(t, err)
	// both entries should be present (duplicate slug values)
	require.Len(t, slugs, 2)
	require.Equal(t, "same", slugs[0])
	require.Equal(t, "same", slugs[1])
}
