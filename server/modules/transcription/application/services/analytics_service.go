package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	meetingRepos "teammate/server/modules/meeting/domain/repositories"
	"teammate/server/modules/transcription/domain/entities"
	"teammate/server/modules/transcription/domain/repositories"
)

// AnalyticsService provides advanced analytics and insights for transcriptions
type AnalyticsService struct {
	transcriptionRepo repositories.TranscriptionRepository
	meetingRepo       meetingRepos.MeetingRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	transcriptionRepo repositories.TranscriptionRepository,
	meetingRepo meetingRepos.MeetingRepository,
) *AnalyticsService {
	return &AnalyticsService{
		transcriptionRepo: transcriptionRepo,
		meetingRepo:       meetingRepo,
	}
}

// AnalyticsData represents comprehensive analytics for a transcription
type AnalyticsData struct {
	TranscriptionID   string             `json:"transcription_id"`
	MeetingID         string             `json:"meeting_id"`
	SpeakerAnalytics  []SpeakerAnalytics `json:"speaker_analytics"`
	TopicAnalysis     []TopicAnalysis    `json:"topic_analysis"`
	SentimentAnalysis *SentimentAnalysis `json:"sentiment_analysis"`
	MeetingMetrics    *MeetingMetrics    `json:"meeting_metrics"`
	KeywordFrequency  []KeywordFrequency `json:"keyword_frequency"`
	TimeDistribution  *TimeDistribution  `json:"time_distribution"`
	QualityMetrics    *QualityMetrics    `json:"quality_metrics"`
	Insights          []Insight          `json:"insights"`
	GeneratedAt       time.Time          `json:"generated_at"`
}

// SpeakerAnalytics provides detailed analysis for each speaker
type SpeakerAnalytics struct {
	Speaker               string   `json:"speaker"`
	SpeakingTime          float64  `json:"speaking_time_seconds"`
	SpeakingTimePercent   float64  `json:"speaking_time_percent"`
	WordCount             int      `json:"word_count"`
	AverageWordsPerMinute float64  `json:"average_words_per_minute"`
	SegmentCount          int      `json:"segment_count"`
	AverageConfidence     float64  `json:"average_confidence"`
	LongestSegment        float64  `json:"longest_segment_seconds"`
	ShortestSegment       float64  `json:"shortest_segment_seconds"`
	SpeakingPeriods       []Period `json:"speaking_periods"`
	TalkingPoints         []string `json:"talking_points"`
}

// TopicAnalysis identifies key topics discussed
type TopicAnalysis struct {
	Topic        string   `json:"topic"`
	Relevance    float64  `json:"relevance"`
	Keywords     []string `json:"keywords"`
	Mentions     int      `json:"mentions"`
	FirstMention float64  `json:"first_mention_time"`
	LastMention  float64  `json:"last_mention_time"`
	Speakers     []string `json:"speakers"`
	Summary      string   `json:"summary,omitempty"`
}

// SentimentAnalysis provides sentiment metrics
type SentimentAnalysis struct {
	OverallSentiment    string               `json:"overall_sentiment"`
	SentimentScore      float64              `json:"sentiment_score"` // -1 to 1
	PositiveSegments    int                  `json:"positive_segments"`
	NeutralSegments     int                  `json:"neutral_segments"`
	NegativeSegments    int                  `json:"negative_segments"`
	SentimentTimeline   []SentimentTimePoint `json:"sentiment_timeline"`
	SpeakerSentiments   []SpeakerSentiment   `json:"speaker_sentiments"`
	EmotionalHighlights []EmotionalHighlight `json:"emotional_highlights"`
}

// MeetingMetrics provides overall meeting statistics
type MeetingMetrics struct {
	TotalDuration      float64 `json:"total_duration_seconds"`
	ActiveSpeakingTime float64 `json:"active_speaking_time_seconds"`
	SilenceTime        float64 `json:"silence_time_seconds"`
	SpeakerCount       int     `json:"speaker_count"`
	SegmentCount       int     `json:"segment_count"`
	WordCount          int     `json:"total_word_count"`
	AverageConfidence  float64 `json:"average_confidence"`
	OverallPace        string  `json:"overall_pace"` // slow, normal, fast
	InterruptionCount  int     `json:"interruption_count"`
	LongestMonologue   float64 `json:"longest_monologue_seconds"`
	SpeakerTurnover    float64 `json:"speaker_turnover_rate"`
}

