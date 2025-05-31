# Implementation Guide

## Domain Entities Implementation

### Meeting Module

```go
// server/modules/meeting/domain/entities/meeting.go
package entities

import (
    "time"
    "teammate/server/seedwork/domain"
)

type MeetingType string
const (
    ZoomMeeting      MeetingType = "zoom"
    GoogleMeet       MeetingType = "google_meet"
    MicrosoftTeams   MeetingType = "microsoft_teams"
    GenericMeeting   MeetingType = "generic"
)

type MeetingStatus string
const (
    Scheduled  MeetingStatus = "scheduled"
    InProgress MeetingStatus = "in_progress"
    Completed  MeetingStatus = "completed"
    Failed     MeetingStatus = "failed"
)

type Meeting struct {
    domain.BaseEntity
    UserID        string        `json:"user_id" gorm:"column:user_id;not null"`
    Title         string        `json:"title" gorm:"column:title;not null"`
    Type          MeetingType   `json:"type" gorm:"column:type;not null"`
    Status        MeetingStatus `json:"status" gorm:"column:status;not null"`
    StartTime     time.Time     `json:"start_time" gorm:"column:start_time;not null"`
    EndTime       *time.Time    `json:"end_time,omitempty" gorm:"column:end_time"`
    MeetingURL    string        `json:"meeting_url" gorm:"column:meeting_url;not null"`
    BotJoinURL    string        `json:"bot_join_url,omitempty" gorm:"column:bot_join_url"`
    RecordingPath string        `json:"recording_path,omitempty" gorm:"column:recording_path"`
    Participants  []Participant `json:"participants" gorm:"foreignKey:MeetingID"`
}

func NewMeeting(userID, title string, meetingType MeetingType, meetingURL string) Meeting {
    meeting := Meeting{
        UserID:     userID,
        Title:      title,
        Type:       meetingType,
        Status:     Scheduled,
        StartTime:  time.Now(),
        MeetingURL: meetingURL,
    }
    meeting.SetID(generateID())
    return meeting
}

func (m *Meeting) StartMeeting(botJoinURL string) {
    m.Status = InProgress
    m.BotJoinURL = botJoinURL
}

func (m *Meeting) CompleteMeeting(recordingPath string) {
    m.Status = Completed
    m.RecordingPath = recordingPath
    now := time.Now()
    m.EndTime = &now
}

func (m *Meeting) FailMeeting() {
    m.Status = Failed
    now := time.Now()
    m.EndTime = &now
}

func (Meeting) TableName() string {
    return "meetings"
}

type Participant struct {
    domain.BaseEntity
    MeetingID string `json:"meeting_id" gorm:"column:meeting_id;not null"`
    Name      string `json:"name" gorm:"column:name;not null"`
    Email     string `json:"email,omitempty" gorm:"column:email"`
    Role      string `json:"role,omitempty" gorm:"column:role"`
}

func (Participant) TableName() string {
    return "participants"
}
```

### Transcription Module

