package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"teammate/server/modules/transcription/domain/entities"
)

// ExportService handles transcription exports in various formats
type ExportService struct {
	storageUploader StorageUploader
}

// StorageUploader interface for uploading export files
type StorageUploader interface {
	UploadDocument(ctx context.Context, data []byte, fileName, contentType string) (string, error)
}

// ExportOptions defines options for transcription export
type ExportOptions struct {
	Format            string                 `json:"format"`                  // pdf, docx, json, txt
	IncludeMetadata   bool                   `json:"include_metadata"`        // Include transcription metadata
	IncludeSpeakers   bool                   `json:"include_speakers"`        // Include speaker names
	IncludeTimestamps bool                   `json:"include_timestamps"`      // Include timestamps
	IncludeStats      bool                   `json:"include_stats"`           // Include statistics
	Template          string                 `json:"template,omitempty"`      // Export template
	CustomFields      map[string]interface{} `json:"custom_fields,omitempty"` // Custom fields
}

// ExportResult contains the result of an export operation
type ExportResult struct {
	DownloadURL string                 `json:"download_url"`
	FileName    string                 `json:"file_name"`
	FileSize    int64                  `json:"file_size"`
	Format      string                 `json:"format"`
	ExpiresAt   time.Time              `json:"expires_at"`
	GeneratedAt time.Time              `json:"generated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewExportService creates a new export service
func NewExportService(storageUploader StorageUploader) *ExportService {
	return &ExportService{
		storageUploader: storageUploader,
	}
}

// ExportTranscription exports a transcription in the specified format
func (s *ExportService) ExportTranscription(ctx context.Context, transcription *TranscriptionHistoryItem, options *ExportOptions) (*ExportResult, error) {
	switch strings.ToLower(options.Format) {
	case "json":
		return s.exportJSON(ctx, transcription, options)
	case "txt":
		return s.exportTXT(ctx, transcription, options)
	case "pdf":
		return s.exportPDF(ctx, transcription, options)
	case "docx":
		return s.exportDOCX(ctx, transcription, options)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// exportJSON exports transcription as JSON
func (s *ExportService) exportJSON(ctx context.Context, transcription *TranscriptionHistoryItem, options *ExportOptions) (*ExportResult, error) {
	exportData := s.buildExportData(transcription, options)

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fileName := fmt.Sprintf("transcription_%s_%s.json", transcription.ID, time.Now().Format("20060102_150405"))

	downloadURL, err := s.storageUploader.UploadDocument(ctx, jsonData, fileName, "application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to upload JSON export: %w", err)
	}

	return &ExportResult{
		DownloadURL: downloadURL,
		FileName:    fileName,
		FileSize:    int64(len(jsonData)),
		Format:      "json",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		GeneratedAt: time.Now(),
		Metadata: map[string]interface{}{
			"transcription_id": transcription.ID,
			"segments_count":   len(transcription.Segments),
		},
	}, nil
}

// exportTXT exports transcription as plain text
func (s *ExportService) exportTXT(ctx context.Context, transcription *TranscriptionHistoryItem, options *ExportOptions) (*ExportResult, error) {
	var buffer bytes.Buffer

	// Add header if metadata is included
	if options.IncludeMetadata {
		buffer.WriteString("TRANSCRIPTION EXPORT\n")
		buffer.WriteString("===================\n\n")
		buffer.WriteString(fmt.Sprintf("Meeting ID: %s\n", transcription.MeetingID))
		buffer.WriteString(fmt.Sprintf("Transcription ID: %s\n", transcription.ID))
		buffer.WriteString(fmt.Sprintf("Provider: %s\n", transcription.Provider))
		buffer.WriteString(fmt.Sprintf("Status: %s\n", transcription.Status))
		buffer.WriteString(fmt.Sprintf("Confidence: %.2f%%\n", transcription.Confidence*100))
		buffer.WriteString(fmt.Sprintf("Created: %s\n", transcription.CreatedAt.Format("2006-01-02 15:04:05")))
		buffer.WriteString("\n")
	}

	// Add content
	if transcription.Content != "" {
		buffer.WriteString("FULL TRANSCRIPT\n")
		buffer.WriteString("---------------\n")
		buffer.WriteString(transcription.Content)
		buffer.WriteString("\n\n")
	}

	// Add segments if available
	if len(transcription.Segments) > 0 {
		buffer.WriteString("DETAILED SEGMENTS\n")
		buffer.WriteString("-----------------\n\n")

		for i, segment := range transcription.Segments {
			if options.IncludeTimestamps {
				buffer.WriteString(fmt.Sprintf("[%s - %s] ",
					s.formatTime(segment.StartTime),
					s.formatTime(segment.EndTime)))
			}

			if options.IncludeSpeakers && segment.Speaker != "" {
				buffer.WriteString(fmt.Sprintf("%s: ", segment.Speaker))
			}

			buffer.WriteString(segment.Text)

			if options.IncludeMetadata {
				buffer.WriteString(fmt.Sprintf(" (Confidence: %.2f%%)", segment.Confidence*100))
			}

			buffer.WriteString("\n")

			if i < len(transcription.Segments)-1 {
				buffer.WriteString("\n")
			}
		}
	}

	// Add statistics if requested
	if options.IncludeStats && transcription.Stats != nil {
		buffer.WriteString("\n\nSTATISTICS\n")
		buffer.WriteString("----------\n")
		statsJSON, _ := json.MarshalIndent(transcription.Stats, "", "  ")
		buffer.Write(statsJSON)
	}

	txtData := buffer.Bytes()
	fileName := fmt.Sprintf("transcription_%s_%s.txt", transcription.ID, time.Now().Format("20060102_150405"))

	downloadURL, err := s.storageUploader.UploadDocument(ctx, txtData, fileName, "text/plain")
	if err != nil {
		return nil, fmt.Errorf("failed to upload TXT export: %w", err)
	}

	return &ExportResult{
		DownloadURL: downloadURL,
		FileName:    fileName,
		FileSize:    int64(len(txtData)),
		Format:      "txt",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		GeneratedAt: time.Now(),
	}, nil
}

// exportPDF exports transcription as PDF (simplified HTML-to-PDF approach)
func (s *ExportService) exportPDF(ctx context.Context, transcription *TranscriptionHistoryItem, options *ExportOptions) (*ExportResult, error) {
	// For now, we'll create an HTML representation and note that PDF generation would require additional libraries
	htmlContent := s.generateHTML(transcription, options)

	// In a full implementation, you'd use a library like wkhtmltopdf, chromedp, or similar
	// For this example, we'll save as HTML with PDF content type
	fileName := fmt.Sprintf("transcription_%s_%s.html", transcription.ID, time.Now().Format("20060102_150405"))

	downloadURL, err := s.storageUploader.UploadDocument(ctx, []byte(htmlContent), fileName, "text/html")
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF export: %w", err)
	}

	return &ExportResult{
		DownloadURL: downloadURL,
		FileName:    fileName,
		FileSize:    int64(len(htmlContent)),
		Format:      "pdf",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		GeneratedAt: time.Now(),
		Metadata: map[string]interface{}{
			"note": "PDF generation requires additional libraries - this is an HTML representation",
		},
	}, nil
}

// exportDOCX exports transcription as Word document (simplified)
func (s *ExportService) exportDOCX(ctx context.Context, transcription *TranscriptionHistoryItem, options *ExportOptions) (*ExportResult, error) {
	// For DOCX generation, you'd typically use a library like github.com/unidoc/unioffice
	// For this example, we'll create a structured text representation
	content := s.generateWordContent(transcription, options)

	fileName := fmt.Sprintf("transcription_%s_%s.docx", transcription.ID, time.Now().Format("20060102_150405"))

	downloadURL, err := s.storageUploader.UploadDocument(ctx, []byte(content), fileName, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		return nil, fmt.Errorf("failed to upload DOCX export: %w", err)
	}

	return &ExportResult{
		DownloadURL: downloadURL,
		FileName:    fileName,
		FileSize:    int64(len(content)),
		Format:      "docx",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		GeneratedAt: time.Now(),
		Metadata: map[string]interface{}{
			"note": "DOCX generation requires additional libraries - this is a text representation",
		},
	}, nil
}

// buildExportData creates the data structure for export
func (s *ExportService) buildExportData(transcription *TranscriptionHistoryItem, options *ExportOptions) map[string]interface{} {
	data := map[string]interface{}{
		"transcription_id": transcription.ID,
		"content":          transcription.Content,
		"format_version":   "1.0",
		"exported_at":      time.Now().Format(time.RFC3339),
	}

	if options.IncludeMetadata {
		data["metadata"] = map[string]interface{}{
			"meeting_id": transcription.MeetingID,
			"provider":   transcription.Provider,
			"status":     transcription.Status,
			"confidence": transcription.Confidence,
			"created_at": transcription.CreatedAt,
			"updated_at": transcription.UpdatedAt,
		}
	}

	if len(transcription.Segments) > 0 {
		segments := make([]map[string]interface{}, len(transcription.Segments))
		for i, segment := range transcription.Segments {
			segmentData := map[string]interface{}{
				"sequence":   segment.SequenceNumber,
				"text":       segment.Text,
				"confidence": segment.Confidence,
			}

			if options.IncludeTimestamps {
				segmentData["start_time"] = segment.StartTime
				segmentData["end_time"] = segment.EndTime
			}

			if options.IncludeSpeakers {
				segmentData["speaker"] = segment.Speaker
			}

			segments[i] = segmentData
		}
		data["segments"] = segments
	}

	if options.IncludeStats && transcription.Stats != nil {
		data["statistics"] = transcription.Stats
	}

	// Add custom fields if provided
	if len(options.CustomFields) > 0 {
		data["custom_fields"] = options.CustomFields
	}

	return data
}

// generateHTML creates HTML content for PDF export
func (s *ExportService) generateHTML(transcription *TranscriptionHistoryItem, options *ExportOptions) string {
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Transcription Export</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { border-bottom: 2px solid #333; margin-bottom: 20px; }
        .metadata { background-color: #f5f5f5; padding: 15px; margin-bottom: 20px; }
        .segment { margin-bottom: 15px; }
        .timestamp { color: #666; font-size: 0.9em; }
        .speaker { font-weight: bold; color: #0066cc; }
        .confidence { color: #999; font-size: 0.8em; }
    </style>
</head>
<body>`)

	html.WriteString("<div class='header'><h1>Transcription Export</h1></div>")

	if options.IncludeMetadata {
		html.WriteString("<div class='metadata'>")
		html.WriteString(fmt.Sprintf("<p><strong>Meeting ID:</strong> %s</p>", transcription.MeetingID))
		html.WriteString(fmt.Sprintf("<p><strong>Transcription ID:</strong> %s</p>", transcription.ID))
		html.WriteString(fmt.Sprintf("<p><strong>Provider:</strong> %s</p>", transcription.Provider))
		html.WriteString(fmt.Sprintf("<p><strong>Confidence:</strong> %.2f%%</p>", transcription.Confidence*100))
		html.WriteString(fmt.Sprintf("<p><strong>Created:</strong> %s</p>", transcription.CreatedAt.Format("2006-01-02 15:04:05")))
		html.WriteString("</div>")
	}

	if transcription.Content != "" {
		html.WriteString("<div class='content'>")
		html.WriteString("<h2>Full Transcript</h2>")
		html.WriteString(fmt.Sprintf("<p>%s</p>", transcription.Content))
		html.WriteString("</div>")
	}

	if len(transcription.Segments) > 0 {
		html.WriteString("<div class='segments'>")
		html.WriteString("<h2>Detailed Segments</h2>")

		for _, segment := range transcription.Segments {
			html.WriteString("<div class='segment'>")

			if options.IncludeTimestamps {
				html.WriteString(fmt.Sprintf("<span class='timestamp'>[%s - %s]</span> ",
					s.formatTime(segment.StartTime),
					s.formatTime(segment.EndTime)))
			}

			if options.IncludeSpeakers && segment.Speaker != "" {
				html.WriteString(fmt.Sprintf("<span class='speaker'>%s:</span> ", segment.Speaker))
			}

			html.WriteString(segment.Text)

			if options.IncludeMetadata {
				html.WriteString(fmt.Sprintf(" <span class='confidence'>(%.2f%%)</span>", segment.Confidence*100))
			}

			html.WriteString("</div>")
		}
		html.WriteString("</div>")
	}

	html.WriteString("</body></html>")
	return html.String()
}