// KeywordFrequency tracks important terms
type KeywordFrequency struct {
	Keyword   string   `json:"keyword"`
	Count     int      `json:"count"`
	Frequency float64  `json:"frequency"`
	TfIdf     float64  `json:"tf_idf"`
	Context   []string `json:"context,omitempty"`
}

// TimeDistribution shows activity over time
type TimeDistribution struct {
	TimeSlots     []TimeSlot `json:"time_slots"`
	PeakActivity  []Period   `json:"peak_activity_periods"`
	QuietPeriods  []Period   `json:"quiet_periods"`
	ActivityCurve []float64  `json:"activity_curve"`
}

// QualityMetrics assesses transcription quality
type QualityMetrics struct {
	OverallQuality         string         `json:"overall_quality"` // excellent, good, fair, poor
	AverageConfidence      float64        `json:"average_confidence"`
	LowConfidenceSegments  int            `json:"low_confidence_segments"`
	ConfidenceDistribution map[string]int `json:"confidence_distribution"`
	AudioQuality           string         `json:"audio_quality_assessment"`
	RecommendedActions     []string       `json:"recommended_actions"`
}

// Insight represents an AI-generated insight
type Insight struct {
	Type        string   `json:"type"` // speaker_behavior, topic_trend, quality_issue, etc.
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Severity    string   `json:"severity"` // low, medium, high
	Timestamp   float64  `json:"timestamp,omitempty"`
	Speaker     string   `json:"speaker,omitempty"`
	ActionItems []string `json:"action_items,omitempty"`
}

// Supporting types
type Period struct {
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Duration  float64 `json:"duration"`
}

type TimeSlot struct {
	StartTime     float64 `json:"start_time"`
	EndTime       float64 `json:"end_time"`
	WordCount     int     `json:"word_count"`
	SpeakerCount  int     `json:"speaker_count"`
	ActivityLevel string  `json:"activity_level"` // low, medium, high
}

type SentimentTimePoint struct {
	Timestamp float64 `json:"timestamp"`
	Sentiment string  `json:"sentiment"`
	Score     float64 `json:"score"`
}

type SpeakerSentiment struct {
	Speaker          string  `json:"speaker"`
	AverageSentiment float64 `json:"average_sentiment"`
	PositiveSegments int     `json:"positive_segments"`
	NegativeSegments int     `json:"negative_segments"`
	MostPositive     string  `json:"most_positive_quote,omitempty"`
	MostNegative     string  `json:"most_negative_quote,omitempty"`
}

type EmotionalHighlight struct {
	Timestamp float64 `json:"timestamp"`
	Speaker   string  `json:"speaker"`
	Emotion   string  `json:"emotion"`
	Intensity float64 `json:"intensity"`
	Text      string  `json:"text"`
	Context   string  `json:"context"`
}

// GetTranscriptionAnalytics generates comprehensive analytics for a transcription
func (s *AnalyticsService) GetTranscriptionAnalytics(ctx context.Context, transcriptionID string) (*AnalyticsData, error) {
	// Get transcription and segments
	transcription, err := s.transcriptionRepo.FindByID(ctx, transcriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transcription: %w", err)
	}

	segments, err := s.transcriptionRepo.FindSegmentsByTranscriptionID(ctx, transcriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find segments: %w", err)
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments found for transcription %s", transcriptionID)
	}

	analytics := &AnalyticsData{
		TranscriptionID: transcriptionID,
		MeetingID:       transcription.MeetingID,
		GeneratedAt:     time.Now(),
	}

	// Generate all analytics components
	analytics.SpeakerAnalytics = s.analyzeSpeakers(segments)
	analytics.TopicAnalysis = s.analyzeTopics(segments)
	analytics.SentimentAnalysis = s.analyzeSentiment(segments)
	analytics.MeetingMetrics = s.calculateMeetingMetrics(segments)
	analytics.KeywordFrequency = s.analyzeKeywords(segments)
	analytics.TimeDistribution = s.analyzeTimeDistribution(segments)
	analytics.QualityMetrics = s.assessQuality(segments)
	analytics.Insights = s.generateInsights(analytics)

	return analytics, nil
}