```go
// server/modules/transcription/domain/entities/transcription.go
package entities

import (
    "teammate/server/seedwork/domain"
)

type TranscriptionStatus string
const (
    Pending    TranscriptionStatus = "pending"
    Processing TranscriptionStatus = "processing"
    Completed  TranscriptionStatus = "completed"
    Failed     TranscriptionStatus = "failed"
)

type Transcription struct {
    domain.BaseEntity
    MeetingID     string              `json:"meeting_id" gorm:"column:meeting_id;not null"`
    AudioFilePath string              `json:"audio_file_path" gorm:"column:audio_file_path;not null"`
    Status        TranscriptionStatus `json:"status" gorm:"column:status;not null"`
    Content       string              `json:"content" gorm:"column:content;type:text"`
    Confidence    float64             `json:"confidence" gorm:"column:confidence"`
    Provider      string              `json:"provider" gorm:"column:provider;not null"`
    Segments      []TranscriptSegment `json:"segments" gorm:"foreignKey:TranscriptionID"`
}

func NewTranscription(meetingID, audioFilePath, provider string) Transcription {
    transcription := Transcription{
        MeetingID:     meetingID,
        AudioFilePath: audioFilePath,
        Status:        Pending,
        Provider:      provider,
    }
    transcription.SetID(generateID())
    return transcription
}

func (t *Transcription) StartProcessing() {
    t.Status = Processing
}

func (t *Transcription) CompleteTranscription(content string, confidence float64, segments []TranscriptSegment) {
    t.Status = Completed
    t.Content = content
    t.Confidence = confidence
    t.Segments = segments
}

func (t *Transcription) FailTranscription() {
    t.Status = Failed
}

func (Transcription) TableName() string {
    return "transcriptions"
}

type TranscriptSegment struct {
    domain.BaseEntity
    TranscriptionID string  `json:"transcription_id" gorm:"column:transcription_id;not null"`
    Speaker         string  `json:"speaker" gorm:"column:speaker"`
    Text            string  `json:"text" gorm:"column:text;type:text;not null"`
    StartTime       float64 `json:"start_time" gorm:"column:start_time;not null"`
    EndTime         float64 `json:"end_time" gorm:"column:end_time;not null"`
    Confidence      float64 `json:"confidence" gorm:"column:confidence"`
    SequenceNumber  int     `json:"sequence_number" gorm:"column:sequence_number;not null"`
}

func (TranscriptSegment) TableName() string {
    return "transcript_segments"
}
```

### Action Item Module

```go
// server/modules/actionitem/domain/entities/action_item.go
package entities

import (
    "time"
    "teammate/server/seedwork/domain"
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

type ActionItem struct {
    domain.BaseEntity
    MeetingID         string           `json:"meeting_id" gorm:"column:meeting_id;not null"`
    TranscriptionID   string           `json:"transcription_id" gorm:"column:transcription_id;not null"`
    Title             string           `json:"title" gorm:"column:title;not null"`
    Description       string           `json:"description" gorm:"column:description;type:text;not null"`
    Assignee          string           `json:"assignee,omitempty" gorm:"column:assignee"`
    Priority          Priority         `json:"priority" gorm:"column:priority;not null"`
    DueDate           *time.Time       `json:"due_date,omitempty" gorm:"column:due_date"`
    Status            ActionItemStatus `json:"status" gorm:"column:status;not null"`
    Context           string           `json:"context" gorm:"column:context;type:text;not null"`
    TicketReferences  []TicketReference `json:"ticket_references" gorm:"foreignKey:ActionItemID"`
}

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
    actionItem.SetID(generateID())
    return actionItem
}

func (a *ActionItem) Approve() {
    a.Status = Approved
}

func (a *ActionItem) Reject() {
    a.Status = Rejected
}

func (a *ActionItem) MarkAsCreated() {
    a.Status = Created
}

func (a *ActionItem) SetAssignee(assignee string) {
    a.Assignee = assignee
}

func (a *ActionItem) SetDueDate(dueDate time.Time) {
    a.DueDate = &dueDate
}

func (ActionItem) TableName() string {
    return "action_items"
}

type TicketReference struct {
    domain.BaseEntity
    ActionItemID    string `json:"action_item_id" gorm:"column:action_item_id;not null"`
    System          string `json:"system" gorm:"column:system;not null"`
    TicketID        string `json:"ticket_id" gorm:"column:ticket_id;not null"`
    TicketURL       string `json:"ticket_url" gorm:"column:ticket_url;not null"`
    ProjectKey      string `json:"project_key,omitempty" gorm:"column:project_key"`
    ReferenceType   string `json:"reference_type" gorm:"column:reference_type;not null"`
    Metadata        string `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

func (TicketReference) TableName() string {
    return "ticket_references"
}
```

## Application Services Implementation

### Meeting Orchestrator

```go
// server/modules/meeting/application/services/meeting_orchestrator.go
package services

