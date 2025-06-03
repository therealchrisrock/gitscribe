package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_firebase_url.go <firebase_url>")
		fmt.Println("Example:")
		fmt.Println("go run test_firebase_url.go 'https://storage.googleapis.com/...'")
		os.Exit(1)
	}

	firebaseURL := os.Args[1]
	fmt.Printf("Testing Firebase URL: %s\n", firebaseURL[:80]+"...")

	// Make HTTP request to Firebase URL
	resp, err := http.Get(firebaseURL)
	if err != nil {
		fmt.Printf("❌ Failed to fetch URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("✅ HTTP Status: %s\n", resp.Status)
	fmt.Printf("📄 Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("📊 Content-Length: %s\n", resp.Header.Get("Content-Length"))

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📦 Downloaded %d bytes\n", len(data))

	if len(data) > 0 {
		fmt.Printf("✅ Audio file downloaded successfully!\n")

		// Check if it starts with WAV header
		if len(data) >= 4 && string(data[:4]) == "RIFF" {
			fmt.Printf("🎵 File appears to be a WAV audio file (starts with RIFF header)\n")
		} else if len(data) >= 8 {
			fmt.Printf("🔍 File header (first 8 bytes): %x\n", data[:8])
		}

		// Save to file for verification
		filename := "downloaded_audio_test.wav"
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to save file: %v\n", err)
		} else {
			fmt.Printf("💾 Audio saved as: %s\n", filename)
		}
	} else {
		fmt.Printf("❌ No data received\n")
		os.Exit(1)
	}

	fmt.Printf("\n🎉 Firebase URL test completed successfully!\n")
}
