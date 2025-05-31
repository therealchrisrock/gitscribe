package assemblyai

import (
	"fmt"
)

// APIError represents an error response from the AssemblyAI API
type APIError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("AssemblyAI API error (status %d): %s", e.StatusCode, e.Message)
}

// TranscriptStatus represents the status of a transcript
type TranscriptStatus string

const (
	StatusQueued     TranscriptStatus = "queued"
	StatusProcessing TranscriptStatus = "processing"
	StatusCompleted  TranscriptStatus = "completed"
	StatusError      TranscriptStatus = "error"
)

// TranscriptRequest represents a request to create a transcript
type TranscriptRequest struct {
	AudioURL               string                 `json:"audio_url"`
	LanguageCode           *string                `json:"language_code,omitempty"`
	Punctuate              *bool                  `json:"punctuate,omitempty"`
	FormatText             *bool                  `json:"format_text,omitempty"`
	DualChannel            *bool                  `json:"dual_channel,omitempty"`
	WebhookURL             *string                `json:"webhook_url,omitempty"`
	WebhookAuthHeaderName  *string                `json:"webhook_auth_header_name,omitempty"`
	WebhookAuthHeaderValue *string                `json:"webhook_auth_header_value,omitempty"`
	AutoHighlights         *bool                  `json:"auto_highlights,omitempty"`
	AudioStartFrom         *int                   `json:"audio_start_from,omitempty"`
	AudioEndAt             *int                   `json:"audio_end_at,omitempty"`
	WordBoost              []string               `json:"word_boost,omitempty"`
	BoostParam             *string                `json:"boost_param,omitempty"`
	FilterProfanity        *bool                  `json:"filter_profanity,omitempty"`
	RedactPII              *bool                  `json:"redact_pii,omitempty"`
	RedactPIIAudio         *bool                  `json:"redact_pii_audio,omitempty"`
	RedactPIIPolicies      []PIIPolicy            `json:"redact_pii_policies,omitempty"`
	RedactPIISub           *string                `json:"redact_pii_sub,omitempty"`
	SpeakerLabels          *bool                  `json:"speaker_labels,omitempty"`
	SpeakersExpected       *int                   `json:"speakers_expected,omitempty"`
	ContentSafety          *bool                  `json:"content_safety,omitempty"`
	IabCategories          *bool                  `json:"iab_categories,omitempty"`
	LanguageDetection      *bool                  `json:"language_detection,omitempty"`
	CustomSpelling         []CustomSpelling       `json:"custom_spelling,omitempty"`
	Disfluencies           *bool                  `json:"disfluencies,omitempty"`
	SentimentAnalysis      *bool                  `json:"sentiment_analysis,omitempty"`
	AutoChapters           *bool                  `json:"auto_chapters,omitempty"`
	EntityDetection        *bool                  `json:"entity_detection,omitempty"`
	SpeechThreshold        *float64               `json:"speech_threshold,omitempty"`
	Summarization          *bool                  `json:"summarization,omitempty"`
	SummaryModel           *string                `json:"summary_model,omitempty"`
	SummaryType            *string                `json:"summary_type,omitempty"`
	CustomTopics           *bool                  `json:"custom_topics,omitempty"`
	Topics                 []string               `json:"topics,omitempty"`
	AdditionalProperties   map[string]interface{} `json:"-"`
}

// PIIPolicy represents a PII redaction policy
type PIIPolicy string