import (
    "context"
    "fmt"
    
    meetingEntities "teammate/server/modules/meeting/domain/entities"
    meetingRepos "teammate/server/modules/meeting/domain/repositories"
    transcriptionServices "teammate/server/modules/transcription/application/services"
    actionitemServices "teammate/server/modules/actionitem/application/services"
    ticketingServices "teammate/server/modules/ticketing/application/services"
    integrationServices "teammate/server/modules/integration/application/services"
)

type ProcessMeetingCommand struct {
    UserID      string                    `json:"user_id"`
    Title       string                    `json:"title"`
    MeetingURL  string                    `json:"meeting_url"`
    MeetingType meetingEntities.MeetingType `json:"meeting_type"`
    BotConfig   BotConfig                 `json:"bot_config"`
}

type BotConfig struct {
    BotName        string `json:"bot_name"`
    RecordAudio    bool   `json:"record_audio"`
    RecordVideo    bool   `json:"record_video"`
    AutoTranscribe bool   `json:"auto_transcribe"`
}

type MeetingOrchestrator struct {
    meetingRepo       meetingRepos.MeetingRepository
    botService        integrationServices.MeetingBotService
    transcriptionSvc  transcriptionServices.TranscriptionService
    actionItemSvc     actionitemServices.ActionItemService
    ticketingSvc      ticketingServices.TicketingService
    eventPublisher    EventPublisher
}

func NewMeetingOrchestrator(
    meetingRepo meetingRepos.MeetingRepository,
    botService integrationServices.MeetingBotService,
    transcriptionSvc transcriptionServices.TranscriptionService,
    actionItemSvc actionitemServices.ActionItemService,
    ticketingSvc ticketingServices.TicketingService,
    eventPublisher EventPublisher,
) *MeetingOrchestrator {
    return &MeetingOrchestrator{
        meetingRepo:      meetingRepo,
        botService:       botService,
        transcriptionSvc: transcriptionSvc,
        actionItemSvc:    actionItemSvc,
        ticketingSvc:     ticketingSvc,
        eventPublisher:   eventPublisher,
    }
}

func (o *MeetingOrchestrator) ProcessMeeting(ctx context.Context, cmd ProcessMeetingCommand) (*meetingEntities.Meeting, error) {
    // 1. Create meeting entity
    meeting := meetingEntities.NewMeeting(cmd.UserID, cmd.Title, cmd.MeetingType, cmd.MeetingURL)
    
    if err := o.meetingRepo.Create(&meeting); err != nil {
        return nil, fmt.Errorf("failed to create meeting: %w", err)
    }
    
    // 2. Join meeting with bot
    botSession, err := o.botService.JoinMeeting(ctx, cmd.MeetingURL, cmd.BotConfig)
    if err != nil {
        meeting.FailMeeting()
        o.meetingRepo.Update(&meeting)
        return nil, fmt.Errorf("failed to join meeting with bot: %w", err)
    }
    
    // 3. Start recording
    err = o.botService.StartRecording(ctx, botSession.SessionID)
    if err != nil {
        meeting.FailMeeting()
        o.meetingRepo.Update(&meeting)
        return nil, fmt.Errorf("failed to start recording: %w", err)
    }
    
    // 4. Update meeting status
    meeting.StartMeeting(botSession.BotJoinURL)
    if err := o.meetingRepo.Update(&meeting); err != nil {
        return nil, fmt.Errorf("failed to update meeting: %w", err)
    }
    
    // 5. Publish event for async processing
    event := MeetingStartedEvent{
        MeetingID: meeting.GetID(),
        SessionID: botSession.SessionID,
        UserID:    cmd.UserID,
    }
    o.eventPublisher.Publish("meeting.started", event)
    
    return &meeting, nil
}

