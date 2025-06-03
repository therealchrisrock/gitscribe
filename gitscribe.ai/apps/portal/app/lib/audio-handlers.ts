/**
 * Abstract interface for handling audio data during recording
 */
export interface AudioHandler {
    initialize(): Promise<void>;
    handleAudioChunk(data: Blob, timestamp: string): Promise<void>;
    finalize(): Promise<void>;
    cleanup(): void;
    isConnected(): boolean;
    setOnAudioReady?(callback: (audioData: AudioPlaybackData) => void): void; // New callback for audio ready
}

export interface AudioPlaybackData {
    audioUrl: string;
    audioBlob: Blob;
    mimeType: string;
    chunkCount: number;
    totalSize: number;
    duration?: number;
}

// Message types for WebSocket communication with server
interface AudioMessage {
    type: string;
    session_id?: string;
    metadata?: AudioStreamMetadata;
    options?: AudioProcessingOptions;
    chunk?: AudioChunk;
    error?: string;
    result?: any;
    message?: string;
}

interface AudioStreamMetadata {
    session_id: string;
    meeting_id: string;
    user_id: string;
    sample_rate: number;
    channels: number;
    bits_per_sample: number;
    mime_type: string;
    start_time: number;
    mode: string;
}

interface AudioProcessingOptions {
    mode: string;
    provider: string;
    language: string;
    cost_optimized: boolean;
    real_time_transcription: boolean;
    speaker_diarization: boolean;
    confidence_threshold: number;
}

interface AudioChunk {
    data: number[]; // Convert Uint8Array to number array for JSON
    timestamp: number;
    sequence_num: number;
    size: number;
    duration?: number;
}

/**
 * Production audio handler that streams audio chunks via WebSocket
 */
export class WebSocketAudioHandler implements AudioHandler {
    private websocket: WebSocket | null = null;
    private totalBytesSent = 0;
    private chunkCount = 0;
    private sessionId: string = '';
    private isSessionStarted = false;
    private meetingId: string = '';
    private userId: string = '';

    constructor(
        private websocketUrl: string,
        private userInfo?: { userId: string; meetingId?: string }
    ) {
        // Generate IDs if not provided
        this.userId = userInfo?.userId || `user_${Math.random().toString(36).substr(2, 9)}`;
        this.meetingId = userInfo?.meetingId || `meeting_${Date.now()}`;
    }

    async initialize(): Promise<void> {
        return new Promise((resolve, reject) => {
            try {
                this.websocket = new WebSocket(this.websocketUrl);

                this.websocket.onopen = () => {
                    console.log("WebSocket connected");
                    this.setupMessageHandlers();
                    resolve();
                };

                this.websocket.onclose = (event) => {
                    console.log("WebSocket disconnected", event.code, event.reason);
                    this.websocket = null;
                    this.isSessionStarted = false;
                };

                this.websocket.onerror = (error) => {
                    console.error("WebSocket error:", error);
                    reject(new Error("Failed to connect to audio streaming server"));
                };
            } catch (err) {
                console.error("Error creating WebSocket:", err);
                reject(new Error("Failed to establish WebSocket connection"));
            }
        });
    }

    private setupMessageHandlers(): void {
        if (!this.websocket) return;

        this.websocket.onmessage = (event) => {
            try {
                const message: AudioMessage = JSON.parse(event.data);
                console.log("📨 Received WebSocket message:", message.type, message.message);

                switch (message.type) {
                    case 'connected':
                        console.log("✅ WebSocket connected and ready");
                        this.startSession();
                        break;
                    case 'connection_established':
                        console.log("✅ Enhanced WebSocket connection established");
                        this.startSession();
                        break;
                    case 'session_started':
                        console.log("✅ Recording session started");
                        this.isSessionStarted = true;
                        this.sessionId = message.session_id || '';
                        break;
                    case 'chunk_processed':
                        console.log("📄 Audio chunk processed");
                        break;
                    case 'chunk_buffered':
                        console.log("📦 Audio chunk buffered (session not started yet)");
                        break;
                    case 'buffered_chunks_processed':
                        console.log("✅ Buffered chunks processed successfully");
                        break;
                    case 'session_ended':
                        console.log("🎯 Recording session ended successfully");
                        if (message.result) {
                            console.log("📄 Final transcription result:", message.result);
                        }
                        this.isSessionStarted = false;
                        break;
                    case 'error':
                        console.error("❌ Server error:", message.error);
                        // Don't reset session on error, let user handle it
                        break;
                    case 'processing_result':
                        console.log("📄 Transcription result:", message.result);
                        break;
                    default:
                        console.log("📨 Unknown message type:", message.type, message);
                }
            } catch (err) {
                console.error("Failed to parse WebSocket message:", err, event.data);
            }
        };
    }