const (
	PIIPolicyMedicalProcess       PIIPolicy = "medical_process"
	PIIPolicyMedicalCondition     PIIPolicy = "medical_condition"
	PIIPolicyBloodType            PIIPolicy = "blood_type"
	PIIPolicyDrug                 PIIPolicy = "drug"
	PIIPolicyInjury               PIIPolicy = "injury"
	PIIPolicyNumberSequence       PIIPolicy = "number_sequence"
	PIIPolicyEmailAddress         PIIPolicy = "email_address"
	PIIPolicyDateOfBirth          PIIPolicy = "date_of_birth"
	PIIPolicyPhoneNumber          PIIPolicy = "phone_number"
	PIIPolicyUSSSN                PIIPolicy = "us_social_security_number"
	PIIPolicyCreditCardNumber     PIIPolicy = "credit_card_number"
	PIIPolicyCreditCardCVV        PIIPolicy = "credit_card_cvv"
	PIIPolicyCreditCardExpiry     PIIPolicy = "credit_card_expiration"
	PIIPolicyPersonName           PIIPolicy = "person_name"
	PIIPolicyPersonAge            PIIPolicy = "person_age"
	PIIPolicyOrganization         PIIPolicy = "organization"
	PIIPolicyLocation             PIIPolicy = "location"
	PIIPolicyEvent                PIIPolicy = "event"
	PIIPolicyLanguage             PIIPolicy = "language"
	PIIPolicyNationality          PIIPolicy = "nationality"
	PIIPolicyReligion             PIIPolicy = "religion"
	PIIPolicyPoliticalAffiliation PIIPolicy = "political_affiliation"
	PIIPolicyOccupation           PIIPolicy = "occupation"
	PIIPolicyUSBankNumber         PIIPolicy = "us_bank_number"
	PIIPolicyUSDriversLicense     PIIPolicy = "us_drivers_license"
	PIIPolicyUSPassportNumber     PIIPolicy = "us_passport_number"
)

// CustomSpelling represents custom spelling configuration
type CustomSpelling struct {
	From []string `json:"from"`
	To   string   `json:"to"`
}

// Transcript represents a transcript response
type Transcript struct {
	ID                          string                    `json:"id"`
	LanguageModel               string                    `json:"language_model"`
	AcousticModel               string                    `json:"acoustic_model"`
	LanguageCode                string                    `json:"language_code"`
	Status                      TranscriptStatus          `json:"status"`
	AudioURL                    string                    `json:"audio_url"`
	Text                        *string                   `json:"text"`
	Words                       []Word                    `json:"words"`
	Utterances                  []Utterance               `json:"utterances"`
	Confidence                  *float64                  `json:"confidence"`
	AudioDuration               *float64                  `json:"audio_duration"`
	Punctuate                   bool                      `json:"punctuate"`
	FormatText                  bool                      `json:"format_text"`
	DualChannel                 *bool                     `json:"dual_channel"`
	WebhookURL                  *string                   `json:"webhook_url"`
	WebhookStatusCode           *int                      `json:"webhook_status_code"`
	WebhookAuthHeaderName       *string                   `json:"webhook_auth_header_name"`
	WebhookAuthHeaderValue      *string                   `json:"webhook_auth_header_value"`
	AutoHighlights              *bool                     `json:"auto_highlights"`
	AutoHighlightsResult        *AutoHighlightsResult     `json:"auto_highlights_result"`
	AudioStartFrom              *int                      `json:"audio_start_from"`
	AudioEndAt                  *int                      `json:"audio_end_at"`
	WordBoost                   []string                  `json:"word_boost"`
	BoostParam                  *string                   `json:"boost_param"`
	FilterProfanity             bool                      `json:"filter_profanity"`
	RedactPII                   bool                      `json:"redact_pii"`
	RedactPIIAudio              bool                      `json:"redact_pii_audio"`
	RedactPIIPolicies           []PIIPolicy               `json:"redact_pii_policies"`
	RedactPIISub                string                    `json:"redact_pii_sub"`
	SpeakerLabels               bool                      `json:"speaker_labels"`
	SpeakersExpected            *int                      `json:"speakers_expected"`
	ContentSafety               bool                      `json:"content_safety"`
	ContentSafetyLabels         *ContentSafetyLabels      `json:"content_safety_labels"`
	IabCategories               bool                      `json:"iab_categories"`
	IabCategoriesResult         *IabCategoriesResult      `json:"iab_categories_result"`
	LanguageDetection           bool                      `json:"language_detection"`
	LanguageConfidenceThreshold *float64                  `json:"language_confidence_threshold"`
	CustomSpelling              []CustomSpelling          `json:"custom_spelling"`
	Disfluencies                bool                      `json:"disfluencies"`
	SentimentAnalysis           bool                      `json:"sentiment_analysis"`
	SentimentAnalysisResults    []SentimentAnalysisResult `json:"sentiment_analysis_results"`
	AutoChapters                bool                      `json:"auto_chapters"`
	Chapters                    []Chapter                 `json:"chapters"`
	EntityDetection             bool                      `json:"entity_detection"`
	Entities                    []Entity                  `json:"entities"`
	SpeechThreshold             *float64                  `json:"speech_threshold"`
	Summarization               bool                      `json:"summarization"`
	Summary                     *string                   `json:"summary"`
	SummaryModel                *string                   `json:"summary_model"`
	SummaryType                 *string                   `json:"summary_type"`
	CustomTopics                bool                      `json:"custom_topics"`
	Topics                      []string                  `json:"topics"`
	TopicDetectionResults       []TopicDetectionResult    `json:"topic_detection_results"`
	Error                       *string                   `json:"error"`
	ThrottledBy                 *string                   `json:"throttled_by"`
}

