package providers

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// FirebaseStorageUploader implements the FirebaseUploader interface
type FirebaseStorageUploader struct {
	client          *storage.Client
	bucketName      string
	credentialsPath string
}

// NewFirebaseStorageUploader creates a new Firebase storage uploader
func NewFirebaseStorageUploader(bucketName, credentialsPath string) (*FirebaseStorageUploader, error) {
	ctx := context.Background()

	// Create storage client with credentials
	opt := option.WithCredentialsFile(credentialsPath)
	storageClient, err := storage.NewClient(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &FirebaseStorageUploader{
		client:          storageClient,
		bucketName:      bucketName,
		credentialsPath: credentialsPath,
	}, nil
}

// UploadAudio uploads audio data to Firebase Storage
func (f *FirebaseStorageUploader) UploadAudio(ctx context.Context, audioData []byte, meetingID, sessionID string) (string, error) {
	// Generate file path
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("meetings/%s/audio/%s_%d.wav", meetingID, sessionID, timestamp)

	// Get bucket handle
	bucket := f.client.Bucket(f.bucketName)

	// Create object handle
	obj := bucket.Object(fileName)

	// Create writer
	w := obj.NewWriter(ctx)
	w.ContentType = "audio/wav"
	w.CacheControl = "public, max-age=86400" // Cache for 1 day

	// Set metadata
	w.Metadata = map[string]string{
		"meetingId":  meetingID,
		"sessionId":  sessionID,
		"uploadedAt": time.Now().Format(time.RFC3339),
	}

	// Write data
	if _, err := w.Write(audioData); err != nil {
		w.Close()
		return "", fmt.Errorf("failed to write audio data: %w", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Generate signed URL for access (valid for 1 hour)
	signedURL, err := bucket.SignedURL(fileName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		log.Printf("Warning: failed to generate signed URL, using public URL: %v", err)
		// Fallback to public URL format
		signedURL = fmt.Sprintf("gs://%s/%s", f.bucketName, fileName)
	}

	log.Printf("Audio uploaded to Firebase Storage: %s (size: %d bytes)", fileName, len(audioData))
	return signedURL, nil
}

// UploadAudioStream uploads audio from a reader to Firebase Storage
func (f *FirebaseStorageUploader) UploadAudioStream(ctx context.Context, reader io.Reader, meetingID, sessionID string) (string, error) {
	// Generate file path
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("meetings/%s/audio/%s_%d.wav", meetingID, sessionID, timestamp)

	// Get bucket handle
	bucket := f.client.Bucket(f.bucketName)

	// Create object handle
	obj := bucket.Object(fileName)

	// Create writer
	w := obj.NewWriter(ctx)
	w.ContentType = "audio/wav"
	w.CacheControl = "public, max-age=86400"

	// Set metadata
	w.Metadata = map[string]string{
		"meetingId":  meetingID,
		"sessionId":  sessionID,
		"uploadedAt": time.Now().Format(time.RFC3339),
	}

	// Copy from reader to writer
	bytesWritten, err := io.Copy(w, reader)
	if err != nil {
		w.Close()
		return "", fmt.Errorf("failed to copy audio stream: %w", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Generate signed URL for access (valid for 1 hour)
	signedURL, err := bucket.SignedURL(fileName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		log.Printf("Warning: failed to generate signed URL, using public URL: %v", err)
		// Fallback to public URL format
		signedURL = fmt.Sprintf("gs://%s/%s", f.bucketName, fileName)
	}

	log.Printf("Audio stream uploaded to Firebase Storage: %s (size: %d bytes)", fileName, bytesWritten)
	return signedURL, nil
}

// DeleteAudio deletes an audio file from Firebase Storage
func (f *FirebaseStorageUploader) DeleteAudio(ctx context.Context, filePath string) error {
	// Extract object name from path (remove gs://bucket/ prefix if present)
	objectName := filePath
	if filepath.HasPrefix(filePath, fmt.Sprintf("gs://%s/", f.bucketName)) {
		objectName = filePath[len(fmt.Sprintf("gs://%s/", f.bucketName)):]
	}

	// Get bucket handle
	bucket := f.client.Bucket(f.bucketName)

	// Delete object
	obj := bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete audio file: %w", err)
	}

	log.Printf("Audio file deleted from Firebase Storage: %s", objectName)
	return nil
}

// Close closes the Firebase storage client
func (f *FirebaseStorageUploader) Close() error {
	return f.client.Close()
}
