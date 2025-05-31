# System Architecture Documentation

## Module Architecture Overview

```mermaid
graph TB
    subgraph "Chrome Extension (WXT)"
        CE[Content Scripts]
        CP[Popup Interface]
        CB[Background Service]
        CE --> CP
        CP --> CB
    end
    
    subgraph "API Gateway"
        AG[HTTP Router]
        AM[Auth Middleware]
        AG --> AM
    end
    
    subgraph "Meeting Module"
        subgraph "Meeting Domain"
            ME[Meeting Entity]
            MP[Participant Entity]
            MBS[Bot Session Entity]
            MR[Meeting Repository]
        end
        
        subgraph "Meeting Application"
            MOC[Meeting Orchestrator]
            MCH[Meeting Command Handlers]
            MQH[Meeting Query Handlers]
        end
        
        subgraph "Meeting Infrastructure"
            MRI[Meeting Repo Implementation]
            MBP[Meeting Bot Providers]
        end
    end
    
    subgraph "Transcription Module"
        subgraph "Transcription Domain"
            TE[Transcription Entity]
            TS[Transcript Segment Entity]
            TR[Transcription Repository]
            TSS[Transcription Service]
        end
        
        subgraph "Transcription Application"
            TCH[Transcription Command Handlers]
            TQH[Transcription Query Handlers]
        end
        
        subgraph "Transcription Infrastructure"
            TRI[Transcription Repo Implementation]
            AAI[Assembly AI Provider]
            ZAP[Zoom API Provider]
            GSP[Google Speech Provider]
        end
    end
    
    subgraph "Action Item Module"
        subgraph "Action Item Domain"
            AIE[Action Item Entity]
            AIR[Action Item Repository]
            AIS[Action Item Service]
        end
        
        subgraph "Action Item Application"
            AICH[Action Item Command Handlers]
            AIQH[Action Item Query Handlers]
            AIES[Action Item Extraction Service]
        end
        
        subgraph "Action Item Infrastructure"
            AIRI[Action Item Repo Implementation]
            AILP[AI Language Processing]
        end
    end
    
    subgraph "Ticketing Module"
        subgraph "Ticketing Domain"
            TRE[Ticket Reference Entity]
            TRR[Ticket Reference Repository]
            TIS[Ticketing Service]
        end
        
        subgraph "Ticketing Application"
            TICH[Ticketing Command Handlers]
            TIQH[Ticketing Query Handlers]
        end
        
        subgraph "Ticketing Infrastructure"
            TIRI[Ticketing Repo Implementation]
            GHP[GitHub Provider]
            JP[Jira Provider]
            LP[Linear Provider]
        end
    end
    
    subgraph "Integration Module"
        subgraph "Integration Domain"
            ICE[Integration Config Entity]
            ICR[Integration Config Repository]
        end
        
        subgraph "Integration Application"
            ICCH[Integration Command Handlers]
            ICQH[Integration Query Handlers]
        end
        
        subgraph "Integration Infrastructure"
            ICRI[Integration Repo Implementation]
            EXT[External API Clients]
        end
    end
    
    subgraph "Shared Infrastructure"
        DB[(PostgreSQL Database)]
        FS[File Storage]
        MQ[Message Queue]
        LOG[Logging Service]
    end
    
    %% Connections
    CB --> AG
    AG --> MOC
    MOC --> TSS
    MOC --> AIS
    MOC --> TIS
    
    TSS --> AAI
    TSS --> ZAP
    TSS --> GSP
    
    TIS --> GHP
    TIS --> JP
    TIS --> LP
    
    MRI --> DB
    TRI --> DB
    AIRI --> DB
    TIRI --> DB
    ICRI --> DB
    
    MBP --> FS
    AAI --> FS
    
    MOC --> MQ
    TSS --> MQ
    AIS --> MQ
```

## Data Flow Sequence

