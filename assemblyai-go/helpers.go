package assemblyai

// Helper functions for creating pointers to basic types
// These are useful for optional fields in request structs

// Bool returns a pointer to the bool value
func Bool(v bool) *bool {
	return &v
}

// String returns a pointer to the string value
func String(v string) *string {
	return &v
}

// Int returns a pointer to the int value
func Int(v int) *int {
	return &v
}

// Float64 returns a pointer to the float64 value
func Float64(v float64) *float64 {
	return &v
}

// NewTranscriptRequest creates a new TranscriptRequest with common defaults
func NewTranscriptRequest(audioURL string) *TranscriptRequest {
	return &TranscriptRequest{
		AudioURL:   audioURL,
		Punctuate:  Bool(true),
		FormatText: Bool(true),
	}
}

// WithSpeakerLabels enables speaker labels for the transcript request
func (tr *TranscriptRequest) WithSpeakerLabels(enabled bool) *TranscriptRequest {
	tr.SpeakerLabels = Bool(enabled)
	return tr
}

// WithLanguageCode sets the language code for the transcript request
func (tr *TranscriptRequest) WithLanguageCode(code string) *TranscriptRequest {
	tr.LanguageCode = String(code)
	return tr
}

// WithPunctuation enables or disables punctuation for the transcript request
func (tr *TranscriptRequest) WithPunctuation(enabled bool) *TranscriptRequest {
	tr.Punctuate = Bool(enabled)
	return tr
}

// WithFormatText enables or disables text formatting for the transcript request
func (tr *TranscriptRequest) WithFormatText(enabled bool) *TranscriptRequest {
	tr.FormatText = Bool(enabled)
	return tr
}

// WithDualChannel enables dual channel processing for the transcript request
func (tr *TranscriptRequest) WithDualChannel(enabled bool) *TranscriptRequest {
	tr.DualChannel = Bool(enabled)
	return tr
}

// WithWebhook sets the webhook URL for the transcript request
func (tr *TranscriptRequest) WithWebhook(url string) *TranscriptRequest {
	tr.WebhookURL = String(url)
	return tr
}

// WithAutoHighlights enables auto highlights for the transcript request
func (tr *TranscriptRequest) WithAutoHighlights(enabled bool) *TranscriptRequest {
	tr.AutoHighlights = Bool(enabled)
	return tr
}

// WithContentSafety enables content safety detection for the transcript request
func (tr *TranscriptRequest) WithContentSafety(enabled bool) *TranscriptRequest {
	tr.ContentSafety = Bool(enabled)
	return tr
}

// WithSentimentAnalysis enables sentiment analysis for the transcript request
func (tr *TranscriptRequest) WithSentimentAnalysis(enabled bool) *TranscriptRequest {
	tr.SentimentAnalysis = Bool(enabled)
	return tr
}

// WithEntityDetection enables entity detection for the transcript request
func (tr *TranscriptRequest) WithEntityDetection(enabled bool) *TranscriptRequest {
	tr.EntityDetection = Bool(enabled)
	return tr
}

// WithAutoChapters enables auto chapters for the transcript request
func (tr *TranscriptRequest) WithAutoChapters(enabled bool) *TranscriptRequest {
	tr.AutoChapters = Bool(enabled)
	return tr
}

// WithSummarization enables summarization for the transcript request
func (tr *TranscriptRequest) WithSummarization(enabled bool) *TranscriptRequest {
	tr.Summarization = Bool(enabled)
	return tr
}

// WithRedactPII enables PII redaction for the transcript request
func (tr *TranscriptRequest) WithRedactPII(enabled bool, policies ...PIIPolicy) *TranscriptRequest {
	tr.RedactPII = Bool(enabled)
	if len(policies) > 0 {
		tr.RedactPIIPolicies = policies
	}
	return tr
}

// WithWordBoost sets word boost for the transcript request
func (tr *TranscriptRequest) WithWordBoost(words []string, param string) *TranscriptRequest {
	tr.WordBoost = words
	if param != "" {
		tr.BoostParam = String(param)
	}
	return tr
}

// WithCustomSpelling sets custom spelling for the transcript request
func (tr *TranscriptRequest) WithCustomSpelling(spelling []CustomSpelling) *TranscriptRequest {
	tr.CustomSpelling = spelling
	return tr
}

// WithAudioSegment sets the audio start and end times for the transcript request
func (tr *TranscriptRequest) WithAudioSegment(startMs, endMs int) *TranscriptRequest {
	tr.AudioStartFrom = Int(startMs)
	tr.AudioEndAt = Int(endMs)
	return tr
}

// NewLemurRequest creates a new LeMUR request
func NewLemurRequest(transcriptIDs []string, prompt string) *LemurRequest {
	return &LemurRequest{
		TranscriptIDs: transcriptIDs,
		Prompt:        prompt,
	}
}

// WithContext sets the context for the LeMUR request
func (lr *LemurRequest) WithContext(context string) *LemurRequest {
	lr.Context = String(context)
	return lr
}

// WithModel sets the model for the LeMUR request
func (lr *LemurRequest) WithModel(model string) *LemurRequest {
	lr.FinalModel = String(model)
	return lr
}

// WithTemperature sets the temperature for the LeMUR request
func (lr *LemurRequest) WithTemperature(temperature float64) *LemurRequest {
	lr.Temperature = Float64(temperature)
	return lr
}

// WithMaxOutputSize sets the max output size for the LeMUR request
func (lr *LemurRequest) WithMaxOutputSize(size int) *LemurRequest {
	lr.MaxOutputSize = Int(size)
	return lr
}

// NewLemurSummaryRequest creates a new LeMUR summary request
func NewLemurSummaryRequest(transcriptIDs []string) *LemurSummaryRequest {
	return &LemurSummaryRequest{
		TranscriptIDs: transcriptIDs,
	}
}

// WithAnswerFormat sets the answer format for the LeMUR summary request
func (lsr *LemurSummaryRequest) WithAnswerFormat(format string) *LemurSummaryRequest {
	lsr.AnswerFormat = String(format)
	return lsr
}

// NewLemurQuestionAnswerRequest creates a new LeMUR Q&A request
func NewLemurQuestionAnswerRequest(transcriptIDs []string, questions []LemurQuestion) *LemurQuestionAnswerRequest {
	return &LemurQuestionAnswerRequest{
		TranscriptIDs: transcriptIDs,
		Questions:     questions,
	}
}

// NewLemurActionItemsRequest creates a new LeMUR action items request
func NewLemurActionItemsRequest(transcriptIDs []string) *LemurActionItemsRequest {
	return &LemurActionItemsRequest{
		TranscriptIDs: transcriptIDs,
	}
}
