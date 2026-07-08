package module

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/pkg/core"
)

type fakeModule struct {
	name        string
	initialized bool
	started     bool
	stopped     bool
}

func (f *fakeModule) Name() string                            { return f.name }
func (f *fakeModule) Version() string                         { return "0.0.1" }
func (f *fakeModule) Dependencies() []string                  { return nil }
func (f *fakeModule) Initialize(ctx core.ModuleContext) error { f.initialized = true; return nil }
func (f *fakeModule) Start(ctx context.Context) error         { f.started = true; return nil }
func (f *fakeModule) Stop(ctx context.Context) error          { f.stopped = true; return nil }
func (f *fakeModule) Entities() []core.Entity                 { return nil }
func (f *fakeModule) Routes() []core.RouteProvider            { return nil }
func (f *fakeModule) EventHandlers() []core.EventHandler      { return nil }
func (f *fakeModule) Middleware() []core.MiddlewareProvider   { return nil }
func (f *fakeModule) Services() []core.ServiceProvider        { return nil }
func (f *fakeModule) SwaggerPaths() []string                  { return nil }

func TestLifecycleManager_BasicFlow(t *testing.T) {
	fm := &fakeModule{name: "fm"}
	lm := NewLifecycleManager([]core.Module{fm}, nil, core.ModuleContext{})

	if err := lm.InitializeAll(); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if !fm.initialized {
		t.Fatalf("module not initialized")
	}

	if err := lm.StartAll(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !fm.started {
		t.Fatalf("module not started")
	}

	if err := lm.StopAll(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if !fm.stopped {
		t.Fatalf("module not stopped")
	}
}
