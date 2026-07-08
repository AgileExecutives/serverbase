package module

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"gorm.io/gorm"
)

// LifecycleManager provides ordered lifecycle operations for core.Modules.
// It centralizes AutoMigrate, Initialize, Start and Stop phases so bootstrapping
// code can call them in a deterministic way.
type LifecycleManager struct {
	modules []core.Module
	db      *gorm.DB
	ctx     core.ModuleContext
}

func NewLifecycleManager(mods []core.Module, db *gorm.DB, ctx core.ModuleContext) *LifecycleManager {
	return &LifecycleManager{modules: mods, db: db, ctx: ctx}
}

// AutoMigrateAll runs AutoMigrate for all module entities in deterministic order.
func (m *LifecycleManager) AutoMigrateAll() error {
	if m.db == nil {
		return nil
	}
	for _, mod := range m.modules {
		for _, e := range mod.Entities() {
			if model := e.GetModel(); model != nil {
				if err := m.db.AutoMigrate(model); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// InitializeAll calls Initialize on every module. Modules should register docs
// and service providers during Initialize.
func (m *LifecycleManager) InitializeAll() error {
	for _, mod := range m.modules {
		if err := mod.Initialize(m.ctx); err != nil {
			return err
		}
	}
	return nil
}

// StartAll calls Start on every module.
func (m *LifecycleManager) StartAll(ctx context.Context) error {
	for _, mod := range m.modules {
		if err := mod.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// StopAll calls Stop on every module (reverse order is intentional).
func (m *LifecycleManager) StopAll(ctx context.Context) error {
	for i := len(m.modules) - 1; i >= 0; i-- {
		if err := m.modules[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
