package module

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"gorm.io/gorm"
)

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

// RegisterCoreModules migrates entities for each module, then adapts and
// registers them with the registry. This is the standard one-liner for wiring
// core.Modules into a server-test or production harness:
//
//	module.RegisterCoreModules(mr, modules, db, coreCtx)
func RegisterCoreModules(mr *ModuleRegistry, modules []core.Module, db *gorm.DB, ctx core.ModuleContext) error {
	for _, m := range modules {
		for _, e := range m.Entities() {
			if model := e.GetModel(); model != nil {
				if err := db.AutoMigrate(model); err != nil {
					return err
				}
			}
		}
		mr.RegisterModule(newCoreAdapter(m, ctx))
	}
	return nil
}

// coreModuleAdapter adapts a core.Module to the serverbase Module interface.
type coreModuleAdapter struct {
	mod core.Module
	ctx core.ModuleContext
}

func newCoreAdapter(m core.Module, ctx core.ModuleContext) Module {
	return &coreModuleAdapter{mod: m, ctx: ctx}
}

func (a *coreModuleAdapter) Name() string { return a.mod.Name() }

// Register initializes the underlying core.Module and registers its routes on
// the gin engine under /api/v1.
func (a *coreModuleAdapter) Register(reg Registry) error {
	if err := a.mod.Initialize(a.ctx); err != nil {
		return err
	}
	apiV1 := a.ctx.Router.Group("/api/v1")
	for _, rp := range a.mod.Routes() {
		group := apiV1.Group(rp.GetPrefix())
		for _, mw := range rp.GetMiddleware() {
			group.Use(mw)
		}
		rp.RegisterRoutes(group, a.ctx)
	}
	return nil
}
