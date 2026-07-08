package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

type fakeTenantCreator struct {
	called  bool
	lastReq models.TenantCreateRequest
}

func (f *fakeTenantCreator) CreateTenant(ctx context.Context, req models.TenantCreateRequest) (*models.Tenant, error) {
	f.called = true
	f.lastReq = req
	return &models.Tenant{ID: 999, Name: req.Name, Slug: req.Slug, CustomerID: req.CustomerID}, nil
}

func TestAuthService_CreateTenant_DelegatesToTenantService(t *testing.T) {
	logger := testutils.NewMockLogger()
	svc := NewAuthServiceWithRepo(nil, nil, nil, nil, logger)
	fake := &fakeTenantCreator{}
	svc.SetTenantService(fake)

	tptr := &models.Tenant{CustomerID: 10, Name: "Acme", Slug: "acme"}
	err := svc.CreateTenant(context.Background(), tptr)
	assert.NoError(t, err)
	assert.True(t, fake.called)
	assert.Equal(t, uint(999), tptr.ID)
	assert.Equal(t, "Acme", fake.lastReq.Name)
}
