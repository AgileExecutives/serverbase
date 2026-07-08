package user

import (
	"context"

	basedocs "github.com/AgileExecutives/serverbase/modules/base/docs"
	"github.com/AgileExecutives/serverbase/modules/user/entities"
	"github.com/AgileExecutives/serverbase/modules/user/events"
	"github.com/AgileExecutives/serverbase/modules/user/handlers"
	"github.com/AgileExecutives/serverbase/modules/user/middleware"
	"github.com/AgileExecutives/serverbase/modules/user/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	settingsentities "github.com/AgileExecutives/serverbase/pkg/settings/entities"
)

// UserModule provides core authentication, user management, and contact functionality
type UserModule struct {
	authHandlers         *handlers.AuthHandlers
	contactHandlers      *handlers.ContactHandlers
	healthHandlers       *handlers.HealthHandlers
	userSettingsHandlers *handlers.UserSettingsHandlers
	authService          *services.AuthService
	eventHandlers        *events.BaseEventHandlers
	authMiddleware       *middleware.AuthMiddleware
	moduleContext        core.ModuleContext
}

// NewUserModule creates a new user module instance
func NewUserModule() core.Module {
	return &UserModule{}
}

func (m *UserModule) Name() string {
	return "user"
}

func (m *UserModule) Version() string {
	return "1.0.0"
}

func (m *UserModule) Dependencies() []string {
	return []string{}
}

func (m *UserModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing user module...")
	m.moduleContext = ctx
	m.authService = services.NewAuthService(ctx.DB, ctx.Logger)
	m.authHandlers = handlers.NewAuthHandlers(ctx.DB, ctx.Logger)
	m.contactHandlers = handlers.NewContactHandlers(ctx.DB, ctx.Logger)
	m.healthHandlers = handlers.NewHealthHandlers(ctx.DB, ctx.Logger)
	m.userSettingsHandlers = handlers.NewUserSettingsHandlers(ctx.DB, ctx.Logger)
	m.eventHandlers = events.NewBaseEventHandlers(ctx.EventBus, ctx.Logger)
	m.authMiddleware = middleware.NewAuthMiddleware(ctx.DB, ctx.Logger)

	// Register pre-generated swagger docs (same handler set as the base module).
	if ctx.DocRegistry != nil {
		ctx.DocRegistry.RegisterDoc(m.Name(), basedocs.SwaggerInfobase.ReadDoc())
	}

	ctx.Logger.Info("User module initialized successfully")
	return nil
}

func (m *UserModule) Start(ctx context.Context) error {
	if m.authHandlers != nil && m.moduleContext.ModuleRegistry != nil {
		m.authHandlers.SetModuleRegistry(m.moduleContext.ModuleRegistry)
		m.moduleContext.Logger.Info("User module started - auth handlers configured with module registry for template seeding")
	} else {
		m.moduleContext.Logger.Warn("User module started - module registry not available for template seeding")
	}
	return nil
}

func (m *UserModule) Stop(ctx context.Context) error {
	return nil
}

func (m *UserModule) Entities() []core.Entity {
	return []core.Entity{
		entities.NewUserEntity(),
		entities.NewTenantEntity(),
		entities.NewContactEntity(),
		entities.NewNewsletterEntity(),
		entities.NewTokenBlacklistEntity(),
		entities.NewUserSettingsEntity(),
		settingsentities.NewSettingDefinitionEntity(),
		settingsentities.NewSettingEntity(),
	}
}

func (m *UserModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{
		handlers.NewAuthRoutes(m.authHandlers),
		handlers.NewContactRoutes(m.contactHandlers),
		handlers.NewUserSettingsRoutes(m.userSettingsHandlers),
		handlers.NewHealthRoutes(m.healthHandlers),
	}
}

func (m *UserModule) EventHandlers() []core.EventHandler {
	return []core.EventHandler{
		events.NewUserCreatedHandler(m.eventHandlers),
		events.NewUserLoginHandler(m.eventHandlers),
		events.NewContactFormSubmittedHandler(m.eventHandlers),
	}
}

func (m *UserModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{
		services.NewAuthServiceProvider(m.authService),
	}
}

func (m *UserModule) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{
		middleware.NewAuthMiddlewareProvider(m.authMiddleware),
	}
}

func (m *UserModule) SwaggerPaths() []string {
	return []string{
		"./modules/user/handlers",
		"./modules/user/entities",
	}
}
