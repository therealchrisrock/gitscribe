package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/therealchrisrock/assemblyai-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		log.Fatal("ASSEMBLYAI_API_KEY environment variable is required")
	}

	// Create a new client
	client := assemblyai.NewClient(apiKey)

	// Audio URL to transcribe
	audioURL := "https://github.com/AssemblyAI-Examples/audio-examples/raw/main/20230607_me_canadian_wildfires.mp3"

	fmt.Println("Starting advanced transcription with AI features...")

	// Configure advanced options using the fluent API
	options := assemblyai.NewTranscriptRequest(audioURL).
		WithSpeakerLabels(true).
		WithSentimentAnalysis(true).
		WithEntityDetection(true).
		WithAutoHighlights(true).
		WithAutoChapters(true).
		WithSummarization(true).
		WithLanguageCode("en")

	// Transcribe the audio
	transcript, err := client.TranscribeFromURL(
		context.Background(),
		audioURL,
		options,
	)
	if err != nil {
		log.Fatal("Transcription failed:", err)
	}

	// Print basic results
	fmt.Println("\n=== Basic Results ===")
	fmt.Printf("ID: %s\n", transcript.ID)
	fmt.Printf("Status: %s\n", transcript.Status)
	fmt.Printf("Audio Duration: %.2f seconds\n", *transcript.AudioDuration)
	fmt.Printf("Confidence: %.2f\n", *transcript.Confidence)

	// Print speaker labels
	if len(transcript.Utterances) > 0 {
		fmt.Println("\n=== Speaker Labels ===")
		for i, utterance := range transcript.Utterances {
			if i >= 5 { // Limit to first 5 utterances for brevity
				fmt.Printf("... and %d more utterances\n", len(transcript.Utterances)-5)
				break
			}
			fmt.Printf("Speaker %s: %s\n", utterance.Speaker, utterance.Text)
		}
	}

	// Print sentiment analysis
	if len(transcript.SentimentAnalysisResults) > 0 {
		fmt.Println("\n=== Sentiment Analysis ===")
		for i, sentiment := range transcript.SentimentAnalysisResults {
			if i >= 3 { // Limit to first 3 for brevity
				fmt.Printf("... and %d more sentiment results\n", len(transcript.SentimentAnalysisResults)-3)
				break
			}
			fmt.Printf("Text: \"%s\"\n", sentiment.Text)
			fmt.Printf("Sentiment: %s (%.2f confidence)\n\n", sentiment.Sentiment, sentiment.Confidence)
		}
	}

	// Print entities
	if len(transcript.Entities) > 0 {
		fmt.Println("\n=== Entities ===")
		for i, entity := range transcript.Entities {
			if i >= 10 { // Limit to first 10 for brevity
				fmt.Printf("... and %d more entities\n", len(transcript.Entities)-10)
				break
			}
			fmt.Printf("Entity: \"%s\" (%s)\n", entity.Text, entity.EntityType)
		}
	}

	// Print auto highlights
	if transcript.AutoHighlightsResult != nil && len(transcript.AutoHighlightsResult.Results) > 0 {
		fmt.Println("\n=== Auto Highlights ===")
		for i, highlight := range transcript.AutoHighlightsResult.Results {
			if i >= 5 { // Limit to first 5 for brevity
				fmt.Printf("... and %d more highlights\n", len(transcript.AutoHighlightsResult.Results)-5)
				break
			}
			fmt.Printf("Highlight: \"%s\" (rank: %.2f, count: %d)\n",
				highlight.Text, highlight.Rank, highlight.Count)
		}
	}

	// Print chapters
	if len(transcript.Chapters) > 0 {
		fmt.Println("\n=== Auto Chapters ===")
		for _, chapter := range transcript.Chapters {
			startTime := chapter.Start / 1000 // Convert to seconds
			endTime := chapter.End / 1000
			fmt.Printf("Chapter: %s (%d:%02d - %d:%02d)\n",
				chapter.Headline,
				startTime/60, startTime%60,
				endTime/60, endTime%60)
			fmt.Printf("Summary: %s\n\n", chapter.Summary)
		}
	}

	// Print summary
	if transcript.Summary != nil {
		fmt.Println("\n=== Summary ===")
		fmt.Printf("%s\n", *transcript.Summary)
	}

	// Print full transcript
	fmt.Println("\n=== Full Transcript ===")
	fmt.Printf("%s\n", *transcript.Text)
}