    private startSession(): void {
        if (!this.websocket || this.websocket.readyState !== WebSocket.OPEN) {
            console.error("Cannot start session: WebSocket not connected");
            return;
        }

        const sessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        const startMessage: AudioMessage = {
            type: "start_session",
            session_id: sessionId,
            metadata: {
                session_id: sessionId,
                meeting_id: this.meetingId,
                user_id: this.userId,
                sample_rate: 44100,
                channels: 1,
                bits_per_sample: 16,
                mime_type: "audio/webm",
                start_time: Date.now(),
                mode: "batch"
            }
        };

        console.log("🚀 Starting transcription session:", sessionId, "for user:", this.userId);
        this.websocket.send(JSON.stringify(startMessage));
        this.sessionId = sessionId;
    }

    async handleAudioChunk(data: Blob, timestamp: string): Promise<void> {
        if (!this.websocket || this.websocket.readyState !== WebSocket.OPEN) {
            console.warn("WebSocket not open, skipping chunk. State:", this.websocket?.readyState);
            return;
        }

        try {
            const buffer = await data.arrayBuffer();
            const uint8Array = new Uint8Array(buffer);
            this.chunkCount++;
            this.totalBytesSent += buffer.byteLength;

            const audioChunk: AudioChunk = {
                data: Array.from(uint8Array), // Convert to number array for JSON
                timestamp: Date.now(),
                sequence_num: this.chunkCount,
                size: buffer.byteLength,
                duration: 0.1 // Assume 100ms chunks
            };

            const chunkMessage: AudioMessage = {
                type: "audio_chunk",
                session_id: this.sessionId,
                chunk: audioChunk
            };

            console.log(`[${timestamp}] Sending audio chunk #${this.chunkCount}:`, {
                chunkSize: buffer.byteLength,
                totalSent: this.totalBytesSent,
                sessionStarted: this.isSessionStarted,
                wsState: this.websocket.readyState
            });

            // Double-check WebSocket is still open before sending
            if (this.websocket.readyState === WebSocket.OPEN) {
                this.websocket.send(JSON.stringify(chunkMessage));
            } else {
                console.warn("WebSocket closed during chunk processing");
            }
        } catch (error) {
            console.error("Error sending audio data:", error);
            // Don't throw error to prevent recording from failing completely
        }
    }

    async finalize(): Promise<void> {
        console.log(`📊 WebSocket streaming complete: ${this.chunkCount} chunks, ${this.totalBytesSent} total bytes`);

        // Check WebSocket state before attempting to send end message
        if (!this.websocket) {
            console.warn("Cannot end session: WebSocket is null");
            return;
        }

        if (this.websocket.readyState !== WebSocket.OPEN) {
            console.warn(`Cannot end session: WebSocket not open (state: ${this.websocket.readyState})`);
            return;
        }

        if (!this.isSessionStarted) {
            console.warn("Cannot end session: Session was never started");
            return;
        }

        try {
            const endMessage: AudioMessage = {
                type: "end_session",
                session_id: this.sessionId
            };

            console.log("🛑 Ending transcription session:", this.sessionId);

            // Final check before sending
            if (this.websocket.readyState === WebSocket.OPEN) {
                this.websocket.send(JSON.stringify(endMessage));

                // Wait a bit for the server to process the end session
                await new Promise(resolve => setTimeout(resolve, 500));
            } else {
                console.warn("WebSocket closed before sending end message");
            }
        } catch (error) {
            console.error("Error ending session:", error);
            // Don't throw error - just log it
        }
    }

