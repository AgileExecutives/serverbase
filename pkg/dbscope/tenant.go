package dbscope

import "gorm.io/gorm"

// TenantScope returns a GORM scope that filters by tenant_id.
func TenantScope(tenantID uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}
