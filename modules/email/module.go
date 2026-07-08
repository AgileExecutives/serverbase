package email

import (
	"context"

	"github.com/AgileExecutives/serverbase/modules/email/entities"
	"github.com/AgileExecutives/serverbase/modules/email/events"
	"github.com/AgileExecutives/serverbase/modules/email/handlers"
	"github.com/AgileExecutives/serverbase/modules/email/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
)

// EmailModule represents the email module
type EmailModule struct {
	emailEntity   *entities.EmailEntity
	emailHandler  *handlers.EmailHandler
	publicHandler *handlers.PublicEmailRoutes
	emailService  *services.EmailService
	eventHandlers []core.EventHandler
}

// NewEmailModule creates a new email module instance
func NewEmailModule() *EmailModule {
	return &EmailModule{}
}

// Name returns the module name
func (m *EmailModule) Name() string { return "email" }

// Version returns the module version
func (m *EmailModule) Version() string { return "1.0.0" }

// Description returns the module description
func (m *EmailModule) Description() string { return "Email management and notification system" }

// Dependencies returns the module dependencies
func (m *EmailModule) Dependencies() []string { return []string{} }

// Initialize initializes the email module
func (m *EmailModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing Email module...")

	m.emailEntity = entities.NewEmailEntity()
	m.emailService = services.NewEmailService()
	m.emailHandler = handlers.NewEmailHandler(ctx.DB, m.emailService)
	m.publicHandler = handlers.NewPublicEmailRoutes(m.emailHandler)

	m.eventHandlers = []core.EventHandler{
		events.NewEmailSentHandler(ctx.DB),
		events.NewEmailFailedHandler(ctx.DB, ctx.Logger),
	}

	if ctx.DocRegistry != nil {
		ctx.DocRegistry.RegisterDoc(m.Name(), EmailSwaggerJSON)
	}

	ctx.Logger.Info("Email module initialized successfully")
	return nil
}

// Start starts the email module
func (m *EmailModule) Start(ctx context.Context) error { return nil }

// Stop stops the email module
func (m *EmailModule) Stop(ctx context.Context) error { return nil }

// Entities returns core entities
func (m *EmailModule) Entities() []core.Entity { return []core.Entity{m.emailEntity} }

// Routes returns route providers
func (m *EmailModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{m.emailHandler, m.publicHandler}
}

// EventHandlers returns event handlers
func (m *EmailModule) EventHandlers() []core.EventHandler { return m.eventHandlers }

// Services returns module services
func (m *EmailModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{
		&EmailServiceProvider{module: m},
	}
}

// Middleware returns middleware providers
func (m *EmailModule) Middleware() []core.MiddlewareProvider { return []core.MiddlewareProvider{} }

// SwaggerPaths returns swagger paths
func (m *EmailModule) SwaggerPaths() []string { return []string{} }

// EmailServiceProvider provides the email service
type EmailServiceProvider struct{ module *EmailModule }

func (p *EmailServiceProvider) ServiceName() string           { return "email-service" }
func (p *EmailServiceProvider) ServiceInterface() interface{} { return p.module.emailService }
func (p *EmailServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	if p.module.emailService == nil {
		p.module.emailService = services.NewEmailService()
	}
	return p.module.emailService, nil
}