    cleanup(): void {
        if (this.websocket) {
            this.websocket.close();
            this.websocket = null;
        }
        this.isSessionStarted = false;
        this.sessionId = '';
    }

    isConnected(): boolean {
        return this.websocket?.readyState === WebSocket.OPEN;
    }

    setOnAudioReady(callback: (audioData: AudioPlaybackData) => void): void {
        // No-op for production mode - only test mode needs this
        console.log("🚫 Audio ready callback not available in production mode");
    }
}

/**
 * Test mode audio handler that records audio locally and provides playback
 */
export class LocalRecordingAudioHandler implements AudioHandler {
    private audioChunks: Blob[] = [];
    private totalSize = 0;
    private chunkCount = 0;
    private mimeType = '';
    private onAudioReadyCallback?: (audioData: AudioPlaybackData) => void;

    async initialize(): Promise<void> {
        console.log("🧪 TEST MODE: Local recording initialized");
        this.audioChunks = [];
        this.totalSize = 0;
        this.chunkCount = 0;
        this.onAudioReadyCallback = undefined;
    }

    setOnAudioReady(callback: (audioData: AudioPlaybackData) => void): void {
        this.onAudioReadyCallback = callback;
        console.log("🎬 Audio ready callback set");
    }

    async handleAudioChunk(data: Blob, timestamp: string): Promise<void> {
        this.audioChunks.push(data);
        this.chunkCount++;
        this.totalSize += data.size;
        this.mimeType = data.type;

        console.log(`🎵 TEST MODE: Audio chunk #${this.chunkCount} recorded:`, {
            chunkNumber: this.chunkCount,
            size: data.size,
            type: data.type,
            timestamp: timestamp,
            totalSize: this.totalSize
        });
    }

    async finalize(): Promise<void> {
        if (this.audioChunks.length === 0) {
            console.warn("No audio chunks recorded");
            return;
        }

        console.log(`🎬 TEST MODE: Creating final audio file from ${this.chunkCount} chunks (${this.totalSize} bytes)`);

        // Combine all chunks into a single blob
        const finalBlob = new Blob(this.audioChunks, { type: this.mimeType });

        // Create object URL for the audio
        const audioUrl = URL.createObjectURL(finalBlob);

        // If we have a callback, use it (this is the new widget approach)
        if (this.onAudioReadyCallback) {
            console.log("🎬 Calling audio ready callback for widget");
            const audioData: AudioPlaybackData = {
                audioUrl,
                audioBlob: finalBlob,
                mimeType: this.mimeType,
                chunkCount: this.chunkCount,
                totalSize: this.totalSize
            };
            this.onAudioReadyCallback(audioData);
            return;
        }

        // Fallback: create download link if no callback is set
        console.log("🎬 No callback set, creating download link as fallback");
        this.createDownloadLink(finalBlob);
    }

    private createDownloadLink(blob: Blob): void {
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `test-recording-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.webm`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        console.log("🎬 TEST MODE: Audio download initiated");
    }

    cleanup(): void {
        this.audioChunks = [];
        this.totalSize = 0;
        this.chunkCount = 0;
        this.onAudioReadyCallback = undefined;
        console.log("🧪 TEST MODE: Local recording cleaned up");
    }

    isConnected(): boolean {
        return true; // Local recording is always "connected"
    }
}

/**
 * Factory function to create the appropriate audio handler based on mode
 */
export function createAudioHandler(
    testMode: boolean,
    websocketUrl: string,
    userInfo?: { userId: string; meetingId?: string }
): AudioHandler {
    if (testMode) {
        return new LocalRecordingAudioHandler();
    } else {
        return new WebSocketAudioHandler(websocketUrl, userInfo);
    }
}