// Word represents a word in the transcript
type Word struct {
	Confidence float64 `json:"confidence"`
	End        int     `json:"end"`
	Start      int     `json:"start"`
	Text       string  `json:"text"`
	Speaker    *string `json:"speaker,omitempty"`
}

// Utterance represents an utterance in the transcript
type Utterance struct {
	Confidence float64 `json:"confidence"`
	End        int     `json:"end"`
	Start      int     `json:"start"`
	Text       string  `json:"text"`
	Words      []Word  `json:"words"`
	Speaker    string  `json:"speaker"`
}

// AutoHighlightsResult represents auto highlights results
type AutoHighlightsResult struct {
	Status  string      `json:"status"`
	Results []Highlight `json:"results"`
}

// Highlight represents a highlight
type Highlight struct {
	Count      int     `json:"count"`
	Rank       float64 `json:"rank"`
	Text       string  `json:"text"`
	Timestamps []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"timestamps"`
}

// ContentSafetyLabels represents content safety labels
type ContentSafetyLabels struct {
	Status  string                     `json:"status"`
	Results []ContentSafetyLabelResult `json:"results"`
	Summary map[string]float64         `json:"summary"`
}

// ContentSafetyLabelResult represents a content safety label result
type ContentSafetyLabelResult struct {
	Text      string                       `json:"text"`
	Labels    []ContentSafetyLabelCategory `json:"labels"`
	Timestamp struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"timestamp"`
}

// ContentSafetyLabelCategory represents a content safety label category
type ContentSafetyLabelCategory struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
	Severity   float64 `json:"severity"`
}

// IabCategoriesResult represents IAB categories result
type IabCategoriesResult struct {
	Status  string              `json:"status"`
	Results []IabCategoryResult `json:"results"`
	Summary map[string]float64  `json:"summary"`
}

// IabCategoryResult represents an IAB category result
type IabCategoryResult struct {
	Text      string             `json:"text"`
	Labels    []IabCategoryLabel `json:"labels"`
	Timestamp struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"timestamp"`
}

// IabCategoryLabel represents an IAB category label
type IabCategoryLabel struct {
	Relevance float64 `json:"relevance"`
	Label     string  `json:"label"`
}

// SentimentAnalysisResult represents sentiment analysis result
type SentimentAnalysisResult struct {
	Text       string  `json:"text"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
	Speaker    *string `json:"speaker,omitempty"`
}

// Chapter represents a chapter
type Chapter struct {
	Summary  string `json:"summary"`
	Headline string `json:"headline"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Gist     string `json:"gist"`
}

// Entity represents an entity
type Entity struct {
	EntityType string `json:"entity_type"`
	Text       string `json:"text"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
}

