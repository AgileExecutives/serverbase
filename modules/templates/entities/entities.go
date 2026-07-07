package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Channel represents the delivery channel for a template
type Channel string

const (
	ChannelEmail Channel = "EMAIL"
	ChannelSMS   Channel = "SMS"
)

// Template represents a stored template metadata record
type Template struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TenantID       uint           `json:"tenant_id" gorm:"index"`
	OrganizationID *uint          `json:"organization_id"`
	Module         string         `json:"module"`
	TemplateKey    string         `json:"template_key" gorm:"index"`
	Channel        Channel        `json:"channel"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	StorageKey     string         `json:"storage_key"`
	Version        int            `json:"version"`
	IsActive       bool           `json:"is_active"`
	IsDefault      bool           `json:"is_default"`
	Variables      datatypes.JSON `json:"variables"`
	SampleData     datatypes.JSON `json:"sample_data"`
	Subject        *string        `json:"subject"`
	TemplateType   string         `json:"template_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TemplateContract represents a contract for a template (variables + sample data)
type TemplateContract struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Module            string         `json:"module" gorm:"index"`
	TemplateKey       string         `json:"template_key" gorm:"index"`
	VariableSchema    datatypes.JSON `json:"variable_schema"`
	DefaultSampleData datatypes.JSON `json:"default_sample_data"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (Template) TableName() string         { return "templates" }
func (TemplateContract) TableName() string { return "template_contracts" }
