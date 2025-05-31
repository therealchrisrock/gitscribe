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

/**
 * Production audio handler that streams audio chunks via WebSocket
 */
export class WebSocketAudioHandler implements AudioHandler {
    private websocket: WebSocket | null = null;
    private totalBytesSent = 0;
    private chunkCount = 0;

    constructor(private websocketUrl: string) { }

    async initialize(): Promise<void> {
        return new Promise((resolve, reject) => {
            try {
                this.websocket = new WebSocket(this.websocketUrl);

                this.websocket.onopen = () => {
                    console.log("WebSocket connected");
                    resolve();
                };

                this.websocket.onclose = () => {
                    console.log("WebSocket disconnected");
                    this.websocket = null;
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

    async handleAudioChunk(data: Blob, timestamp: string): Promise<void> {
        if (this.websocket?.readyState === WebSocket.OPEN) {
            try {
                const buffer = await data.arrayBuffer();
                this.chunkCount++;
                this.totalBytesSent += buffer.byteLength;

                console.log(`[${timestamp}] Sending audio chunk #${this.chunkCount}:`, {
                    chunkSize: buffer.byteLength,
                    totalSent: this.totalBytesSent
                });

                this.websocket.send(buffer);
            } catch (error) {
                console.error("Error sending audio data:", error);
                throw error;
            }
        } else {
            console.warn("WebSocket not open, skipping chunk");
        }
    }

    async finalize(): Promise<void> {
        console.log(`📊 WebSocket streaming complete: ${this.chunkCount} chunks, ${this.totalBytesSent} total bytes`);
    }

    cleanup(): void {
        if (this.websocket) {
            this.websocket.close();
            this.websocket = null;
        }
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
export function createAudioHandler(testMode: boolean, websocketUrl: string): AudioHandler {
    if (testMode) {
        return new LocalRecordingAudioHandler();
    } else {
        return new WebSocketAudioHandler(websocketUrl);
    }
}