func (o *MeetingOrchestrator) HandleMeetingCompletion(ctx context.Context, event MeetingCompletedEvent) error {
    // 1. Get meeting
    meeting, err := o.meetingRepo.FindByID(event.MeetingID)
    if err != nil {
        return fmt.Errorf("failed to find meeting: %w", err)
    }
    
    // 2. Update meeting with recording path
    meeting.CompleteMeeting(event.RecordingPath)
    if err := o.meetingRepo.Update(meeting); err != nil {
        return fmt.Errorf("failed to update meeting: %w", err)
    }
    
    // 3. Start transcription
    transcription, err := o.transcriptionSvc.StartTranscription(ctx, StartTranscriptionCommand{
        MeetingID:     event.MeetingID,
        AudioFilePath: event.RecordingPath,
        Provider:      "assembly_ai", // Could be configurable
    })
    if err != nil {
        return fmt.Errorf("failed to start transcription: %w", err)
    }
    
    // 4. Publish transcription started event
    transcriptionEvent := TranscriptionStartedEvent{
        TranscriptionID: transcription.GetID(),
        MeetingID:       event.MeetingID,
    }
    o.eventPublisher.Publish("transcription.started", transcriptionEvent)
    
    return nil
}

func (o *MeetingOrchestrator) HandleTranscriptionCompletion(ctx context.Context, event TranscriptionCompletedEvent) error {
    // 1. Extract action items
    actionItems, err := o.actionItemSvc.ExtractActionItems(ctx, ExtractActionItemsCommand{
        TranscriptionID: event.TranscriptionID,
        MeetingID:       event.MeetingID,
    })
    if err != nil {
        return fmt.Errorf("failed to extract action items: %w", err)
    }
    
    // 2. Process action items for ticketing
    for _, actionItem := range actionItems {
        err := o.ticketingSvc.ProcessActionItem(ctx, ProcessActionItemCommand{
            ActionItemID: actionItem.GetID(),
            UserID:       event.UserID,
        })
        if err != nil {
            // Log error but continue processing other items
            fmt.Printf("Failed to process action item %s: %v\n", actionItem.GetID(), err)
        }
    }
    
    // 3. Publish completion event
    completionEvent := MeetingProcessingCompletedEvent{
        MeetingID:       event.MeetingID,
        TranscriptionID: event.TranscriptionID,
        ActionItemCount: len(actionItems),
    }
    o.eventPublisher.Publish("meeting.processing.completed", completionEvent)
    
    return nil
}
```

## Provider Implementations

### Assembly AI Transcription Provider

```go
// server/modules/transcription/infrastructure/providers/assembly_ai_provider.go
package providers

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "teammate/server/modules/transcription/domain/entities"
    "teammate/server/modules/transcription/domain/services"
)

