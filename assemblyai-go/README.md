# AssemblyAI Go SDK

A comprehensive Go client library for the AssemblyAI API. This package provides a simple and idiomatic way to interact with AssemblyAI's speech-to-text and audio intelligence services.

## Features

- **Speech-to-Text**: Transcribe audio files from URLs or local files
- **Audio Intelligence**: Speaker labels, sentiment analysis, entity detection, auto-highlights, and more
- **LeMUR**: Large Language Model powered analysis of transcripts
- **File Upload**: Upload local audio files to AssemblyAI
- **Polling**: Automatic polling for transcript completion
- **Error Handling**: Comprehensive error handling with detailed API errors
- **Context Support**: Full context.Context support for cancellation and timeouts
- **Fluent API**: Chainable methods for easy configuration

## Installation

```bash
go get github.com/therealchrisrock/assemblyai-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/therealchrisrock/assemblyai-go"
)

func main() {
    // Create a new client
    client := assemblyai.NewClient(os.Getenv("ASSEMBLYAI_API_KEY"))

    // Transcribe from URL
    transcript, err := client.TranscribeFromURL(
        context.Background(),
        "https://example.com/audio.mp3",
        nil, // Use default options
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Transcript:", *transcript.Text)
}
```

## Usage Examples

### Basic Transcription

#### From URL

```go
client := assemblyai.NewClient("your-api-key")

transcript, err := client.TranscribeFromURL(
    context.Background(),
    "https://example.com/audio.mp3",
    nil,
)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Text:", *transcript.Text)
fmt.Println("Confidence:", *transcript.Confidence)
```

#### From Local File

```go
file, err := os.Open("audio.mp3")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

transcript, err := client.TranscribeFromReader(
    context.Background(),
    file,
    nil,
)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Text:", *transcript.Text)
```

### Advanced Transcription Options

```go
// Using the fluent API
options := assemblyai.NewTranscriptRequest("https://example.com/audio.mp3").
    WithSpeakerLabels(true).
    WithSentimentAnalysis(true).
    WithEntityDetection(true).
    WithAutoHighlights(true).
    WithLanguageCode("en")

transcript, err := client.TranscribeFromURL(
    context.Background(),
    "https://example.com/audio.mp3",
    options,
)
if err != nil {
    log.Fatal(err)
}

// Access speaker labels
for _, utterance := range transcript.Utterances {
    fmt.Printf("Speaker %s: %s\n", utterance.Speaker, utterance.Text)
}

// Access sentiment analysis
for _, sentiment := range transcript.SentimentAnalysisResults {
    fmt.Printf("Sentiment: %s (%.2f confidence)\n", 
        sentiment.Sentiment, sentiment.Confidence)
}

// Access entities
for _, entity := range transcript.Entities {
    fmt.Printf("Entity: %s (%s)\n", entity.Text, entity.EntityType)
}
```

### Manual Transcript Management

```go
// Create transcript without waiting
request := &assemblyai.TranscriptRequest{
    AudioURL:      "https://example.com/audio.mp3",
    SpeakerLabels: assemblyai.Bool(true),
}

transcript, err := client.CreateTranscript(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Transcript ID:", transcript.ID)
fmt.Println("Status:", transcript.Status)

// Poll for completion manually
for transcript.Status == assemblyai.StatusQueued || transcript.Status == assemblyai.StatusProcessing {
    time.Sleep(5 * time.Second)
    
    transcript, err = client.GetTranscript(context.Background(), transcript.ID)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Status:", transcript.Status)
}

if transcript.Status == assemblyai.StatusCompleted {
    fmt.Println("Text:", *transcript.Text)
} else {
    fmt.Println("Error:", *transcript.Error)
}
```

### File Upload

```go
// Upload a file and get the upload URL
file, err := os.Open("audio.mp3")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

uploadResp, err := client.UploadFile(context.Background(), file)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Upload URL:", uploadResp.UploadURL)

// Use the upload URL for transcription
transcript, err := client.TranscribeFromURL(
    context.Background(),
    uploadResp.UploadURL,
    nil,
)
```

