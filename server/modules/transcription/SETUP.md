# Transcription Feature Setup

This document outlines how to set up the transcription feature with Firebase Storage and AssemblyAI integration.

## Phase 1: Core Provider Implementation - Complete ✅

### What's Implemented

1. **AssemblyAI Provider** (`infrastructure/providers/assemblyai_provider.go`)
   - Full integration with the standalone `assemblyai-go` package
   - Support for speaker diarization
   - Real-time and batch processing modes
   - Audio upload to Firebase Storage
   - Transcript polling and completion handling

2. **Firebase Storage Uploader** (`infrastructure/providers/firebase_uploader.go`)
   - Audio file upload to Firebase Storage
   - Signed URL generation for secure access
   - Stream-based uploads for large files
   - File deletion capabilities

3. **Mock Provider** (`infrastructure/providers/assemblyai_mock.go`)
   - Development/testing implementation
   - Realistic mock data with speaker diarization
   - Firebase Storage integration for testing

4. **Updated Factory** (`application/services/audio_processor_factory.go`)
   - Automatic provider selection based on configuration
   - Graceful fallback to mock when API keys are missing
   - Firebase Storage initialization

## Environment Configuration

Create a `.env` file in the `server` directory with the following variables:

```env
# AssemblyAI Configuration
ASSEMBLYAI_API_KEY=your-assemblyai-api-key-here

# Firebase Configuration
FIREBASE_STORAGE_BUCKET=your-firebase-storage-bucket
FIREBASE_CREDENTIALS_PATH=./firebase-credentials/serviceAccountKey.json

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/gitscribe_dev

# Server Configuration
PORT=8080
GIN_MODE=debug
```

## Firebase Setup

1. **Create Firebase Project**
   - Go to [Firebase Console](https://console.firebase.google.com/)
   - Create a new project or use existing one
   - Enable Storage in the project

2. **Generate Service Account Key**
   - Go to Project Settings > Service Accounts
   - Generate new private key
   - Download the JSON file as `serviceAccountKey.json`
   - Place it in `server/firebase-credentials/serviceAccountKey.json`

3. **Configure Storage Bucket**
   - Note your storage bucket name (usually `project-id.appspot.com`)
   - Set the `FIREBASE_STORAGE_BUCKET` environment variable

## AssemblyAI Setup

1. **Get API Key**
   - Sign up at [AssemblyAI](https://www.assemblyai.com/)
   - Get your API key from the dashboard
   - Set the `ASSEMBLYAI_API_KEY` environment variable

2. **Features Enabled**
   - Speaker diarization (speaker identification)
   - Sentiment analysis
   - Entity detection
   - Auto highlights
   - Punctuation and formatting

## Testing the Integration

### Using Mock Provider (No API Keys Required)

```bash
# Start the server without API keys
cd server
go run main.go
```

The system will automatically use the mock provider and log:
```
Warning: ASSEMBLYAI_API_KEY not set, using mock provider
Warning: Failed to initialize Firebase storage uploader: ...
```

### Using Real Providers

```bash
# Set environment variables
export ASSEMBLYAI_API_KEY="your-key"
export FIREBASE_STORAGE_BUCKET="your-bucket"
export FIREBASE_CREDENTIALS_PATH="./firebase-credentials/serviceAccountKey.json"

# Start the server
cd server
go run main.go
```

### WebSocket Testing

Connect to the WebSocket endpoint with speaker diarization enabled:

```
ws://localhost:8080/ws/enhanced-audio?provider=assemblyai&speaker_diarization=true&mode=batch
```