type AssemblyAIProvider struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewAssemblyAIProvider(apiKey string) *AssemblyAIProvider {
    return &AssemblyAIProvider{
        apiKey:  apiKey,
        baseURL: "https://api.assemblyai.com/v2",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (p *AssemblyAIProvider) GetProviderName() string {
    return "assembly_ai"
}

func (p *AssemblyAIProvider) SupportsRealTime() bool {
    return true
}

func (p *AssemblyAIProvider) GetSupportedFormats() []string {
    return []string{"mp3", "mp4", "wav", "flac", "m4a"}
}

func (p *AssemblyAIProvider) Transcribe(ctx context.Context, audioFilePath string) (*entities.Transcription, error) {
    // 1. Upload audio file
    uploadURL, err := p.uploadAudioFile(ctx, audioFilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to upload audio file: %w", err)
    }
    
    // 2. Submit transcription request
    transcriptID, err := p.submitTranscriptionRequest(ctx, uploadURL)
    if err != nil {
        return nil, fmt.Errorf("failed to submit transcription request: %w", err)
    }
    
    // 3. Poll for completion
    result, err := p.pollForCompletion(ctx, transcriptID)
    if err != nil {
        return nil, fmt.Errorf("failed to get transcription result: %w", err)
    }
    
    // 4. Convert to domain entity
    return p.convertToTranscription(result, audioFilePath), nil
}

type UploadResponse struct {
    UploadURL string `json:"upload_url"`
}

type TranscriptionRequest struct {
    AudioURL           string `json:"audio_url"`
    SpeakerLabels      bool   `json:"speaker_labels"`
    AutoHighlights     bool   `json:"auto_highlights"`
    SentimentAnalysis  bool   `json:"sentiment_analysis"`
    EntityDetection    bool   `json:"entity_detection"`
}

type TranscriptionResponse struct {
    ID         string  `json:"id"`
    Status     string  `json:"status"`
    Text       string  `json:"text"`
    Confidence float64 `json:"confidence"`
    Words      []Word  `json:"words"`
    Utterances []Utterance `json:"utterances"`
}

type Word struct {
    Text       string  `json:"text"`
    Start      int     `json:"start"`
    End        int     `json:"end"`
    Confidence float64 `json:"confidence"`
    Speaker    string  `json:"speaker"`
}

type Utterance struct {
    Text       string  `json:"text"`
    Start      int     `json:"start"`
    End        int     `json:"end"`
    Confidence float64 `json:"confidence"`
    Speaker    string  `json:"speaker"`
}

func (p *AssemblyAIProvider) uploadAudioFile(ctx context.Context, filePath string) (string, error) {
    // Implementation for uploading audio file to Assembly AI
    // Returns upload URL for the file
    
    req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/upload", nil)
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", p.apiKey)
    req.Header.Set("Content-Type", "application/octet-stream")
    
    // Add file upload logic here
    
    resp, err := p.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var uploadResp UploadResponse
    if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
        return "", err
    }
    
    return uploadResp.UploadURL, nil
}

func (p *AssemblyAIProvider) submitTranscriptionRequest(ctx context.Context, audioURL string) (string, error) {
    reqBody := TranscriptionRequest{
        AudioURL:          audioURL,
        SpeakerLabels:     true,
        AutoHighlights:    true,
        SentimentAnalysis: false,
        EntityDetection:   false,
    }
    
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return "", err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/transcript", bytes.NewBuffer(jsonBody))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", p.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var transcriptResp TranscriptionResponse
    if err := json.NewDecoder(resp.Body).Decode(&transcriptResp); err != nil {
        return "", err
    }
    
    return transcriptResp.ID, nil
}

func (p *AssemblyAIProvider) pollForCompletion(ctx context.Context, transcriptID string) (*TranscriptionResponse, error) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-ticker.C:
            req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/transcript/"+transcriptID, nil)
            if err != nil {
                return nil, err
            }
            
            req.Header.Set("Authorization", p.apiKey)
            
            resp, err := p.httpClient.Do(req)
            if err != nil {
                return nil, err
            }
            
            var transcriptResp TranscriptionResponse
            if err := json.NewDecoder(resp.Body).Decode(&transcriptResp); err != nil {
                resp.Body.Close()
                return nil, err
            }
            resp.Body.Close()
            
            switch transcriptResp.Status {
            case "completed":
                return &transcriptResp, nil
            case "error":
                return nil, fmt.Errorf("transcription failed")
            case "processing", "queued":
                // Continue polling
                continue
            default:
                return nil, fmt.Errorf("unknown status: %s", transcriptResp.Status)
            }
        }
    }
}

func (p *AssemblyAIProvider) convertToTranscription(result *TranscriptionResponse, audioFilePath string) *entities.Transcription {
    transcription := entities.Transcription{
        AudioFilePath: audioFilePath,
        Status:        entities.Completed,
        Content:       result.Text,
        Confidence:    result.Confidence,
        Provider:      p.GetProviderName(),
    }
    
    // Convert utterances to segments
    segments := make([]entities.TranscriptSegment, len(result.Utterances))
    for i, utterance := range result.Utterances {
        segments[i] = entities.TranscriptSegment{
            Speaker:        utterance.Speaker,
            Text:           utterance.Text,
            StartTime:      float64(utterance.Start) / 1000.0, // Convert ms to seconds
            EndTime:        float64(utterance.End) / 1000.0,
            Confidence:     utterance.Confidence,
            SequenceNumber: i + 1,
        }
    }
    
    transcription.Segments = segments
    return &transcription
}
```

### GitHub Ticketing Provider

```go
// server/modules/ticketing/infrastructure/providers/github_provider.go
package providers

import (
    "context"
    "fmt"
    "strings"
    
    "github.com/google/go-github/v45/github"
    "golang.org/x/oauth2"
    
    "teammate/server/modules/ticketing/domain/services"
)

