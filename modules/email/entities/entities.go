package entities

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
)

type EmailEntity struct{}

func NewEmailEntity() *EmailEntity                     { return &EmailEntity{} }
func (e *EmailEntity) Name() string                    { return "email" }
func (e *EmailEntity) TableName() string               { return "emails" }
func (e *EmailEntity) GetModel() interface{}           { return &models.Email{} }
func (e *EmailEntity) GetMigrations() []core.Migration { return []core.Migration{&EmailMigration001{}} }

type EmailMigration001 struct{}

func (m *EmailMigration001) Version() string        { return "001_create_emails_table" }
func (m *EmailMigration001) Up(db *gorm.DB) error   { return db.AutoMigrate(&models.Email{}) }
func (m *EmailMigration001) Down(db *gorm.DB) error { return db.Migrator().DropTable(&models.Email{}) }