```mermaid
sequenceDiagram
    participant CE as Chrome Extension
    participant API as API Gateway
    participant MO as Meeting Orchestrator
    participant BS as Bot Service
    participant TS as Transcription Service
    participant AIS as Action Item Service
    participant TIS as Ticketing Service
    participant DB as Database
    participant EXT as External APIs
    
    Note over CE,EXT: Meeting Processing Flow
    
    CE->>API: POST /meetings (meeting URL, config)
    API->>MO: ProcessMeetingCommand
    MO->>DB: Create Meeting entity
    MO->>BS: JoinMeeting(url, config)
    BS->>EXT: Join meeting via bot API
    EXT-->>BS: Bot session created
    BS-->>MO: BotSession details
    MO->>DB: Update meeting with bot session
    MO->>BS: StartRecording(sessionId)
    BS->>EXT: Start recording
    MO-->>API: Meeting started response
    API-->>CE: Success response
    
    Note over CE,EXT: Recording Completion (Async)
    
    EXT->>BS: Recording completed webhook
    BS->>MO: RecordingCompleted event
    MO->>TS: TranscribeAudio(audioPath)
    TS->>EXT: Send to Assembly AI
    EXT-->>TS: Transcription result
    TS->>DB: Save transcription & segments
    TS-->>MO: Transcription completed
    
    MO->>AIS: ExtractActionItems(transcriptionId)
    AIS->>EXT: AI processing for action items
    EXT-->>AIS: Extracted action items
    AIS->>DB: Save action items
    AIS-->>MO: Action items extracted
    
    MO->>TIS: ProcessActionItems(actionItemIds)
    TIS->>EXT: Search existing tickets
    EXT-->>TIS: Existing tickets found
    TIS->>DB: Link to existing tickets
    TIS->>EXT: Create new tickets
    EXT-->>TIS: New tickets created
    TIS->>DB: Save ticket references
    TIS-->>MO: Ticketing completed
    
    MO->>DB: Update meeting status to completed
    MO->>CE: Send completion notification
```

## Provider Strategy Pattern

```mermaid
classDiagram
    class TranscriptionProvider {
        <<interface>>
        +Transcribe(audioPath string) Transcription
        +GetProviderName() string
        +SupportsRealTime() bool
        +GetSupportedFormats() []string
    }
    
    class AssemblyAIProvider {
        -apiKey string
        -client HTTPClient
        +Transcribe(audioPath string) Transcription
        +GetProviderName() string
        +SupportsRealTime() bool
        +GetSupportedFormats() []string
    }
    
    class ZoomAPIProvider {
        -apiKey string
        -apiSecret string
        +Transcribe(audioPath string) Transcription
        +GetProviderName() string
        +SupportsRealTime() bool
        +GetSupportedFormats() []string
    }
    
    class GoogleSpeechProvider {
        -credentials ServiceAccount
        +Transcribe(audioPath string) Transcription
        +GetProviderName() string
        +SupportsRealTime() bool
        +GetSupportedFormats() []string
    }
    
    class TranscriptionService {
        -providers map[string]TranscriptionProvider
        -defaultProvider string
        +RegisterProvider(provider TranscriptionProvider)
        +Transcribe(meetingType string, audioPath string) Transcription
        +SelectProvider(meetingType string) TranscriptionProvider
    }
    
    class TicketingProvider {
        <<interface>>
        +GetProviderName() string
        +SearchExistingTickets(query string) []ExistingTicket
        +CreateTicket(request CreateTicketRequest) CreatedTicket
        +UpdateTicket(ticketId string, updates UpdateTicketRequest) UpdatedTicket
        +GetTicketURL(ticketId string) string
    }
    
    class GitHubProvider {
        -token string
        -client GitHubClient
        +GetProviderName() string
        +SearchExistingTickets(query string) []ExistingTicket
        +CreateTicket(request CreateTicketRequest) CreatedTicket
        +UpdateTicket(ticketId string, updates UpdateTicketRequest) UpdatedTicket
        +GetTicketURL(ticketId string) string
    }
    
    class JiraProvider {
        -baseURL string
        -username string
        -token string
        +GetProviderName() string
        +SearchExistingTickets(query string) []ExistingTicket
        +CreateTicket(request CreateTicketRequest) CreatedTicket
        +UpdateTicket(ticketId string, updates UpdateTicketRequest) UpdatedTicket
        +GetTicketURL(ticketId string) string
    }
    
    class TicketingService {
        -providers map[string]TicketingProvider
        -userConfigs map[string]ProviderConfig
        +RegisterProvider(provider TicketingProvider)
        +ProcessActionItems(actionItems []ActionItem) []TicketResult
        +GetProviderForUser(userId string) TicketingProvider
    }
    
    TranscriptionProvider <|.. AssemblyAIProvider
    TranscriptionProvider <|.. ZoomAPIProvider
    TranscriptionProvider <|.. GoogleSpeechProvider
    TranscriptionService --> TranscriptionProvider
    
    TicketingProvider <|.. GitHubProvider
    TicketingProvider <|.. JiraProvider
    TicketingService --> TicketingProvider
```