type GitHubProvider struct {
    client *github.Client
    token  string
}

func NewGitHubProvider(token string) *GitHubProvider {
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(context.Background(), ts)
    
    return &GitHubProvider{
        client: github.NewClient(tc),
        token:  token,
    }
}

func (g *GitHubProvider) GetProviderName() string {
    return "github"
}

func (g *GitHubProvider) SearchExistingTickets(ctx context.Context, query string, projectKey string) ([]services.ExistingTicket, error) {
    // Parse project key (format: "owner/repo")
    parts := strings.Split(projectKey, "/")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid project key format, expected 'owner/repo'")
    }
    owner, repo := parts[0], parts[1]
    
    // Search issues
    searchQuery := fmt.Sprintf("repo:%s %s", projectKey, query)
    opts := &github.SearchOptions{
        ListOptions: github.ListOptions{PerPage: 10},
    }
    
    result, _, err := g.client.Search.Issues(ctx, searchQuery, opts)
    if err != nil {
        return nil, fmt.Errorf("failed to search GitHub issues: %w", err)
    }
    
    tickets := make([]services.ExistingTicket, len(result.Issues))
    for i, issue := range result.Issues {
        labels := make([]string, len(issue.Labels))
        for j, label := range issue.Labels {
            labels[j] = *label.Name
        }
        
        assignee := ""
        if issue.Assignee != nil {
            assignee = *issue.Assignee.Login
        }
        
        tickets[i] = services.ExistingTicket{
            ID:          fmt.Sprintf("%d", *issue.Number),
            Title:       *issue.Title,
            Description: getStringValue(issue.Body),
            Status:      *issue.State,
            Assignee:    assignee,
            Labels:      labels,
            Metadata: map[string]string{
                "url":    *issue.HTMLURL,
                "number": fmt.Sprintf("%d", *issue.Number),
                "owner":  owner,
                "repo":   repo,
            },
        }
    }
    
    return tickets, nil
}

func (g *GitHubProvider) CreateTicket(ctx context.Context, request services.CreateTicketRequest) (*services.CreatedTicket, error) {
    // Parse project key
    parts := strings.Split(request.ProjectKey, "/")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid project key format, expected 'owner/repo'")
    }
    owner, repo := parts[0], parts[1]
    
    // Prepare issue request
    issueRequest := &github.IssueRequest{
        Title: &request.Title,
        Body:  &request.Description,
    }
    
    if request.Assignee != "" {
        issueRequest.Assignees = &[]string{request.Assignee}
    }
    
    if len(request.Labels) > 0 {
        issueRequest.Labels = &request.Labels
    }
    
    // Create issue
    issue, _, err := g.client.Issues.Create(ctx, owner, repo, issueRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to create GitHub issue: %w", err)
    }
    
    return &services.CreatedTicket{
        ID:        fmt.Sprintf("%d", *issue.Number),
        URL:       *issue.HTMLURL,
        Title:     *issue.Title,
        Status:    *issue.State,
        Metadata: map[string]string{
            "number": fmt.Sprintf("%d", *issue.Number),
            "owner":  owner,
            "repo":   repo,
        },
    }, nil
}

func (g *GitHubProvider) UpdateTicket(ctx context.Context, ticketID string, updates services.UpdateTicketRequest) (*services.UpdatedTicket, error) {
    // Parse project key from metadata
    parts := strings.Split(updates.ProjectKey, "/")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid project key format, expected 'owner/repo'")
    }
    owner, repo := parts[0], parts[1]
    
    // Convert ticket ID to issue number
    issueNumber := parseInt(ticketID)
    
    // Prepare update request
    issueRequest := &github.IssueRequest{}
    
    if updates.Title != "" {
        issueRequest.Title = &updates.Title
    }
    
    if updates.Description != "" {
        issueRequest.Body = &updates.Description
    }
    
    if updates.Status != "" {
        issueRequest.State = &updates.Status
    }
    
    if updates.Assignee != "" {
        issueRequest.Assignees = &[]string{updates.Assignee}
    }
    
    if len(updates.Labels) > 0 {
        issueRequest.Labels = &updates.Labels
    }
    
    // Update issue
    issue, _, err := g.client.Issues.Edit(ctx, owner, repo, issueNumber, issueRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to update GitHub issue: %w", err)
    }
    
    return &services.UpdatedTicket{
        ID:     fmt.Sprintf("%d", *issue.Number),
        URL:    *issue.HTMLURL,
        Title:  *issue.Title,
        Status: *issue.State,
    }, nil
}