// analyzeSpeakers provides detailed speaker analytics
func (s *AnalyticsService) analyzeSpeakers(segments []entities.TranscriptSegment) []SpeakerAnalytics {
	speakerData := make(map[string]*SpeakerAnalytics)
	totalDuration := s.calculateTotalDuration(segments)

	for _, segment := range segments {
		speaker := segment.Speaker
		// Only normalize truly unknown speakers - preserve differentiation
		if speaker == "" || speaker == "speaker_unknown" {
			speaker = "Speaker Unknown"
		}

		if speakerData[speaker] == nil {
			speakerData[speaker] = &SpeakerAnalytics{
				Speaker:         speaker,
				SpeakingPeriods: []Period{},
				TalkingPoints:   []string{},
				LongestSegment:  0,
				ShortestSegment: math.MaxFloat64,
			}
		}

		data := speakerData[speaker]
		duration := segment.EndTime - segment.StartTime
		words := len(strings.Fields(segment.Text))

		// Update metrics
		data.SpeakingTime += duration
		data.WordCount += words
		data.SegmentCount++
		data.AverageConfidence = (data.AverageConfidence*float64(data.SegmentCount-1) + segment.Confidence) / float64(data.SegmentCount)

		if duration > data.LongestSegment {
			data.LongestSegment = duration
		}
		if duration < data.ShortestSegment {
			data.ShortestSegment = duration
		}

		// Add speaking period
		data.SpeakingPeriods = append(data.SpeakingPeriods, Period{
			StartTime: segment.StartTime,
			EndTime:   segment.EndTime,
			Duration:  duration,
		})

		// Extract talking points (sentences that are longer and have good confidence)
		if len(segment.Text) > 50 && segment.Confidence > 0.8 {
			data.TalkingPoints = append(data.TalkingPoints, segment.Text)
		}
	}

	// Calculate derived metrics
	var result []SpeakerAnalytics
	for _, data := range speakerData {
		if data.SpeakingTime > 0 {
			data.AverageWordsPerMinute = float64(data.WordCount) / (data.SpeakingTime / 60.0)
		}
		if totalDuration > 0 {
			data.SpeakingTimePercent = (data.SpeakingTime / totalDuration) * 100
		}
		if data.ShortestSegment == math.MaxFloat64 {
			data.ShortestSegment = 0
		}

		// Limit talking points to top 5
		if len(data.TalkingPoints) > 5 {
			data.TalkingPoints = data.TalkingPoints[:5]
		}

		result = append(result, *data)
	}

	// Sort by speaking time
	sort.Slice(result, func(i, j int) bool {
		return result[i].SpeakingTime > result[j].SpeakingTime
	})

	return result
}

// analyzeTopics identifies key topics (simplified implementation)
func (s *AnalyticsService) analyzeTopics(segments []entities.TranscriptSegment) []TopicAnalysis {
	// This is a simplified topic analysis
	// In a real implementation, you'd use NLP libraries or services

	topicKeywords := map[string][]string{
		"Meeting Management":   {"meeting", "agenda", "schedule", "plan", "organize"},
		"Technical Discussion": {"technical", "development", "code", "system", "architecture", "bug", "feature"},
		"Business Strategy":    {"business", "strategy", "market", "customer", "revenue", "growth"},
		"Project Planning":     {"project", "timeline", "deadline", "milestone", "deliverable", "task"},
		"Team Coordination":    {"team", "collaboration", "communication", "assign", "responsibility"},
		"Decision Making":      {"decision", "approve", "choose", "select", "vote", "agree"},
		"Problem Solving":      {"problem", "issue", "solution", "fix", "resolve", "troubleshoot"},
		"Review & Feedback":    {"review", "feedback", "comment", "suggestion", "improvement"},
	}

	topicScores := make(map[string]*TopicAnalysis)

	for topic, keywords := range topicKeywords {
		topicScores[topic] = &TopicAnalysis{
			Topic:    topic,
			Keywords: keywords,
			Speakers: []string{},
		}
	}

	// Analyze segments for topic relevance
	for _, segment := range segments {
		text := strings.ToLower(segment.Text)
		speakerMap := make(map[string]bool)

		for _, analysis := range topicScores {
			matches := 0
			for _, keyword := range analysis.Keywords {
				if strings.Contains(text, keyword) {
					matches++
				}
			}

			if matches > 0 {
				analysis.Mentions++
				analysis.Relevance += float64(matches)

				if analysis.FirstMention == 0 || segment.StartTime < analysis.FirstMention {
					analysis.FirstMention = segment.StartTime
				}
				if segment.EndTime > analysis.LastMention {
					analysis.LastMention = segment.EndTime
				}

				if !speakerMap[segment.Speaker] && segment.Speaker != "" {
					analysis.Speakers = append(analysis.Speakers, segment.Speaker)
					speakerMap[segment.Speaker] = true
				}
			}
		}
	}

	// Convert to slice and filter relevant topics
	var result []TopicAnalysis
	for _, analysis := range topicScores {
		if analysis.Mentions > 0 {
			// Normalize relevance score
			analysis.Relevance = analysis.Relevance / float64(len(segments))
			result = append(result, *analysis)
		}
	}

	// Sort by relevance
	sort.Slice(result, func(i, j int) bool {
		return result[i].Relevance > result[j].Relevance
	})

	// Limit to top 10 topics
	if len(result) > 10 {
		result = result[:10]
	}

	return result
}

