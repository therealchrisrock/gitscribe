# Domain Model Documentation

## Entity Relationship Diagram

```mermaid
erDiagram
    User {
        string id PK
        string name
        string email
        timestamp created_at
        timestamp updated_at
    }
    
    Meeting {
        string id PK
        string user_id FK
        string title
        string type "zoom|google_meet|microsoft_teams|generic"
        string status "scheduled|in_progress|completed|failed"
        timestamp start_time
        timestamp end_time
        string meeting_url
        string bot_join_url
        string recording_path
        timestamp created_at
        timestamp updated_at
    }
    
    Participant {
        string id PK
        string meeting_id FK
        string name
        string email
        string role
        timestamp created_at
    }
    
    BotSession {
        string id PK
        string meeting_id FK
        string session_id
        string status "joining|active|recording|completed|failed"
        timestamp joined_at
        timestamp left_at
        string bot_user_id
        json metadata
        timestamp created_at
    }
    
    Transcription {
        string id PK
        string meeting_id FK
        string audio_file_path
        string status "pending|processing|completed|failed"
        text content
        float confidence
        string provider "assembly_ai|zoom_api|google_speech"
        timestamp created_at
        timestamp updated_at
    }
    
    TranscriptSegment {
        string id PK
        string transcription_id FK
        string speaker
        text text
        float start_time
        float end_time
        float confidence
        int sequence_number
        timestamp created_at
    }
    
    ActionItem {
        string id PK
        string meeting_id FK
        string transcription_id FK
        string title
        text description
        string assignee
        string priority "low|medium|high|urgent"
        timestamp due_date
        string status "extracted|pending|approved|created|rejected"
        text context
        timestamp created_at
        timestamp updated_at
    }
    
    TicketReference {
        string id PK
        string action_item_id FK
        string system "github|jira|linear|asana"
        string ticket_id
        string ticket_url
        string project_key
        string reference_type "existing|created"
        json metadata
        timestamp created_at
    }
    
    IntegrationConfig {
        string id PK
        string user_id FK
        string provider_type "ticketing|transcription|meeting_bot"
        string provider_name
        json config
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }
    
    ProcessingJob {
        string id PK
        string entity_type "meeting|transcription|action_item"
        string entity_id
        string job_type "transcribe|extract_actions|create_tickets"
        string status "pending|processing|completed|failed"
        json payload
        text error_message
        int retry_count
        timestamp scheduled_at
        timestamp started_at
        timestamp completed_at
        timestamp created_at
    }

    %% Relationships
    User ||--o{ Meeting : "creates"
    User ||--o{ IntegrationConfig : "configures"
    
    Meeting ||--o{ Participant : "has"
    Meeting ||--o{ BotSession : "has"
    Meeting ||--o{ Transcription : "generates"
    Meeting ||--o{ ActionItem : "produces"
    
    Transcription ||--o{ TranscriptSegment : "contains"
    Transcription ||--o{ ActionItem : "generates"
    
    ActionItem ||--o{ TicketReference : "references"
    
    Meeting ||--o{ ProcessingJob : "triggers"
    Transcription ||--o{ ProcessingJob : "triggers"
    ActionItem ||--o{ ProcessingJob : "triggers"
```

## Domain Entities

### User
The root entity representing a system user.

**Attributes:**
- `id`: Unique identifier
- `name`: User's display name
- `email`: User's email address (unique)
- `created_at`, `updated_at`: Audit timestamps

**Relationships:**
- One-to-many with Meetings
- One-to-many with IntegrationConfigs

### Meeting
Core entity representing a meeting session.

**Attributes:**
- `id`: Unique identifier
- `user_id`: Reference to the meeting creator
- `title`: Meeting title/subject
- `type`: Meeting platform (zoom, google_meet, microsoft_teams, generic)
- `status`: Current meeting state
- `start_time`, `end_time`: Meeting duration
- `meeting_url`: Original meeting URL
- `bot_join_url`: Bot-specific join URL (if different)
- `recording_path`: Path to recorded audio/video file

**Status Values:**
- `scheduled`: Meeting created but not started
- `in_progress`: Meeting is active
- `completed`: Meeting finished successfully
- `failed`: Meeting processing failed

