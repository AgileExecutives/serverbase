package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/stretchr/testify/require"
)

func TestInMemoryEmailRepo_BasicFlow(t *testing.T) {
	r := NewInMemoryEmailRepo()
	ctx := context.Background()

	e := &models.Email{To: "a@x.com", From: "b@x.com", Subject: "s", Body: "b", Status: "pending"}
	require.NoError(t, r.Create(ctx, e))
	require.NotZero(t, e.ID)

	got, err := r.FindByID(ctx, e.ID)
	require.NoError(t, err)
	require.Equal(t, e.To, got.To) // ensure stored

	// update status
	require.NoError(t, r.UpdateStatus(ctx, e.ID, "sent", ""))
	got2, err := r.FindByID(ctx, e.ID)
	require.NoError(t, err)
	require.Equal(t, "sent", got2.Status)

	list, total, err := r.List(ctx, 0, 10, "")
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	stats, err := r.Stats(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), stats["total"])
	require.Equal(t, int64(1), stats["sent"])
}