// analyzeSentiment provides sentiment analysis (simplified)
func (s *AnalyticsService) analyzeSentiment(segments []entities.TranscriptSegment) *SentimentAnalysis {
	// Simplified sentiment analysis using keyword-based approach
	// In production, you'd use ML models or services like Google Cloud Natural Language

	positiveWords := []string{"good", "great", "excellent", "positive", "agree", "success", "happy", "wonderful", "perfect", "amazing"}
	negativeWords := []string{"bad", "terrible", "awful", "negative", "disagree", "failure", "sad", "problem", "issue", "wrong"}

	sentiment := &SentimentAnalysis{
		SentimentTimeline:   []SentimentTimePoint{},
		SpeakerSentiments:   []SpeakerSentiment{},
		EmotionalHighlights: []EmotionalHighlight{},
	}

	totalScore := 0.0
	speakerSentiments := make(map[string]*SpeakerSentiment)

	for _, segment := range segments {
		text := strings.ToLower(segment.Text)
		segmentScore := 0.0

		// Calculate sentiment score for segment
		for _, word := range positiveWords {
			if strings.Contains(text, word) {
				segmentScore += 1.0
			}
		}
		for _, word := range negativeWords {
			if strings.Contains(text, word) {
				segmentScore -= 1.0
			}
		}

		// Normalize score
		words := strings.Fields(text)
		if len(words) > 0 {
			segmentScore = segmentScore / float64(len(words))
		}

		totalScore += segmentScore

		// Categorize segment
		if segmentScore > 0.1 {
			sentiment.PositiveSegments++
		} else if segmentScore < -0.1 {
			sentiment.NegativeSegments++
		} else {
			sentiment.NeutralSegments++
		}

		// Track timeline
		sentimentType := "neutral"
		if segmentScore > 0.1 {
			sentimentType = "positive"
		} else if segmentScore < -0.1 {
			sentimentType = "negative"
		}

		sentiment.SentimentTimeline = append(sentiment.SentimentTimeline, SentimentTimePoint{
			Timestamp: segment.StartTime,
			Sentiment: sentimentType,
			Score:     segmentScore,
		})

		// Track speaker sentiment
		if segment.Speaker != "" {
			if speakerSentiments[segment.Speaker] == nil {
				speakerSentiments[segment.Speaker] = &SpeakerSentiment{
					Speaker: segment.Speaker,
				}
			}

			speakerData := speakerSentiments[segment.Speaker]
			speakerData.AverageSentiment = (speakerData.AverageSentiment + segmentScore) / 2.0

			if segmentScore > 0.1 {
				speakerData.PositiveSegments++
				if speakerData.MostPositive == "" || segmentScore > 0.3 {
					speakerData.MostPositive = segment.Text
				}
			} else if segmentScore < -0.1 {
				speakerData.NegativeSegments++
				if speakerData.MostNegative == "" || segmentScore < -0.3 {
					speakerData.MostNegative = segment.Text
				}
			}
		}

		// Identify emotional highlights
		if math.Abs(segmentScore) > 0.3 {
			emotion := "neutral"
			intensity := math.Abs(segmentScore)

			if segmentScore > 0.3 {
				emotion = "enthusiasm"
			} else if segmentScore < -0.3 {
				emotion = "concern"
			}

			sentiment.EmotionalHighlights = append(sentiment.EmotionalHighlights, EmotionalHighlight{
				Timestamp: segment.StartTime,
				Speaker:   segment.Speaker,
				Emotion:   emotion,
				Intensity: intensity,
				Text:      segment.Text,
				Context:   fmt.Sprintf("Segment with %s sentiment", emotion),
			})
		}
	}

	// Calculate overall sentiment
	if len(segments) > 0 {
		sentiment.SentimentScore = totalScore / float64(len(segments))
	}

	if sentiment.SentimentScore > 0.1 {
		sentiment.OverallSentiment = "positive"
	} else if sentiment.SentimentScore < -0.1 {
		sentiment.OverallSentiment = "negative"
	} else {
		sentiment.OverallSentiment = "neutral"
	}

	// Convert speaker sentiments map to slice
	for _, speakerData := range speakerSentiments {
		sentiment.SpeakerSentiments = append(sentiment.SpeakerSentiments, *speakerData)
	}

	return sentiment
}

