package module

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
)

// AdapterModule is a small helper to adapt internal packages into a core.Module
// by composing entities, routes, services and lifecycle hooks. It is intended
// for internal packages that don't need a full module implementation but want
// to expose a `NewModule()` compatible factory for the module registry.
type AdapterModule struct {
	name         string
	version      string
	dependencies []string

	entities      []core.Entity
	routes        []core.RouteProvider
	services      []core.ServiceProvider
	eventHandlers []core.EventHandler
	middleware    []core.MiddlewareProvider
	swaggerPaths  []string

	initFn  func(ctx core.ModuleContext) error
	startFn func(ctx context.Context) error
	stopFn  func(ctx context.Context) error
}

// Option configures an AdapterModule
type Option func(*AdapterModule)

// NewAdapterModule creates a new AdapterModule with required metadata and options.
func NewAdapterModule(name, version string, deps []string, opts ...Option) core.Module {
	m := &AdapterModule{name: name, version: version, dependencies: deps}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// WithEntities attaches entities to the module
func WithEntities(es ...core.Entity) Option {
	return func(m *AdapterModule) { m.entities = append(m.entities, es...) }
}

// WithRoutes attaches route providers
func WithRoutes(rps ...core.RouteProvider) Option {
	return func(m *AdapterModule) { m.routes = append(m.routes, rps...) }
}

// WithServices attaches service providers
func WithServices(sps ...core.ServiceProvider) Option {
	return func(m *AdapterModule) { m.services = append(m.services, sps...) }
}

// WithEventHandlers attaches event handlers
func WithEventHandlers(handlers ...core.EventHandler) Option {
	return func(m *AdapterModule) { m.eventHandlers = append(m.eventHandlers, handlers...) }
}

// WithMiddleware attaches middleware providers
func WithMiddleware(mws ...core.MiddlewareProvider) Option {
	return func(m *AdapterModule) { m.middleware = append(m.middleware, mws...) }
}

// WithSwaggerPaths attaches swagger doc paths
func WithSwaggerPaths(paths ...string) Option {
	return func(m *AdapterModule) { m.swaggerPaths = append(m.swaggerPaths, paths...) }
}

// WithInit sets an Initialize hook
func WithInit(fn func(ctx core.ModuleContext) error) Option {
	return func(m *AdapterModule) { m.initFn = fn }
}

// WithStart sets a Start hook
func WithStart(fn func(ctx context.Context) error) Option {
	return func(m *AdapterModule) { m.startFn = fn }
}

// WithStop sets a Stop hook
func WithStop(fn func(ctx context.Context) error) Option {
	return func(m *AdapterModule) { m.stopFn = fn }
}

// core.Module implementation
func (m *AdapterModule) Name() string           { return m.name }
func (m *AdapterModule) Version() string        { return m.version }
func (m *AdapterModule) Dependencies() []string { return append([]string{}, m.dependencies...) }

func (m *AdapterModule) Initialize(ctx core.ModuleContext) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *AdapterModule) Start(ctx context.Context) error {
	if m.startFn != nil {
		return m.startFn(ctx)
	}
	return nil
}

func (m *AdapterModule) Stop(ctx context.Context) error {
	if m.stopFn != nil {
		return m.stopFn(ctx)
	}
	return nil
}

func (m *AdapterModule) Entities() []core.Entity { return append([]core.Entity{}, m.entities...) }
func (m *AdapterModule) Routes() []core.RouteProvider {
	return append([]core.RouteProvider{}, m.routes...)
}
func (m *AdapterModule) Services() []core.ServiceProvider {
	return append([]core.ServiceProvider{}, m.services...)
}
func (m *AdapterModule) EventHandlers() []core.EventHandler {
	return append([]core.EventHandler{}, m.eventHandlers...)
}
func (m *AdapterModule) Middleware() []core.MiddlewareProvider {
	return append([]core.MiddlewareProvider{}, m.middleware...)
}
func (m *AdapterModule) SwaggerPaths() []string { return append([]string{}, m.swaggerPaths...) }
