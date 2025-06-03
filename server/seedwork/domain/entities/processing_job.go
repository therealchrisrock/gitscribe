package entities

import (
	"teammate/server/seedwork/domain"
	"time"
)

// ProcessingJob represents an asynchronous processing task
// This is a general-purpose entity for managing background jobs across all modules
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

// Job types - these can be used across all modules
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
