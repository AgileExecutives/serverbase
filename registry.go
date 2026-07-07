package serverbase

import "github.com/gin-gonic/gin"

// ModuleRegistry holds modules to be registered on a server.
type ModuleRegistry struct {
	modules []Module
}

// NewModuleRegistry creates a new registry.
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{modules: make([]Module, 0)}
}

// Register adds a module to the registry.
func (mr *ModuleRegistry) Register(m Module) {
	mr.modules = append(mr.modules, m)
}

// RegisterAll registers all modules on the given router group.
func (mr *ModuleRegistry) RegisterAll(rg *gin.RouterGroup) error {
	for _, m := range mr.modules {
		m.RegisterRoutes(rg)
	}
	return nil
}

// Modules returns registered modules (read-only slice).
func (mr *ModuleRegistry) Modules() []Module { return mr.modules }
