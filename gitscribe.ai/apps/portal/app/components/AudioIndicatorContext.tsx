import React, { createContext, useContext, useEffect, useState } from 'react';

// Types for audio tracks (simplified versions of what might come from a video SDK)
export interface AudioTrack {
    on(event: string, callback: () => void): void;
    off(event: string, callback: () => void): void;
}

export interface LocalAudioTrack extends AudioTrack {
    // Local audio track specific properties
}

export interface RemoteAudioTrack extends AudioTrack {
    // Remote audio track specific properties
}

// Utility functions
export const isIOS = typeof navigator !== 'undefined' && /iPad|iPhone|iPod/.test(navigator.userAgent);

let clipIdCounter = 0;
export const getUniqueClipId = (): string => {
    return `${Date.now()}-${++clipIdCounter}`;
};

// Timer utility similar to what might be used in video SDKs
export const interval = (callback: () => void, delay: number) => {
    const intervalId = setInterval(callback, delay);
    return {
        stop: () => clearInterval(intervalId)
    };
};

// Initialize analyser utility
export const initializeAnalyser = (mediaStream: MediaStream): AnalyserNode | undefined => {
    try {
        const audioContext = new AudioContext();
        const analyser = audioContext.createAnalyser();

        analyser.fftSize = 256;
        analyser.smoothingTimeConstant = 0.8;
        analyser.minDecibels = -90;
        analyser.maxDecibels = -10;

        const source = audioContext.createMediaStreamSource(mediaStream);
        source.connect(analyser);

        // Store the audio context for cleanup
        (analyser as AnalyserNode & { _audioContext: AudioContext })._audioContext = audioContext;

        return analyser;
    } catch (error) {
        console.error('Failed to initialize audio analyser:', error);
        return undefined;
    }
};

// Hook to check if track is enabled
export const useIsTrackEnabled = (audioTrack?: LocalAudioTrack | RemoteAudioTrack): boolean => {
    const [isEnabled, setIsEnabled] = useState(false);

    useEffect(() => {
        if (audioTrack) {
            // For now, assume track is enabled if it exists
            // In a real implementation, this would check the track's enabled state
            setIsEnabled(true);
        } else {
            setIsEnabled(false);
        }
    }, [audioTrack]);

    return isEnabled;
};

// Hook to get media stream track
export const useMediaStreamTrack = (audioTrack?: AudioTrack): MediaStreamTrack | undefined => {
    const [mediaStreamTrack, setMediaStreamTrack] = useState<MediaStreamTrack | undefined>(undefined);

    useEffect(() => {
        if (audioTrack) {
            // In a real implementation, this would extract the MediaStreamTrack from the AudioTrack
            // For now, we'll return undefined since we don't have a real audio track
            setMediaStreamTrack(undefined);
        }
    }, [audioTrack]);

    return mediaStreamTrack;
};

// Context for audio indicator
interface AudioIndicatorContextValue {
    audioTrack?: AudioTrack;
    setAudioTrack: (track?: AudioTrack) => void;
    isInitialized: boolean;
}

const AudioIndicatorContext = createContext<AudioIndicatorContextValue | undefined>(undefined);

export const AudioIndicatorProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [isInitialized, setIsInitialized] = useState(false);
    const [audioTrack, setAudioTrack] = useState<AudioTrack | undefined>(undefined);

    useEffect(() => {
        setIsInitialized(true);
    }, []);

    return (
        <AudioIndicatorContext.Provider value={{
            audioTrack,
            setAudioTrack,
            isInitialized
        }}>
            {children}
        </AudioIndicatorContext.Provider>
    );
};

export function useAudioIndicator() {
    const ctx = useContext(AudioIndicatorContext);
    if (!ctx) throw new Error('useAudioIndicator must be used within AudioIndicatorProvider');
    return ctx;
} 