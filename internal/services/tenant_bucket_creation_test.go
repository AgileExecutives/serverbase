package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/modules/tenant/repo"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeBucketSvc struct {
	created []uint
}

func (f *fakeBucketSvc) CreateTenantBucket(ctx context.Context, tenantID uint) error {
	f.created = append(f.created, tenantID)
	return nil
}

func (f *fakeBucketSvc) BucketExists(ctx context.Context, tenantID uint) (bool, error) {
	return false, nil
}

func TestTenantService_Create_InvokesBucketCreation(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(db)

	r := repo.NewGormTenantRepo(db)
	fake := &fakeBucketSvc{}
	svc := NewTenantService(r, fake)

	req := models.TenantCreateRequest{Name: "Bucketed Co"}
	tenant, err := svc.CreateTenant(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, tenant)
	assert.NotZero(t, tenant.ID)
	// fake should have been called with the created tenant ID
	assert.Len(t, fake.created, 1)
	assert.Equal(t, tenant.ID, fake.created[0])
}
