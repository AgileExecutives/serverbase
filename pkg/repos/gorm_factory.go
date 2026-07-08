package repos

import (
	"github.com/AgileExecutives/serverbase/modules/customers/repo"
	tenantrepo "github.com/AgileExecutives/serverbase/modules/tenant/repo"
	userrepo "github.com/AgileExecutives/serverbase/modules/user/repo"
	"gorm.io/gorm"
)

// GormRepoFactory centralizes creation of GORM-backed repo implementations.
type GormRepoFactory struct {
	db *gorm.DB
}

// NewGormRepoFactory returns a factory bound to the provided DB.
func NewGormRepoFactory(db *gorm.DB) *GormRepoFactory { return &GormRepoFactory{db: db} }

func (f *GormRepoFactory) TenantRepo() tenantrepo.TenantRepo {
	return tenantrepo.NewGormTenantRepo(f.db)
}
func (f *GormRepoFactory) CustomerRepo() repo.CustomerRepo { return repo.NewGormCustomerRepo(f.db) }
func (f *GormRepoFactory) UserRepo() userrepo.UserRepo     { return userrepo.NewGormUserRepo(f.db) }