### LeMUR (Large Language Model) Features

#### Custom Task

```go
request := assemblyai.NewLemurRequest(
    []string{"transcript-id-1", "transcript-id-2"},
    "Summarize the key points discussed in this meeting.",
).WithTemperature(0.7)

response, err := client.LemurTask(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Response:", response.Response)
fmt.Println("Usage:", response.Usage)
```

#### Summary

```go
request := assemblyai.NewLemurSummaryRequest([]string{"transcript-id"}).
    WithAnswerFormat("bullet points")

response, err := client.LemurSummary(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Summary:", response.Response)
```

#### Question & Answer

```go
questions := []assemblyai.LemurQuestion{
    {
        Question: "What were the main topics discussed?",
        AnswerFormat: assemblyai.String("bullet points"),
    },
    {
        Question: "What action items were mentioned?",
    },
}

request := assemblyai.NewLemurQuestionAnswerRequest(
    []string{"transcript-id"},
    questions,
)

response, err := client.LemurQuestionAnswer(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

for _, answer := range response.Response {
    fmt.Printf("Q: %s\nA: %s\n\n", answer.Question, answer.Answer)
}
```

#### Action Items

```go
request := assemblyai.NewLemurActionItemsRequest([]string{"transcript-id"})

response, err := client.LemurActionItems(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Action Items:")
for _, item := range response.ActionItems {
    fmt.Println("-", item)
}
```

### List Transcripts

```go
// List recent transcripts
limit := 10
listResp, err := client.ListTranscripts(context.Background(), &limit, nil, nil)
if err != nil {
    log.Fatal(err)
}

for _, transcript := range listResp.Transcripts {
    fmt.Printf("ID: %s, Status: %s\n", transcript.ID, transcript.Status)
}

// Pagination
if listResp.PageDetails.NextURL != nil {
    fmt.Println("Next page available:", *listResp.PageDetails.NextURL)
}
```

### PII Redaction

```go
options := assemblyai.NewTranscriptRequest("https://example.com/audio.mp3").
    WithRedactPII(true, 
        assemblyai.PIIPolicyPersonName,
        assemblyai.PIIPolicyPhoneNumber,
        assemblyai.PIIPolicyEmailAddress,
    )

transcript, err := client.TranscribeFromURL(
    context.Background(),
    "https://example.com/audio.mp3",
    options,
)
```

### Custom Spelling

```go
customSpelling := []assemblyai.CustomSpelling{
    {
        From: []string{"AssemblyAI", "Assembly AI"},
        To:   "AssemblyAI",
    },
}

options := assemblyai.NewTranscriptRequest("https://example.com/audio.mp3").
    WithCustomSpelling(customSpelling)
```

## Configuration

### Client Options

```go
// Use EU server
client := assemblyai.NewClient("api-key", assemblyai.WithEUServer())

// Custom HTTP client
httpClient := &http.Client{Timeout: 60 * time.Second}
client := assemblyai.NewClient("api-key", assemblyai.WithHTTPClient(httpClient))

// Custom base URL
client := assemblyai.NewClient("api-key", assemblyai.WithBaseURL("https://custom.api.com"))
```

### Context and Timeouts

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

transcript, err := client.TranscribeFromURL(ctx, audioURL, nil)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(30 * time.Second)
    cancel() // Cancel after 30 seconds
}()

transcript, err := client.TranscribeFromURL(ctx, audioURL, nil)
```

## Error Handling

```go
transcript, err := client.TranscribeFromURL(ctx, audioURL, nil)
if err != nil {
    var apiErr *assemblyai.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.StatusCode)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

## Supported Audio Formats

AssemblyAI supports many audio formats including:
- MP3
- MP4
- WAV
- FLAC
- OGG
- M4A
- WebM
- And many more

## Rate Limits

The LeMUR endpoints have rate limits. Check the response headers for rate limit information:
- `X-RateLimit-Limit`: Maximum requests per minute
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Seconds until rate limit resets

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support with this SDK, please open an issue on GitHub.
For AssemblyAI API support, visit [AssemblyAI Support](https://www.assemblyai.com/support). 