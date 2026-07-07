package services

import (
	"context"
)

// Template represents a minimal template metadata struct
type Template struct {
	ID   uint
	Name string
}

// TemplateService provides template rendering and retrieval.
type TemplateService struct{}

func NewTemplateService() *TemplateService { return &TemplateService{} }

func (s *TemplateService) RenderTemplate(ctx context.Context, tenantID int, templateID uint, data interface{}) (string, error) {
	// minimal renderer: return placeholder HTML
	return "<html><body>Rendered</body></html>", nil
}

func (s *TemplateService) GetTemplate(ctx context.Context, tenantID int, templateID uint) (*Template, error) {
	return &Template{ID: templateID, Name: "default"}, nil
}

// CopyTemplatesFromTenant2Org2 copies default templates from tenant 2/org 2 to a new org.
// This is a lightweight stub used by the server-test harness.
func (s *TemplateService) CopyTemplatesFromTenant2Org2(ctx context.Context, tenantID uint, orgID uint) error {
	// In the real implementation this would duplicate templates and storage entries.
	return nil
}

// ContractRegistrar is a helper for registering template contracts.
type ContractRegistrar struct{}

// NewContractRegistrar creates a new ContractRegistrar.
func NewContractRegistrar() *ContractRegistrar { return &ContractRegistrar{} }

// RegisterContractFromFile registers a contract file for a tenant and module.
func (r *ContractRegistrar) RegisterContractFromFile(tenantID uint, module string, path string) error {
	// Stub: in real code this would read the file and register it.
	return nil
}
