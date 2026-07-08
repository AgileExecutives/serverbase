// Package swagger provides runtime aggregation of per-module Swagger/OpenAPI 2.0
// documentation. Each module that has pre-generated docs registers its rendered
// JSON string during Initialize; the Application merges everything into a single
// spec and serves it at /swagger/index.html via gin-swagger.
package swagger

import "sync"

// Registry implements core.DocRegistry. It is safe for concurrent use.
type Registry struct {
	mu   sync.RWMutex
	docs map[string]string
}

// NewRegistry returns a ready-to-use Registry.
func NewRegistry() *Registry {
	return &Registry{docs: make(map[string]string)}
}

// RegisterDoc stores a fully-rendered (no template variables) Swagger 2.0 JSON
// document produced by a module. Calling again with the same name overwrites.
func (r *Registry) RegisterDoc(moduleName, docJSON string) {
	r.mu.Lock()
	r.docs[moduleName] = docJSON
	r.mu.Unlock()
}

// Docs returns a snapshot of all registered documents keyed by module name.
func (r *Registry) Docs() map[string]string {
	r.mu.RLock()
	out := make(map[string]string, len(r.docs))
	for k, v := range r.docs {
		out[k] = v
	}
	r.mu.RUnlock()
	return out
}
