package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/stretchr/testify/require"
)

func TestInMemoryUserRepo_NonDedupByDefault(t *testing.T) {
	r := NewInMemoryUserRepo()
	ctx := context.Background()

	n := &models.Newsletter{Name: "A", Email: "dup@example.com", Source: "t"}
	require.NoError(t, r.SaveNewsletter(ctx, n))
	require.NoError(t, r.SaveNewsletter(ctx, n))

	list := r.ListNewsletters()
	require.Len(t, list, 2)
}

func TestInMemoryUserRepo_DedupeEnabled(t *testing.T) {
	r := NewInMemoryUserRepo()
	r.SetDeduplicateNewsletters(true)
	ctx := context.Background()

	n1 := &models.Newsletter{Name: "First", Email: "uniq@example.com", Source: "t1"}
	n2 := &models.Newsletter{Name: "Second", Email: "uniq@example.com", Source: "t2"}

	require.NoError(t, r.SaveNewsletter(ctx, n1))
	require.NoError(t, r.SaveNewsletter(ctx, n2))

	list := r.ListNewsletters()
	require.Len(t, list, 1)
	require.Equal(t, "Second", list[0].Name)
	require.Equal(t, "t2", list[0].Source)
}
