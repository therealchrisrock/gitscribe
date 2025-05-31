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

	// First, we need to transcribe some audio to get transcript IDs
	audioURL := "https://github.com/AssemblyAI-Examples/audio-examples/raw/main/20230607_me_canadian_wildfires.mp3"

	fmt.Println("Creating transcript for LeMUR analysis...")

	// Create a basic transcript
	transcript, err := client.TranscribeFromURL(
		context.Background(),
		audioURL,
		nil,
	)
	if err != nil {
		log.Fatal("Transcription failed:", err)
	}

	fmt.Printf("Transcript created with ID: %s\n", transcript.ID)

	// Now use LeMUR for various AI-powered analyses
	transcriptIDs := []string{transcript.ID}

	// 1. Custom Task
	fmt.Println("\n=== LeMUR Custom Task ===")
	taskRequest := assemblyai.NewLemurRequest(
		transcriptIDs,
		"Provide a detailed summary of the main topics discussed in this audio, including any key facts or statistics mentioned.",
	).WithTemperature(0.3)

	taskResponse, err := client.LemurTask(context.Background(), taskRequest)
	if err != nil {
		log.Printf("LeMUR task failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", taskResponse.Response)
		fmt.Printf("Usage: %d input tokens, %d output tokens\n",
			taskResponse.Usage.InputTokens, taskResponse.Usage.OutputTokens)
	}

	// 2. Summary
	fmt.Println("\n=== LeMUR Summary ===")
	summaryRequest := assemblyai.NewLemurSummaryRequest(transcriptIDs).
		WithAnswerFormat("bullet points")

	summaryResponse, err := client.LemurSummary(context.Background(), summaryRequest)
	if err != nil {
		log.Printf("LeMUR summary failed: %v", err)
	} else {
		fmt.Printf("Summary: %s\n", summaryResponse.Response)
	}

	// 3. Question & Answer
	fmt.Println("\n=== LeMUR Question & Answer ===")
	questions := []assemblyai.LemurQuestion{
		{
			Question:     "What is the main topic of this audio?",
			AnswerFormat: assemblyai.String("one sentence"),
		},
		{
			Question: "Are there any specific locations mentioned?",
		},
		{
			Question: "What are the key concerns or issues discussed?",
		},
	}

	qaRequest := assemblyai.NewLemurQuestionAnswerRequest(transcriptIDs, questions)

	qaResponse, err := client.LemurQuestionAnswer(context.Background(), qaRequest)
	if err != nil {
		log.Printf("LeMUR Q&A failed: %v", err)
	} else {
		for _, answer := range qaResponse.Response {
			fmt.Printf("Q: %s\n", answer.Question)
			fmt.Printf("A: %s\n\n", answer.Answer)
		}
	}

	// 4. Action Items
	fmt.Println("\n=== LeMUR Action Items ===")
	actionItemsRequest := assemblyai.NewLemurActionItemsRequest(transcriptIDs)

	actionItemsResponse, err := client.LemurActionItems(context.Background(), actionItemsRequest)
	if err != nil {
		log.Printf("LeMUR action items failed: %v", err)
	} else {
		if len(actionItemsResponse.ActionItems) > 0 {
			fmt.Println("Action Items:")
			for _, item := range actionItemsResponse.ActionItems {
				fmt.Printf("- %s\n", item)
			}
		} else {
			fmt.Println("No action items found in this audio.")
		}
	}

	fmt.Println("\n=== LeMUR Analysis Complete ===")
}
