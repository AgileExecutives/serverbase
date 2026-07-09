package services

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestAuthService_ConcurrentBlacklist(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	count := 50
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		i := i
		go func(n int) {
			defer wg.Done()
			tb := &models.TokenBlacklist{TokenID: "jti-" + strconv.Itoa(n), UserID: uint(n)}
			_ = svc.BlacklistToken(ctx, tb)
		}(i)
	}
	wg.Wait()

	bl := repo.ListBlacklists()
	require.Len(t, bl, count)
}

func TestAuthService_ConcurrentNewsletterSaves(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := userrepo.NewInMemoryUserRepo()
	svc := NewAuthServiceWithRepo(repo, nil, repo, repo, logger)

	ctx := context.Background()
	count := 50
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			n := &models.Newsletter{Name: "N", Email: "race@example.com", Source: "test"}
			_ = svc.SaveNewsletter(ctx, n)
		}()
	}
	wg.Wait()

	list := repo.ListNewsletters()
	require.Len(t, list, count)
}
