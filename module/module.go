package module

import "net/http"

// Registry is a minimal surface exposed to modules so they can register routes
// and access basic platform services. Expand as needed.
type Registry interface {
	RegisterRoute(pattern string, handler http.Handler)
}

// Module is the base interface all modules should implement.
type Module interface {
	Name() string
	// Register is called during server setup so modules can register routes and
	// wire their own dependencies using the provided Registry.
	Register(reg Registry) error
}

// ModuleRegistry keeps a list of modules and initializes them against a
// Registry implementation (the platform's router).
type ModuleRegistry struct {
	modules []Module
}

func NewRegistry() *ModuleRegistry { return &ModuleRegistry{} }

// RegisterModule appends the module to the registry.
func (r *ModuleRegistry) RegisterModule(m Module) { r.modules = append(r.modules, m) }

// InitializeAll calls Register on every module with the provided Registry.
func (r *ModuleRegistry) InitializeAll(reg Registry) error {
	for _, m := range r.modules {
		if err := m.Register(reg); err != nil {
			return err
		}
	}
	return nil
}
