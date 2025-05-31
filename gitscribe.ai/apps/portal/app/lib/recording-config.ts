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

export const DEFAULT_RECORDING_CONFIG: RecordingConfig = {
    testMode: false, // Default to production mode
    websocketUrl: "wss://api.gitscribe.ai/ws/audio",
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
    websocketUrl: import.meta.env.VITE_WEBSOCKET_URL || "wss://api.gitscribe.ai/ws/audio",
};

export const DEVELOPMENT_CONFIG: Partial<RecordingConfig> = {
    testMode: false, // Turn off test mode to use websocket handler
    websocketUrl: "ws://localhost:8080/ws/audio",
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