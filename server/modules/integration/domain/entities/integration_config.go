package entities

import (
	"teammate/server/seedwork/domain"
	"time"
)

type ProviderType string

const (
	TicketingProvider     ProviderType = "ticketing"
	TranscriptionProvider ProviderType = "transcription"
	MeetingBotProvider    ProviderType = "meeting_bot"
)

// IntegrationConfig represents a user's configuration for external service integrations
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

// ProcessingJob represents an asynchronous processing task
type ProcessingJob struct {
	domain.BaseEntity
	EntityType   string                 `json:"entity_type" gorm:"column:entity_type;not null"`
	EntityID     string                 `json:"entity_id" gorm:"column:entity_id;not null"`
	JobType      string                 `json:"job_type" gorm:"column:job_type;not null"`
	Status       ProcessingJobStatus    `json:"status" gorm:"column:status;not null"`
	Payload      map[string]interface{} `json:"payload" gorm:"column:payload;type:jsonb;not null"`
	ErrorMessage string                 `json:"error_message,omitempty" gorm:"column:error_message;type:text"`
	RetryCount   int                    `json:"retry_count" gorm:"column:retry_count;default:0"`
	ScheduledAt  *time.Time             `json:"scheduled_at" gorm:"column:scheduled_at;not null"`
	StartedAt    *time.Time             `json:"started_at,omitempty" gorm:"column:started_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty" gorm:"column:completed_at"`
}

type ProcessingJobStatus string

const (
	JobPending    ProcessingJobStatus = "pending"
	JobProcessing ProcessingJobStatus = "processing"
	JobCompleted  ProcessingJobStatus = "completed"
	JobFailed     ProcessingJobStatus = "failed"
)

// Job types
const (
	TranscribeJobType     = "transcribe"
	ExtractActionsJobType = "extract_actions"
	CreateTicketsJobType  = "create_tickets"
	ProcessMeetingJobType = "process_meeting"
)

// NewProcessingJob creates a new ProcessingJob entity
func NewProcessingJob(entityType, entityID, jobType string, payload map[string]interface{}) ProcessingJob {
	if payload == nil {
		payload = make(map[string]interface{})
	}

	now := time.Now()
	job := ProcessingJob{
		EntityType:  entityType,
		EntityID:    entityID,
		JobType:     jobType,
		Status:      JobPending,
		Payload:     payload,
		RetryCount:  0,
		ScheduledAt: &now,
	}
	job.SetID(domain.GenerateID())
	return job
}

// Start transitions the job to processing status
func (pj *ProcessingJob) Start() {
	pj.Status = JobProcessing
	now := time.Now()
	pj.StartedAt = &now
}

// Complete transitions the job to completed status
func (pj *ProcessingJob) Complete() {
	pj.Status = JobCompleted
	now := time.Now()
	pj.CompletedAt = &now
}

// Fail transitions the job to failed status
func (pj *ProcessingJob) Fail(errorMessage string) {
	pj.Status = JobFailed
	pj.ErrorMessage = errorMessage
	now := time.Now()
	pj.CompletedAt = &now
}

// Retry increments the retry count and resets the job to pending
func (pj *ProcessingJob) Retry() {
	pj.RetryCount++
	pj.Status = JobPending
	pj.ErrorMessage = ""
	pj.StartedAt = nil
	pj.CompletedAt = nil
}

// CanRetry returns true if the job can be retried (max 3 retries)
func (pj *ProcessingJob) CanRetry() bool {
	return pj.RetryCount < 3
}

// IsCompleted returns true if the job has completed successfully
func (pj *ProcessingJob) IsCompleted() bool {
	return pj.Status == JobCompleted
}

// IsFailed returns true if the job has failed
func (pj *ProcessingJob) IsFailed() bool {
	return pj.Status == JobFailed
}

// IsProcessing returns true if the job is currently being processed
func (pj *ProcessingJob) IsProcessing() bool {
	return pj.Status == JobProcessing
}

// IsPending returns true if the job is pending execution
func (pj *ProcessingJob) IsPending() bool {
	return pj.Status == JobPending
}

// GetPayloadValue returns a payload value by key
func (pj *ProcessingJob) GetPayloadValue(key string) (interface{}, bool) {
	if pj.Payload == nil {
		return nil, false
	}
	value, exists := pj.Payload[key]
	return value, exists
}

// SetPayloadValue sets a payload value
func (pj *ProcessingJob) SetPayloadValue(key string, value interface{}) {
	if pj.Payload == nil {
		pj.Payload = make(map[string]interface{})
	}
	pj.Payload[key] = value
}

// TableName sets the table name for GORM
func (ProcessingJob) TableName() string {
	return "processing_jobs"
}
