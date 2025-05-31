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

// Transcription represents a meeting transcription entity
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

// NewTranscription creates a new Transcription entity
func NewTranscription(meetingID, audioFilePath, provider string) Transcription {
	transcription := Transcription{
		MeetingID:     meetingID,
		AudioFilePath: audioFilePath,
		Status:        Pending,
		Provider:      provider,
	}
	transcription.SetID(domain.GenerateID())
	return transcription
}

// StartProcessing transitions the transcription to processing status
func (t *Transcription) StartProcessing() {
	t.Status = Processing
}

// CompleteTranscription transitions the transcription to completed status
func (t *Transcription) CompleteTranscription(content string, confidence float64, segments []TranscriptSegment) {
	t.Status = Completed
	t.Content = content
	t.Confidence = confidence
	t.Segments = segments
}

// FailTranscription transitions the transcription to failed status
func (t *Transcription) FailTranscription() {
	t.Status = Failed
}

// IsCompleted returns true if the transcription has completed successfully
func (t *Transcription) IsCompleted() bool {
	return t.Status == Completed
}

// IsProcessing returns true if the transcription is currently being processed
func (t *Transcription) IsProcessing() bool {
	return t.Status == Processing
}

// HasContent returns true if the transcription has content
func (t *Transcription) HasContent() bool {
	return t.Content != ""
}

// GetWordCount returns the approximate word count of the transcription
func (t *Transcription) GetWordCount() int {
	if t.Content == "" {
		return 0
	}
	// Simple word count - split by spaces
	words := 0
	inWord := false
	for _, char := range t.Content {
		if char == ' ' || char == '\t' || char == '\n' {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	return words
}

// AddSegment adds a transcript segment to the transcription
func (t *Transcription) AddSegment(speaker, text string, startTime, endTime, confidence float64, sequenceNumber int) {
	segment := NewTranscriptSegment(t.GetID(), speaker, text, startTime, endTime, confidence, sequenceNumber)
	t.Segments = append(t.Segments, segment)
}

// TableName sets the table name for GORM
func (Transcription) TableName() string {
	return "transcriptions"
}

// TranscriptSegment represents an individual segment of a transcription
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

// NewTranscriptSegment creates a new TranscriptSegment entity
func NewTranscriptSegment(transcriptionID, speaker, text string, startTime, endTime, confidence float64, sequenceNumber int) TranscriptSegment {
	segment := TranscriptSegment{
		TranscriptionID: transcriptionID,
		Speaker:         speaker,
		Text:            text,
		StartTime:       startTime,
		EndTime:         endTime,
		Confidence:      confidence,
		SequenceNumber:  sequenceNumber,
	}
	segment.SetID(domain.GenerateID())
	return segment
}

// GetDuration returns the duration of the segment in seconds
func (ts *TranscriptSegment) GetDuration() float64 {
	return ts.EndTime - ts.StartTime
}

// GetWordCount returns the approximate word count of the segment
func (ts *TranscriptSegment) GetWordCount() int {
	if ts.Text == "" {
		return 0
	}
	words := 0
	inWord := false
	for _, char := range ts.Text {
		if char == ' ' || char == '\t' || char == '\n' {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	return words
}

// IsHighConfidence returns true if the segment has high confidence (>= 0.8)
func (ts *TranscriptSegment) IsHighConfidence() bool {
	return ts.Confidence >= 0.8
}

// TableName sets the table name for GORM
func (TranscriptSegment) TableName() string {
	return "transcript_segments"
}
