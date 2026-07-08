package services

import (
	"context"

	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
)

// ContactService provides operations for contacts and newsletters.
type ContactService struct {
	repo   repo.ContactRepo
	logger core.Logger
}

func NewContactServiceWithRepo(r repo.ContactRepo, logger core.Logger) *ContactService {
	return &ContactService{repo: r, logger: logger}
}

func (s *ContactService) ListContacts(ctx context.Context, offset, limit int, active *bool, contactType string) ([]models.Contact, int64, error) {
	return s.repo.ListContacts(ctx, offset, limit, active, contactType)
}

func (s *ContactService) GetContact(ctx context.Context, id uint) (*models.Contact, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ContactService) CreateContact(ctx context.Context, c *models.Contact) error {
	return s.repo.CreateContact(ctx, c)
}

func (s *ContactService) UpdateContact(ctx context.Context, c *models.Contact) error {
	return s.repo.UpdateContact(ctx, c)
}

func (s *ContactService) DeleteContact(ctx context.Context, c *models.Contact) error {
	return s.repo.DeleteContact(ctx, c)
}

func (s *ContactService) UpsertNewsletter(ctx context.Context, n *basemodels.Newsletter) (bool, error) {
	return s.repo.UpsertNewsletter(ctx, n)
}

func (s *ContactService) ListNewsletters(ctx context.Context) ([]basemodels.Newsletter, error) {
	return s.repo.ListNewsletters(ctx)
}

func (s *ContactService) DeleteNewsletterByEmail(ctx context.Context, email string) (int64, error) {
	return s.repo.DeleteNewsletterByEmail(ctx, email)
}

type ContactServiceProvider struct{ service *ContactService }

func NewContactServiceProvider(service *ContactService) core.ServiceProvider {
	return &ContactServiceProvider{service: service}
}

func (p *ContactServiceProvider) ServiceName() string           { return "contact" }
func (p *ContactServiceProvider) ServiceInterface() interface{} { return (*ContactService)(nil) }
func (p *ContactServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	return p.service, nil
}
