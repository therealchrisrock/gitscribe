export interface RecordingConfig {
    testMode: boolean;
    websocketUrl: string;
    timeSlice: number;
    audioConstraints: MediaStreamConstraints['audio'];
    visualization: {
        minIntensity: number;
        scaleMultiplier: number;
        opacityBase: number;
        opacityMultiplier: number;
        audioFloor: number; // Minimum audio level threshold
        audioCeiling: number; // Maximum audio level for normalization
    };
}

// Default WebSocket connection options
export interface WebSocketOptions {
    provider: string;
    mode: string;
    language: string;
    speaker_diarization: boolean;
    real_time?: boolean;
    cost_optimized?: boolean;
}

export const DEFAULT_WEBSOCKET_OPTIONS: WebSocketOptions = {
    provider: "assemblyai",
    mode: "batch",
    language: "en",
    speaker_diarization: true,
    real_time: false,
    cost_optimized: false,
};

// Helper function to build WebSocket URL with options
const buildWebSocketUrl = (baseUrl: string, options: WebSocketOptions): string => {
    const params = new URLSearchParams();

    Object.entries(options).forEach(([key, value]) => {
        if (value !== undefined) {
            params.append(key, value.toString());
        }
    });

    return `${baseUrl}?${params.toString()}`;
};

export const DEFAULT_RECORDING_CONFIG: RecordingConfig = {
    testMode: false, // Default to production mode
    websocketUrl: buildWebSocketUrl("wss://api.gitscribe.ai/ws/audio", DEFAULT_WEBSOCKET_OPTIONS),
    timeSlice: 100, // Generate data every 100ms
    audioConstraints: {
        echoCancellation: true,
        noiseSuppression: true,
        sampleRate: 44100,
    },
    visualization: {
        minIntensity: 0.05,
        scaleMultiplier: 0.8,
        opacityBase: 0.3,
        opacityMultiplier: 0.7,
        audioFloor: 0.005,
        audioCeiling: 0.35
    }
};

// Environment-specific configurations
export const PRODUCTION_CONFIG: Partial<RecordingConfig> = {
    testMode: false,
    websocketUrl: buildWebSocketUrl(
        import.meta.env.VITE_WEBSOCKET_URL || "wss://api.gitscribe.ai/ws/audio",
        DEFAULT_WEBSOCKET_OPTIONS
    ),
};

export const DEVELOPMENT_CONFIG: Partial<RecordingConfig> = {
    testMode: false, // Turn off test mode to use websocket handler
    websocketUrl: buildWebSocketUrl("ws://localhost:8080/ws/enhanced-audio", DEFAULT_WEBSOCKET_OPTIONS),
};

export const getRecordingConfig = (): RecordingConfig => {
    const baseConfig = DEFAULT_RECORDING_CONFIG;

    if (typeof window !== 'undefined') {
        // Client-side environment detection
        const isDevelopment = window.location.hostname === 'localhost' ||
            window.location.hostname === '127.0.0.1' ||
            window.location.hostname.includes('dev') ||
            window.location.port === '3000'; // Common dev port

        if (isDevelopment) {
            return { ...baseConfig, ...DEVELOPMENT_CONFIG };
        } else {
            return { ...baseConfig, ...PRODUCTION_CONFIG };
        }
    }

    return baseConfig;
}; 