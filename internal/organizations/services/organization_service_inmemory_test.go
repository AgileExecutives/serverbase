package services

import (
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	orgrepo "github.com/AgileExecutives/serverbase/internal/organizations/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeCreateOrgReqInMem() models.CreateOrganizationRequest {
	return models.CreateOrganizationRequest{
		Name:      "InMemory Org",
		OwnerName: "Owner",
		Email:     "inmem@example.com",
	}
}

func TestOrganizationService_InMemory_CRUD(t *testing.T) {
	r := orgrepo.NewInMemoryOrganizationRepo()
	svc := NewOrganizationServiceWithRepo(r)

	// Create
	org, err := svc.CreateOrganization(makeCreateOrgReqInMem(), 10)
	require.NoError(t, err)
	require.NotNil(t, org)
	assert.Equal(t, uint(10), org.TenantID)

	// Get
	fetched, err := svc.GetOrganizationByID(org.ID, 10)
	require.NoError(t, err)
	assert.Equal(t, org.ID, fetched.ID)

	// List
	list, total, err := svc.GetOrganizations(1, 10, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)

	// Update
	newName := "Updated InMemory Org"
	updated, err := svc.UpdateOrganization(org.ID, 10, models.UpdateOrganizationRequest{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// Delete wrong tenant should fail
	err = svc.DeleteOrganization(org.ID, 11)
	require.Error(t, err)

	// Delete correct tenant
	err = svc.DeleteOrganization(org.ID, 10)
	require.NoError(t, err)

	// Ensure deleted
	_, err = svc.GetOrganizationByID(org.ID, 10)
	require.Error(t, err)
}
