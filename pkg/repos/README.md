GORM Repo Factory

Purpose:
- Centralize creation of GORM-backed repository adapters so bootstrapping code can request typed repositories from a single place.

Usage:
- Create a factory bound to your `*gorm.DB`:

  rf := repos.NewGormRepoFactory(db)

- Obtain repositories:

  tenantRepo := rf.TenantRepo()
  customerRepo := rf.CustomerRepo()
  userRepo := rf.UserRepo()

Notes:
- Repositories returned by the factory implement the module-local `repo` interfaces defined in each module (for example, `modules/customers/repo.CustomerRepo`).
- Use the factory in bootstrap or application wiring code to keep repo construction in one place and simplify tests/mocking.

Example (bootstrap):

  rf := repos.NewGormRepoFactory(ctx.DB)
  tenantSvc := services.NewTenantService(rf.TenantRepo(), ...)

Testing:
- For unit tests, prefer the module `inmemory` repo implementations (e.g. `modules/customers/repo.NewInMemoryCustomerRepo()`) instead of the GORM factory.
