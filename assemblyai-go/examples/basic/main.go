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

	fmt.Println("Starting transcription...")

	// Transcribe the audio
	transcript, err := client.TranscribeFromURL(
		context.Background(),
		audioURL,
		nil, // Use default options
	)
	if err != nil {
		log.Fatal("Transcription failed:", err)
	}

	// Print the results
	fmt.Println("\n=== Transcription Results ===")
	fmt.Printf("ID: %s\n", transcript.ID)
	fmt.Printf("Status: %s\n", transcript.Status)
	fmt.Printf("Audio Duration: %.2f seconds\n", *transcript.AudioDuration)
	fmt.Printf("Confidence: %.2f\n", *transcript.Confidence)
	fmt.Printf("\nText:\n%s\n", *transcript.Text)
}