// calculateMeetingMetrics provides overall meeting statistics
func (s *AnalyticsService) calculateMeetingMetrics(segments []entities.TranscriptSegment) *MeetingMetrics {
	if len(segments) == 0 {
		return &MeetingMetrics{}
	}

	metrics := &MeetingMetrics{
		SegmentCount: len(segments),
	}

	speakers := make(map[string]bool)
	totalConfidence := 0.0
	totalWords := 0
	activeSpeakingTime := 0.0
	lastEndTime := 0.0
	speakerTurns := 0
	lastSpeaker := ""

	// Calculate basic metrics
	for i, segment := range segments {
		duration := segment.EndTime - segment.StartTime
		activeSpeakingTime += duration

		words := len(strings.Fields(segment.Text))
		totalWords += words
		totalConfidence += segment.Confidence

		if segment.Speaker != "" {
			speakers[segment.Speaker] = true

			// Track speaker turns
			if i > 0 && segment.Speaker != lastSpeaker {
				speakerTurns++
			}
			lastSpeaker = segment.Speaker
		}

		if segment.EndTime > lastEndTime {
			lastEndTime = segment.EndTime
		}

		// Track longest monologue
		if segment.Speaker != "" && duration > metrics.LongestMonologue {
			metrics.LongestMonologue = duration
		}
	}

	metrics.TotalDuration = lastEndTime
	metrics.ActiveSpeakingTime = activeSpeakingTime
	metrics.SilenceTime = metrics.TotalDuration - activeSpeakingTime
	metrics.SpeakerCount = len(speakers)
	metrics.WordCount = totalWords
	metrics.AverageConfidence = totalConfidence / float64(len(segments))

	// Calculate speaker turnover rate (turns per minute)
	if metrics.TotalDuration > 0 {
		metrics.SpeakerTurnover = float64(speakerTurns) / (metrics.TotalDuration / 60.0)
	}

	// Determine overall pace
	if metrics.TotalDuration > 0 {
		wordsPerMinute := float64(totalWords) / (metrics.TotalDuration / 60.0)
		if wordsPerMinute < 100 {
			metrics.OverallPace = "slow"
		} else if wordsPerMinute > 160 {
			metrics.OverallPace = "fast"
		} else {
			metrics.OverallPace = "normal"
		}
	}

	// Estimate interruptions (simplified - when segments overlap or are very short)
	for i := 1; i < len(segments); i++ {
		prevSegment := segments[i-1]
		currSegment := segments[i]

		// If current segment starts before previous ends, it might be an interruption
		if currSegment.StartTime < prevSegment.EndTime {
			metrics.InterruptionCount++
		}

		// Very short segments might indicate interruptions
		if currSegment.EndTime-currSegment.StartTime < 1.0 && len(strings.Fields(currSegment.Text)) < 3 {
			metrics.InterruptionCount++
		}
	}

	return metrics
}