// TopicDetectionResult represents topic detection result
type TopicDetectionResult struct {
	Text      string                `json:"text"`
	Labels    []TopicDetectionLabel `json:"labels"`
	Timestamp struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"timestamp"`
}

// TopicDetectionLabel represents a topic detection label
type TopicDetectionLabel struct {
	Relevance float64 `json:"relevance"`
	Label     string  `json:"label"`
}

// UploadResponse represents the response from uploading a file
type UploadResponse struct {
	UploadURL string `json:"upload_url"`
}

// ListTranscriptsResponse represents the response from listing transcripts
type ListTranscriptsResponse struct {
	PageDetails PageDetails  `json:"page_details"`
	Transcripts []Transcript `json:"transcripts"`
}

// PageDetails represents pagination details
type PageDetails struct {
	Limit       int     `json:"limit"`
	ResultCount int     `json:"result_count"`
	CurrentURL  string  `json:"current_url"`
	PrevURL     *string `json:"prev_url"`
	NextURL     *string `json:"next_url"`
}

// LemurRequest represents a LeMUR request
type LemurRequest struct {
	TranscriptIDs        []string               `json:"transcript_ids"`
	Prompt               string                 `json:"prompt"`
	Context              *string                `json:"context,omitempty"`
	FinalModel           *string                `json:"final_model,omitempty"`
	MaxOutputSize        *int                   `json:"max_output_size,omitempty"`
	Temperature          *float64               `json:"temperature,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurResponse represents a LeMUR response
type LemurResponse struct {
	RequestID            string                 `json:"request_id"`
	Response             string                 `json:"response"`
	Usage                LemurUsage             `json:"usage"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurUsage represents LeMUR usage statistics
type LemurUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// LemurSummaryRequest represents a LeMUR summary request
type LemurSummaryRequest struct {
	TranscriptIDs        []string               `json:"transcript_ids"`
	Context              *string                `json:"context,omitempty"`
	AnswerFormat         *string                `json:"answer_format,omitempty"`
	FinalModel           *string                `json:"final_model,omitempty"`
	MaxOutputSize        *int                   `json:"max_output_size,omitempty"`
	Temperature          *float64               `json:"temperature,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurQuestionAnswerRequest represents a LeMUR Q&A request
type LemurQuestionAnswerRequest struct {
	TranscriptIDs        []string               `json:"transcript_ids"`
	Questions            []LemurQuestion        `json:"questions"`
	Context              *string                `json:"context,omitempty"`
	FinalModel           *string                `json:"final_model,omitempty"`
	MaxOutputSize        *int                   `json:"max_output_size,omitempty"`
	Temperature          *float64               `json:"temperature,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurQuestion represents a LeMUR question
type LemurQuestion struct {
	Question      string   `json:"question"`
	Context       *string  `json:"context,omitempty"`
	AnswerFormat  *string  `json:"answer_format,omitempty"`
	AnswerOptions []string `json:"answer_options,omitempty"`
}

// LemurQuestionAnswerResponse represents a LeMUR Q&A response
type LemurQuestionAnswerResponse struct {
	RequestID            string                 `json:"request_id"`
	Response             []LemurAnswer          `json:"response"`
	Usage                LemurUsage             `json:"usage"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurAnswer represents a LeMUR answer
type LemurAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// LemurActionItemsRequest represents a LeMUR action items request
type LemurActionItemsRequest struct {
	TranscriptIDs        []string               `json:"transcript_ids"`
	Context              *string                `json:"context,omitempty"`
	AnswerFormat         *string                `json:"answer_format,omitempty"`
	FinalModel           *string                `json:"final_model,omitempty"`
	MaxOutputSize        *int                   `json:"max_output_size,omitempty"`
	Temperature          *float64               `json:"temperature,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// LemurActionItemsResponse represents a LeMUR action items response
type LemurActionItemsResponse struct {
	RequestID            string                 `json:"request_id"`
	Response             string                 `json:"response"`
	ActionItems          []string               `json:"action_items"`
	Usage                LemurUsage             `json:"usage"`
	AdditionalProperties map[string]interface{} `json:"-"`
}
