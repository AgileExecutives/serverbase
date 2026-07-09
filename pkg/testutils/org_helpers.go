package testutils

import orgrepo "github.com/AgileExecutives/serverbase/internal/organizations/repo"

// NewInMemoryOrganizationRepo returns an OrganizationRepo suitable for unit tests.
func NewInMemoryOrganizationRepo() *orgrepo.InMemoryOrganizationRepo {
	return orgrepo.NewInMemoryOrganizationRepo()
}
