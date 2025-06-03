package entities

import (
	"teammate/server/seedwork/domain"
)

type ProviderType string

const (
	TicketingProvider     ProviderType = "ticketing"
	TranscriptionProvider ProviderType = "transcription"
	MeetingBotProvider    ProviderType = "meeting_bot"
)

// IntegrationConfig represents a user's configuration for external service integrations
// This entity belongs to the User aggregate as users own their integration configurations
type IntegrationConfig struct {
	domain.BaseEntity
	UserID       string                 `json:"user_id" gorm:"column:user_id;not null"`
	ProviderType ProviderType           `json:"provider_type" gorm:"column:provider_type;not null"`
	ProviderName string                 `json:"provider_name" gorm:"column:provider_name;not null"`
	Config       map[string]interface{} `json:"config" gorm:"column:config;type:jsonb;not null"`
	IsActive     bool                   `json:"is_active" gorm:"column:is_active;default:true"`
}

// NewIntegrationConfig creates a new IntegrationConfig entity
func NewIntegrationConfig(userID string, providerType ProviderType, providerName string, config map[string]interface{}) IntegrationConfig {
	if config == nil {
		config = make(map[string]interface{})
	}

	integrationConfig := IntegrationConfig{
		UserID:       userID,
		ProviderType: providerType,
		ProviderName: providerName,
		Config:       config,
		IsActive:     true,
	}
	integrationConfig.SetID(domain.GenerateID())
	return integrationConfig
}

// Activate enables the integration configuration
func (ic *IntegrationConfig) Activate() {
	ic.IsActive = true
}

// Deactivate disables the integration configuration
func (ic *IntegrationConfig) Deactivate() {
	ic.IsActive = false
}

// UpdateConfig updates the configuration settings
func (ic *IntegrationConfig) UpdateConfig(config map[string]interface{}) {
	if config != nil {
		ic.Config = config
	}
}

// GetConfigValue returns a configuration value by key
func (ic *IntegrationConfig) GetConfigValue(key string) (interface{}, bool) {
	if ic.Config == nil {
		return nil, false
	}
	value, exists := ic.Config[key]
	return value, exists
}

// SetConfigValue sets a configuration value
func (ic *IntegrationConfig) SetConfigValue(key string, value interface{}) {
	if ic.Config == nil {
		ic.Config = make(map[string]interface{})
	}
	ic.Config[key] = value
}

// RemoveConfigValue removes a configuration value
func (ic *IntegrationConfig) RemoveConfigValue(key string) {
	if ic.Config != nil {
		delete(ic.Config, key)
	}
}

// HasConfigValue checks if a configuration value exists
func (ic *IntegrationConfig) HasConfigValue(key string) bool {
	if ic.Config == nil {
		return false
	}
	_, exists := ic.Config[key]
	return exists
}

// IsTicketingProvider returns true if this is a ticketing provider configuration
func (ic *IntegrationConfig) IsTicketingProvider() bool {
	return ic.ProviderType == TicketingProvider
}

// IsTranscriptionProvider returns true if this is a transcription provider configuration
func (ic *IntegrationConfig) IsTranscriptionProvider() bool {
	return ic.ProviderType == TranscriptionProvider
}

// IsMeetingBotProvider returns true if this is a meeting bot provider configuration
func (ic *IntegrationConfig) IsMeetingBotProvider() bool {
	return ic.ProviderType == MeetingBotProvider
}

// GetUniqueKey returns a unique key for this configuration (user + provider type + provider name)
func (ic *IntegrationConfig) GetUniqueKey() string {
	return ic.UserID + ":" + string(ic.ProviderType) + ":" + ic.ProviderName
}

// Validate performs basic validation on the integration configuration
func (ic *IntegrationConfig) Validate() error {
	if ic.UserID == "" {
		return domain.NewDomainError("INVALID_USER_ID", "User ID is required", nil)
	}

	if ic.ProviderType == "" {
		return domain.NewDomainError("INVALID_PROVIDER_TYPE", "Provider type is required", nil)
	}

	if ic.ProviderName == "" {
		return domain.NewDomainError("INVALID_PROVIDER_NAME", "Provider name is required", nil)
	}

	// Validate provider type
	switch ic.ProviderType {
	case TicketingProvider, TranscriptionProvider, MeetingBotProvider:
		// Valid provider types
	default:
		return domain.NewDomainError("INVALID_PROVIDER_TYPE", "Invalid provider type", nil)
	}

	return nil
}

// TableName sets the table name for GORM
func (IntegrationConfig) TableName() string {
	return "integration_configs"
}