// generateWordContent creates content for DOCX export
func (s *ExportService) generateWordContent(transcription *TranscriptionHistoryItem, options *ExportOptions) string {
	// This would be replaced with actual DOCX generation in a full implementation
	return fmt.Sprintf(`TRANSCRIPTION EXPORT

Meeting ID: %s
Transcription ID: %s
Provider: %s
Created: %s

CONTENT:
%s

SEGMENTS:
%s`,
		transcription.MeetingID,
		transcription.ID,
		transcription.Provider,
		transcription.CreatedAt.Format("2006-01-02 15:04:05"),
		transcription.Content,
		s.formatSegmentsForText(transcription.Segments, options))
}

// formatSegmentsForText formats segments as text
func (s *ExportService) formatSegmentsForText(segments []entities.TranscriptSegment, options *ExportOptions) string {
	var buffer strings.Builder

	for _, segment := range segments {
		if options.IncludeTimestamps {
			buffer.WriteString(fmt.Sprintf("[%s - %s] ",
				s.formatTime(segment.StartTime),
				s.formatTime(segment.EndTime)))
		}

		if options.IncludeSpeakers && segment.Speaker != "" {
			buffer.WriteString(fmt.Sprintf("%s: ", segment.Speaker))
		}

		buffer.WriteString(segment.Text)
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// formatTime formats duration in seconds to MM:SS format
func (s *ExportService) formatTime(seconds float64) string {
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

// GetSupportedFormats returns list of supported export formats
func (s *ExportService) GetSupportedFormats() []string {
	return []string{"json", "txt", "pdf", "docx"}
}

// ValidateExportOptions validates export options
func (s *ExportService) ValidateExportOptions(options *ExportOptions) error {
	if options.Format == "" {
		return fmt.Errorf("format is required")
	}

	supportedFormats := s.GetSupportedFormats()
	format := strings.ToLower(options.Format)

	for _, supported := range supportedFormats {
		if format == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported format: %s (supported: %s)", options.Format, strings.Join(supportedFormats, ", "))
}