func (g *GitHubProvider) GetTicketURL(ticketID string, projectKey string) string {
    return fmt.Sprintf("https://github.com/%s/issues/%s", projectKey, ticketID)
}

func getStringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func parseInt(s string) int {
    // Implementation to convert string to int
    // Handle error appropriately
    return 0
}
```

## Chrome Extension Implementation

### Content Script for Meeting Detection

```typescript
// webext/entrypoints/content/meeting-detector.ts
export interface MeetingInfo {
  type: MeetingType | null;
  url: string;
  meetingId: string;
  title?: string;
  participants?: string[];
}

export type MeetingType = 'zoom' | 'google_meet' | 'microsoft_teams' | 'generic';

export class MeetingDetector {
  private static readonly MEETING_PATTERNS = {
    zoom: {
      urlPattern: /zoom\.us\/j\/(\d+)/,
      titleSelector: '.meeting-topic',
      participantSelector: '.participants-item'
    },
    google_meet: {
      urlPattern: /meet\.google\.com\/([a-z-]+)/,
      titleSelector: '[data-meeting-title]',
      participantSelector: '[data-participant-id]'
    },
    microsoft_teams: {
      urlPattern: /teams\.microsoft\.com\/l\/meetup-join/,
      titleSelector: '.ts-calling-thread-header',
      participantSelector: '.calling-roster-item'
    }
  };

  static detectMeetingType(url: string): MeetingType | null {
    for (const [type, config] of Object.entries(this.MEETING_PATTERNS)) {
      if (config.urlPattern.test(url)) {
        return type as MeetingType;
      }
    }
    return null;
  }

  static extractMeetingInfo(url: string): MeetingInfo {
    const type = this.detectMeetingType(url);
    
    if (!type) {
      return {
        type: null,
        url,
        meetingId: '',
      };
    }

    const config = this.MEETING_PATTERNS[type];
    const match = url.match(config.urlPattern);
    const meetingId = match ? match[1] : '';

    return {
      type,
      url,
      meetingId,
      title: this.extractTitle(config.titleSelector),
      participants: this.extractParticipants(config.participantSelector),
    };
  }

  private static extractTitle(selector: string): string | undefined {
    const element = document.querySelector(selector);
    return element?.textContent?.trim();
  }

  private static extractParticipants(selector: string): string[] {
    const elements = document.querySelectorAll(selector);
    return Array.from(elements).map(el => el.textContent?.trim() || '');
  }

  static observeMeetingChanges(callback: (meetingInfo: MeetingInfo) => void): () => void {
    const observer = new MutationObserver(() => {
      const meetingInfo = this.extractMeetingInfo(window.location.href);
      callback(meetingInfo);
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
    });

    // Initial detection
    callback(this.extractMeetingInfo(window.location.href));

    return () => observer.disconnect();
  }
}

// Initialize meeting detection
let currentMeeting: MeetingInfo | null = null;

const stopObserving = MeetingDetector.observeMeetingChanges((meetingInfo) => {
  if (meetingInfo.type && meetingInfo.type !== currentMeeting?.type) {
    currentMeeting = meetingInfo;
    
    // Send to background script
    browser.runtime.sendMessage({
      type: 'MEETING_DETECTED',
      payload: meetingInfo,
    });
  }
});

// Cleanup on page unload
window.addEventListener('beforeunload', stopObserving);
```

This implementation provides a solid foundation for the meeting transcription and action item extraction system, following DDD principles and maintaining loose coupling between components through well-defined interfaces and the strategy pattern. 