package entities

import (
	"teammate/server/seedwork/domain"
	"time"
)

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
	Urgent Priority = "urgent"
)

type ActionItemStatus string

const (
	Extracted ActionItemStatus = "extracted"
	Pending   ActionItemStatus = "pending"
	Approved  ActionItemStatus = "approved"
	Created   ActionItemStatus = "created"
	Rejected  ActionItemStatus = "rejected"
)

// ActionItem represents an extracted action item from a meeting
type ActionItem struct {
	domain.BaseEntity
	MeetingID        string            `json:"meeting_id" gorm:"column:meeting_id;not null"`
	TranscriptionID  string            `json:"transcription_id" gorm:"column:transcription_id;not null"`
	Title            string            `json:"title" gorm:"column:title;not null"`
	Description      string            `json:"description" gorm:"column:description;type:text;not null"`
	Assignee         string            `json:"assignee,omitempty" gorm:"column:assignee"`
	Priority         Priority          `json:"priority" gorm:"column:priority;not null"`
	DueDate          *time.Time        `json:"due_date,omitempty" gorm:"column:due_date"`
	Status           ActionItemStatus  `json:"status" gorm:"column:status;not null"`
	Context          string            `json:"context" gorm:"column:context;type:text;not null"`
	TicketReferences []TicketReference `json:"ticket_references" gorm:"foreignKey:ActionItemID"`
}

// NewActionItem creates a new ActionItem entity
func NewActionItem(meetingID, transcriptionID, title, description, context string, priority Priority) ActionItem {
	actionItem := ActionItem{
		MeetingID:       meetingID,
		TranscriptionID: transcriptionID,
		Title:           title,
		Description:     description,
		Priority:        priority,
		Status:          Extracted,
		Context:         context,
	}
	actionItem.SetID(domain.GenerateID())
	return actionItem
}

// Approve transitions the action item to approved status
func (a *ActionItem) Approve() {
	a.Status = Approved
}

// Reject transitions the action item to rejected status
func (a *ActionItem) Reject() {
	a.Status = Rejected
}

// MarkAsCreated transitions the action item to created status
func (a *ActionItem) MarkAsCreated() {
	a.Status = Created
}

// SetPending transitions the action item to pending status
func (a *ActionItem) SetPending() {
	a.Status = Pending
}

// SetAssignee sets the assignee for the action item
func (a *ActionItem) SetAssignee(assignee string) {
	a.Assignee = assignee
}

// SetDueDate sets the due date for the action item
func (a *ActionItem) SetDueDate(dueDate time.Time) {
	a.DueDate = &dueDate
}

// ClearDueDate removes the due date from the action item
func (a *ActionItem) ClearDueDate() {
	a.DueDate = nil
}

// IsApproved returns true if the action item has been approved
func (a *ActionItem) IsApproved() bool {
	return a.Status == Approved
}

// IsCreated returns true if the action item has been created as a ticket
func (a *ActionItem) IsCreated() bool {
	return a.Status == Created
}

// IsPending returns true if the action item is pending approval
func (a *ActionItem) IsPending() bool {
	return a.Status == Pending
}

// IsRejected returns true if the action item has been rejected
func (a *ActionItem) IsRejected() bool {
	return a.Status == Rejected
}

// HasAssignee returns true if the action item has an assignee
func (a *ActionItem) HasAssignee() bool {
	return a.Assignee != ""
}

// HasDueDate returns true if the action item has a due date
func (a *ActionItem) HasDueDate() bool {
	return a.DueDate != nil
}

// IsOverdue returns true if the action item is past its due date
func (a *ActionItem) IsOverdue() bool {
	if a.DueDate == nil {
		return false
	}
	return time.Now().After(*a.DueDate)
}

// GetPriorityLevel returns a numeric representation of priority for sorting
func (a *ActionItem) GetPriorityLevel() int {
	switch a.Priority {
	case Urgent:
		return 4
	case High:
		return 3
	case Medium:
		return 2
	case Low:
		return 1
	default:
		return 0
	}
}

// AddTicketReference adds a ticket reference to the action item
func (a *ActionItem) AddTicketReference(system, ticketID, ticketURL, projectKey, referenceType string, metadata map[string]interface{}) {
	ticketRef := NewTicketReference(a.GetID(), system, ticketID, ticketURL, projectKey, referenceType, metadata)
	a.TicketReferences = append(a.TicketReferences, ticketRef)
}

// HasTicketReferences returns true if the action item has any ticket references
func (a *ActionItem) HasTicketReferences() bool {
	return len(a.TicketReferences) > 0
}

// GetTicketReferencesBySystem returns ticket references for a specific system
func (a *ActionItem) GetTicketReferencesBySystem(system string) []TicketReference {
	var refs []TicketReference
	for _, ref := range a.TicketReferences {
		if ref.System == system {
			refs = append(refs, ref)
		}
	}
	return refs
}

// TableName sets the table name for GORM
func (ActionItem) TableName() string {
	return "action_items"
}

// TicketReference represents a reference to an external ticket
type TicketReference struct {
	domain.BaseEntity
	ActionItemID  string                 `json:"action_item_id" gorm:"column:action_item_id;not null"`
	System        string                 `json:"system" gorm:"column:system;not null"`
	TicketID      string                 `json:"ticket_id" gorm:"column:ticket_id;not null"`
	TicketURL     string                 `json:"ticket_url" gorm:"column:ticket_url;not null"`
	ProjectKey    string                 `json:"project_key,omitempty" gorm:"column:project_key"`
	ReferenceType string                 `json:"reference_type" gorm:"column:reference_type;not null"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

const (
	ExistingTicketReference = "existing"
	CreatedTicketReference  = "created"
)

// NewTicketReference creates a new TicketReference entity
func NewTicketReference(actionItemID, system, ticketID, ticketURL, projectKey, referenceType string, metadata map[string]interface{}) TicketReference {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	ticketRef := TicketReference{
		ActionItemID:  actionItemID,
		System:        system,
		TicketID:      ticketID,
		TicketURL:     ticketURL,
		ProjectKey:    projectKey,
		ReferenceType: referenceType,
		Metadata:      metadata,
	}
	ticketRef.SetID(domain.GenerateID())
	return ticketRef
}

// IsExistingTicket returns true if this references an existing ticket
func (tr *TicketReference) IsExistingTicket() bool {
	return tr.ReferenceType == ExistingTicketReference
}

// IsCreatedTicket returns true if this references a newly created ticket
func (tr *TicketReference) IsCreatedTicket() bool {
	return tr.ReferenceType == CreatedTicketReference
}

// GetMetadataValue returns a metadata value by key
func (tr *TicketReference) GetMetadataValue(key string) (interface{}, bool) {
	if tr.Metadata == nil {
		return nil, false
	}
	value, exists := tr.Metadata[key]
	return value, exists
}

// SetMetadataValue sets a metadata value
func (tr *TicketReference) SetMetadataValue(key string, value interface{}) {
	if tr.Metadata == nil {
		tr.Metadata = make(map[string]interface{})
	}
	tr.Metadata[key] = value
}

// TableName sets the table name for GORM
func (TicketReference) TableName() string {
	return "ticket_references"
}