## Meeting State Machine

```mermaid
stateDiagram-v2
    [*] --> Scheduled
    
    Scheduled --> Joining: Start Meeting
    Joining --> InProgress: Bot Joined Successfully
    Joining --> Failed: Bot Join Failed
    
    InProgress --> Recording: Start Recording
    Recording --> Processing: Recording Completed
    InProgress --> Completed: Manual Stop
    InProgress --> Failed: Connection Lost
    
    Processing --> Transcribing: Audio Processing
    Transcribing --> ExtractingActions: Transcription Complete
    Transcribing --> Failed: Transcription Failed
    
    ExtractingActions --> MatchingTickets: Actions Extracted
    ExtractingActions --> Failed: Extraction Failed
    
    MatchingTickets --> CreatingTickets: Matches Found
    MatchingTickets --> Failed: Matching Failed
    
    CreatingTickets --> Completed: Tickets Created
    CreatingTickets --> PartiallyCompleted: Some Tickets Failed
    CreatingTickets --> Failed: All Tickets Failed
    
    Failed --> [*]
    Completed --> [*]
    PartiallyCompleted --> [*]
    
    note right of Scheduled
        Meeting created but not started
    end note
    
    note right of InProgress
        Bot actively in meeting
        Can transition to multiple states
    end note
    
    note right of Processing
        Async processing pipeline
        Multiple stages can fail independently
    end note
```

## Layer Responsibilities

### Domain Layer
- **Entities**: Core business objects with identity and lifecycle
- **Value Objects**: Immutable objects representing concepts
- **Domain Services**: Business logic that doesn't belong to a single entity
- **Repository Interfaces**: Contracts for data access
- **Domain Events**: Notifications of important business occurrences

### Application Layer
- **Command Handlers**: Process write operations and business workflows
- **Query Handlers**: Handle read operations and data retrieval
- **Application Services**: Orchestrate complex business processes
- **DTOs**: Data transfer objects for API communication
- **Validation**: Input validation and business rule enforcement

### Infrastructure Layer
- **Repository Implementations**: Concrete data access implementations
- **External Service Providers**: Integrations with third-party APIs
- **Database Configuration**: ORM setup and migrations
- **Message Queue**: Async processing infrastructure
- **File Storage**: Audio/video file management

### Interface Layer
- **HTTP Handlers**: REST API endpoints
- **Middleware**: Authentication, logging, error handling
- **WebSocket Handlers**: Real-time communication
- **Background Jobs**: Scheduled and async task processing

## Key Architectural Decisions

### 1. Modular Monolith
- Single deployable unit with clear module boundaries
- Easier development and debugging than microservices
- Can be split into microservices later if needed

### 2. Event-Driven Processing
- Async processing for long-running tasks
- Resilient to failures with retry mechanisms
- Scalable processing pipeline

### 3. Provider Pattern for Integrations
- Easy to add new meeting platforms and ticketing systems
- Configuration-driven provider selection
- Testable through interface mocking

### 4. CQRS (Command Query Responsibility Segregation)
- Separate read and write operations
- Optimized query models for different use cases
- Clear separation of concerns

### 5. Repository Pattern
- Abstract data access behind interfaces
- Testable business logic
- Database-agnostic domain layer

## Scalability Considerations

### Horizontal Scaling
- Stateless application servers
- Database connection pooling
- Load balancer distribution

### Async Processing
- Message queue for background tasks
- Worker processes for heavy operations
- Retry mechanisms for failed jobs

### Caching Strategy
- Redis for session management
- Application-level caching for frequent queries
- CDN for static assets

### Database Optimization
- Proper indexing strategy
- Read replicas for query scaling
- Partitioning for large tables 