// analyzeKeywords extracts and analyzes keyword frequency
func (s *AnalyticsService) analyzeKeywords(segments []entities.TranscriptSegment) []KeywordFrequency {
	wordCount := make(map[string]int)
	wordContext := make(map[string][]string)
	totalWords := 0

	// Common stop words to exclude
	stopWords := map[string]bool{
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true, "to": true, "for": true, "of": true, "with": true, "by": true,
		"a": true, "an": true, "as": true, "are": true, "was": true, "is": true, "been": true, "be": true, "have": true, "has": true, "had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true, "may": true, "might": true, "can": true, "this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "it": true, "we": true, "they": true, "me": true, "him": true, "her": true, "us": true, "them": true,
		"my": true, "your": true, "his": true, "hers": true, "its": true, "our": true, "their": true, "yes": true, "no": true, "not": true, "so": true, "very": true,
		"just": true, "now": true, "then": true, "here": true, "there": true, "where": true, "when": true, "why": true, "how": true, "what": true, "who": true,
		"um": true, "uh": true, "yeah": true, "okay": true, "ok": true, "well": true, "like": true, "know": true, "think": true, "get": true, "go": true, "see": true,
	}

	// Extract words
	for _, segment := range segments {
		words := strings.Fields(strings.ToLower(segment.Text))
		totalWords += len(words)

		for _, word := range words {
			// Clean word (remove punctuation)
			cleaned := strings.Trim(word, ".,!?;:\"'()[]{}/-")

			// Skip short words, numbers, and stop words
			if len(cleaned) < 3 || stopWords[cleaned] {
				continue
			}

			wordCount[cleaned]++

			// Store context for top words
			if len(wordContext[cleaned]) < 3 {
				wordContext[cleaned] = append(wordContext[cleaned], segment.Text)
			}
		}
	}

	// Convert to frequency analysis
	var keywords []KeywordFrequency
	for word, count := range wordCount {
		if count > 1 { // Only include words mentioned more than once
			frequency := float64(count) / float64(totalWords)

			// Simple TF-IDF calculation (simplified)
			tf := frequency
			idf := math.Log(float64(len(segments)) / float64(count))
			tfidf := tf * idf

			keywords = append(keywords, KeywordFrequency{
				Keyword:   word,
				Count:     count,
				Frequency: frequency,
				TfIdf:     tfidf,
				Context:   wordContext[word],
			})
		}
	}

	// Sort by count
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Count > keywords[j].Count
	})

	// Limit to top 20 keywords
	if len(keywords) > 20 {
		keywords = keywords[:20]
	}

	return keywords
}

// analyzeTimeDistribution shows activity distribution over time
func (s *AnalyticsService) analyzeTimeDistribution(segments []entities.TranscriptSegment) *TimeDistribution {
	if len(segments) == 0 {
		return &TimeDistribution{}
	}

	totalDuration := s.calculateTotalDuration(segments)
	slotDuration := 30.0 // 30-second slots
	numSlots := int(math.Ceil(totalDuration / slotDuration))

	timeSlots := make([]TimeSlot, numSlots)

	// Initialize time slots
	for i := 0; i < numSlots; i++ {
		timeSlots[i] = TimeSlot{
			StartTime: float64(i) * slotDuration,
			EndTime:   float64(i+1) * slotDuration,
		}
	}

	// Populate time slots with activity
	for _, segment := range segments {
		startSlot := int(segment.StartTime / slotDuration)
		endSlot := int(segment.EndTime / slotDuration)

		for slot := startSlot; slot <= endSlot && slot < numSlots; slot++ {
			timeSlots[slot].WordCount += len(strings.Fields(segment.Text))

			// Track unique speakers in this slot
			speakerMap := make(map[string]bool)
			speakerMap[segment.Speaker] = true
			timeSlots[slot].SpeakerCount = len(speakerMap)
		}
	}

	// Determine activity levels and create activity curve
	activityCurve := make([]float64, numSlots)
	maxWords := 0

	for _, slot := range timeSlots {
		if slot.WordCount > maxWords {
			maxWords = slot.WordCount
		}
	}

	var peakPeriods []Period
	var quietPeriods []Period

	for slotIndex, slot := range timeSlots {
		// Normalize activity
		if maxWords > 0 {
			activityCurve[slotIndex] = float64(slot.WordCount) / float64(maxWords)
		}

		// Determine activity level
		if slot.WordCount == 0 {
			timeSlots[slotIndex].ActivityLevel = "low"
		} else {
			ratio := float64(slot.WordCount) / float64(maxWords)
			if ratio > 0.7 {
				timeSlots[slotIndex].ActivityLevel = "high"

				// Track peak periods
				peakPeriods = append(peakPeriods, Period{
					StartTime: slot.StartTime,
					EndTime:   slot.EndTime,
					Duration:  slotDuration,
				})
			} else if ratio > 0.3 {
				timeSlots[slotIndex].ActivityLevel = "medium"
			} else {
				timeSlots[slotIndex].ActivityLevel = "low"

				// Track quiet periods
				quietPeriods = append(quietPeriods, Period{
					StartTime: slot.StartTime,
					EndTime:   slot.EndTime,
					Duration:  slotDuration,
				})
			}
		}
	}

	return &TimeDistribution{
		TimeSlots:     timeSlots,
		PeakActivity:  peakPeriods,
		QuietPeriods:  quietPeriods,
		ActivityCurve: activityCurve,
	}
}

