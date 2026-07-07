package models

import internalmodels "github.com/AgileExecutives/serverbase/internal/models"

// Public aliases to internal models so external modules (shared-modules)
// can import stable types without referencing internal packages.

type User = internalmodels.User
type UserResponse = internalmodels.UserResponse
type UserCreateRequest = internalmodels.UserCreateRequest
type UserUpdateRequest = internalmodels.UserUpdateRequest
type LoginRequest = internalmodels.LoginRequest
type LoginResponse = internalmodels.LoginResponse
type Tenant = internalmodels.Tenant
type TenantResponse = internalmodels.TenantResponse
type Email = internalmodels.Email
type EmailResponse = internalmodels.EmailResponse
type EmailSendRequest = internalmodels.EmailSendRequest
type Contact = internalmodels.Contact
type ContactResponse = internalmodels.ContactResponse
type ContactCreateRequest = internalmodels.ContactCreateRequest
type ContactUpdateRequest = internalmodels.ContactUpdateRequest
type ContactFormRequest = internalmodels.ContactFormRequest
type ContactFormResponse = internalmodels.ContactFormResponse
type UserSettings = internalmodels.UserSettings
type UserSettingsResponse = internalmodels.UserSettingsResponse
type TokenBlacklist = internalmodels.TokenBlacklist
type Organization = internalmodels.Organization
type Plan = internalmodels.Plan
