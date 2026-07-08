package testutils

import "context"

// FixedTenantProvider returns a constant tenant id for tests.
type FixedTenantProvider struct {
	tenantID string
}

func NewFixedTenantProvider(id string) *FixedTenantProvider {
	return &FixedTenantProvider{tenantID: id}
}

// Current implements the TenantProvider interface used in production code.
func (f *FixedTenantProvider) Current(ctx context.Context) (string, error) {
	return f.tenantID, nil
}