// assessQuality evaluates transcription quality
func (s *AnalyticsService) assessQuality(segments []entities.TranscriptSegment) *QualityMetrics {
	if len(segments) == 0 {
		return &QualityMetrics{
			OverallQuality: "poor",
			AudioQuality:   "unknown",
		}
	}

	totalConfidence := 0.0
	lowConfidenceCount := 0
	confidenceDistribution := map[string]int{
		"excellent": 0, // 0.9+
		"good":      0, // 0.7-0.9
		"fair":      0, // 0.5-0.7
		"poor":      0, // <0.5
	}

	for _, segment := range segments {
		totalConfidence += segment.Confidence

		if segment.Confidence < 0.6 {
			lowConfidenceCount++
		}

		// Categorize confidence
		if segment.Confidence >= 0.9 {
			confidenceDistribution["excellent"]++
		} else if segment.Confidence >= 0.7 {
			confidenceDistribution["good"]++
		} else if segment.Confidence >= 0.5 {
			confidenceDistribution["fair"]++
		} else {
			confidenceDistribution["poor"]++
		}
	}

	averageConfidence := totalConfidence / float64(len(segments))

	// Determine overall quality
	var overallQuality string
	if averageConfidence >= 0.9 {
		overallQuality = "excellent"
	} else if averageConfidence >= 0.75 {
		overallQuality = "good"
	} else if averageConfidence >= 0.6 {
		overallQuality = "fair"
	} else {
		overallQuality = "poor"
	}

	// Audio quality assessment based on confidence patterns
	audioQuality := "good"
	if lowConfidenceCount > len(segments)/3 {
		audioQuality = "poor"
	} else if lowConfidenceCount > len(segments)/10 {
		audioQuality = "fair"
	}

	// Generate recommendations
	var recommendations []string
	if averageConfidence < 0.7 {
		recommendations = append(recommendations, "Consider improving audio quality for better transcription accuracy")
	}
	if lowConfidenceCount > len(segments)/4 {
		recommendations = append(recommendations, "Review and manually correct low-confidence segments")
	}
	if audioQuality == "poor" {
		recommendations = append(recommendations, "Use better microphones or reduce background noise")
	}

	return &QualityMetrics{
		OverallQuality:         overallQuality,
		AverageConfidence:      averageConfidence,
		LowConfidenceSegments:  lowConfidenceCount,
		ConfidenceDistribution: confidenceDistribution,
		AudioQuality:           audioQuality,
		RecommendedActions:     recommendations,
	}
}