### Participant
Represents individuals who attended the meeting.

**Attributes:**
- `meeting_id`: Reference to the meeting
- `name`: Participant's name
- `email`: Participant's email (optional)
- `role`: Participant's role in the meeting (optional)

### BotSession
Tracks the bot's participation in a meeting.

**Attributes:**
- `meeting_id`: Reference to the meeting
- `session_id`: External bot session identifier
- `status`: Bot's current state
- `joined_at`, `left_at`: Bot participation timeframe
- `bot_user_id`: Bot's user ID in the meeting platform
- `metadata`: Additional bot-specific data

### Transcription
Contains the transcribed content of a meeting.

**Attributes:**
- `meeting_id`: Reference to the source meeting
- `audio_file_path`: Path to the audio file
- `status`: Processing status
- `content`: Full transcribed text
- `confidence`: Overall transcription confidence score
- `provider`: Transcription service used

### TranscriptSegment
Individual segments of the transcription with speaker attribution.

**Attributes:**
- `transcription_id`: Reference to the parent transcription
- `speaker`: Speaker identification
- `text`: Segment content
- `start_time`, `end_time`: Temporal boundaries
- `confidence`: Segment-specific confidence score
- `sequence_number`: Order within the transcription

### ActionItem
Extracted actionable items from the meeting.

**Attributes:**
- `meeting_id`: Reference to the source meeting
- `transcription_id`: Reference to the source transcription
- `title`: Action item summary
- `description`: Detailed description
- `assignee`: Person responsible (optional)
- `priority`: Urgency level
- `due_date`: Target completion date (optional)
- `status`: Processing status
- `context`: Original transcript context

**Priority Values:**
- `low`, `medium`, `high`, `urgent`

**Status Values:**
- `extracted`: Identified from transcript
- `pending`: Awaiting user approval
- `approved`: User approved for ticket creation
- `created`: Ticket successfully created
- `rejected`: User rejected the action item

### TicketReference
Links action items to external ticketing systems.

**Attributes:**
- `action_item_id`: Reference to the action item
- `system`: Ticketing platform identifier
- `ticket_id`: External ticket identifier
- `ticket_url`: Direct link to the ticket
- `project_key`: Project/repository identifier
- `reference_type`: Whether ticket was existing or newly created
- `metadata`: System-specific additional data

### IntegrationConfig
User's configuration for external service integrations.

**Attributes:**
- `user_id`: Reference to the user
- `provider_type`: Type of integration (ticketing, transcription, meeting_bot)
- `provider_name`: Specific provider (github, jira, assembly_ai, etc.)
- `config`: Provider-specific configuration (API keys, endpoints, etc.)
- `is_active`: Whether this configuration is currently active

### ProcessingJob
Tracks asynchronous processing tasks.

**Attributes:**
- `entity_type`: Type of entity being processed
- `entity_id`: Identifier of the entity
- `job_type`: Type of processing task
- `status`: Current job status
- `payload`: Job-specific data
- `error_message`: Error details if failed
- `retry_count`: Number of retry attempts
- `scheduled_at`, `started_at`, `completed_at`: Job lifecycle timestamps

## Value Objects

### Email (User Module)
Encapsulates email validation and formatting.

### Priority (ActionItem Module)
Represents action item priority with validation.

### MeetingType (Meeting Module)
Enumeration of supported meeting platforms.

### TranscriptionProvider (Transcription Module)
Identifies the transcription service used.

## Domain Services

### MeetingOrchestrator
Coordinates the entire meeting processing workflow.

### TranscriptionService
Manages transcription provider selection and processing.

### ActionItemExtractionService
Handles AI-powered action item extraction from transcripts.

### TicketingService
Manages ticket creation and linking across multiple platforms.

## Repository Interfaces

Each domain entity has a corresponding repository interface:

- `UserRepository`
- `MeetingRepository`
- `TranscriptionRepository`
- `ActionItemRepository`
- `TicketReferenceRepository`
- `IntegrationConfigRepository`
- `ProcessingJobRepository`

All repositories extend the base `Repository[T]` interface with entity-specific query methods. 