// generateInsights creates AI-generated insights
func (s *AnalyticsService) generateInsights(analytics *AnalyticsData) []Insight {
	var insights []Insight

	// Speaker behavior insights
	if len(analytics.SpeakerAnalytics) > 1 {
		maxSpeaker := analytics.SpeakerAnalytics[0]
		if maxSpeaker.SpeakingTimePercent > 60 {
			insights = append(insights, Insight{
				Type:        "speaker_behavior",
				Title:       "Dominant Speaker Detected",
				Description: fmt.Sprintf("%s spoke for %.1f%% of the meeting, which may indicate an imbalanced discussion", maxSpeaker.Speaker, maxSpeaker.SpeakingTimePercent),
				Confidence:  0.8,
				Severity:    "medium",
				Speaker:     maxSpeaker.Speaker,
				ActionItems: []string{"Encourage more participation from other attendees", "Consider time-boxing discussions"},
			})
		}
	}

	// Meeting pace insights
	if analytics.MeetingMetrics.OverallPace == "fast" {
		insights = append(insights, Insight{
			Type:        "meeting_pace",
			Title:       "Fast-Paced Discussion",
			Description: "The meeting had a fast pace, which might make it difficult for participants to follow",
			Confidence:  0.7,
			Severity:    "low",
			ActionItems: []string{"Consider slowing down key discussions", "Provide meeting summaries"},
		})
	}

	// Quality insights
	if analytics.QualityMetrics.OverallQuality == "poor" {
		insights = append(insights, Insight{
			Type:        "quality_issue",
			Title:       "Poor Audio Quality Detected",
			Description: fmt.Sprintf("Average confidence is %.2f, indicating audio quality issues", analytics.QualityMetrics.AverageConfidence),
			Confidence:  0.9,
			Severity:    "high",
			ActionItems: analytics.QualityMetrics.RecommendedActions,
		})
	}

	// Sentiment insights
	if analytics.SentimentAnalysis.OverallSentiment == "negative" {
		insights = append(insights, Insight{
			Type:        "sentiment_analysis",
			Title:       "Negative Sentiment Detected",
			Description: "The overall tone of the meeting was negative, which may indicate concerns or conflicts",
			Confidence:  0.6,
			Severity:    "medium",
			ActionItems: []string{"Follow up on concerns raised", "Address any conflicts or issues"},
		})
	}

	// Topic insights
	if len(analytics.TopicAnalysis) > 0 {
		topTopic := analytics.TopicAnalysis[0]
		insights = append(insights, Insight{
			Type:        "topic_analysis",
			Title:       "Primary Discussion Topic",
			Description: fmt.Sprintf("'%s' was the main topic discussed, with %.2f relevance score", topTopic.Topic, topTopic.Relevance),
			Confidence:  0.8,
			Severity:    "low",
		})
	}

	// Interruption insights
	if analytics.MeetingMetrics.InterruptionCount > 10 {
		insights = append(insights, Insight{
			Type:        "meeting_dynamics",
			Title:       "High Interruption Rate",
			Description: fmt.Sprintf("Detected %d interruptions, which may indicate heated discussions or poor meeting management", analytics.MeetingMetrics.InterruptionCount),
			Confidence:  0.7,
			Severity:    "medium",
			ActionItems: []string{"Implement better turn-taking protocols", "Use meeting facilitation techniques"},
		})
	}

	return insights
}

// Helper functions

func (s *AnalyticsService) calculateTotalDuration(segments []entities.TranscriptSegment) float64 {
	if len(segments) == 0 {
		return 0
	}

	maxEndTime := 0.0
	for _, segment := range segments {
		if segment.EndTime > maxEndTime {
			maxEndTime = segment.EndTime
		}
	}

	return maxEndTime
}

// GetMeetingAnalyticsSummary provides a high-level analytics summary for a meeting
func (s *AnalyticsService) GetMeetingAnalyticsSummary(ctx context.Context, meetingID string) (*AnalyticsData, error) {
	// Get all transcriptions for the meeting
	transcriptions, err := s.transcriptionRepo.FindByMeetingID(ctx, meetingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transcriptions for meeting: %w", err)
	}

	if len(transcriptions) == 0 {
		return nil, fmt.Errorf("no transcriptions found for meeting %s", meetingID)
	}

	// Use the latest transcription for analytics
	latestTranscription := transcriptions[0]
	for _, t := range transcriptions {
		if t.CreatedAt.After(latestTranscription.CreatedAt) {
			latestTranscription = t
		}
	}

	return s.GetTranscriptionAnalytics(ctx, latestTranscription.ID)
}
