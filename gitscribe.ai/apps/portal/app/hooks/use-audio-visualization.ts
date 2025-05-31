import { useMemo } from "react";
import type { AudioRecordingState } from './use-audio-recording';

export interface AudioVisualizationState {
    microphoneStyle: React.CSSProperties;
    audioIndicatorStyle: React.CSSProperties;
    statusMessage: string;
    connectionStatus: 'connected' | 'disconnected' | 'connecting';
}

export function useAudioVisualization(
    state: AudioRecordingState,
    config: {
        minIntensity: number;
        scaleMultiplier: number;
        opacityBase: number;
        opacityMultiplier: number;
    }
): AudioVisualizationState {
    const microphoneStyle = useMemo(() => {
        if (!state.isRecording) return {};

        const intensity = Math.max(state.audioLevel, config.minIntensity);
        const scale = 1 + (intensity * config.scaleMultiplier);

        // Add a subtle glow effect based on audio level
        const glowIntensity = intensity * 20;

        return {
            transform: `scale(${scale})`,
            transition: 'none', // Remove transition for instant response
            boxShadow: `0 0 ${glowIntensity}px rgba(239, 68, 68, ${intensity * 0.6})` // Red glow
        };
    }, [state.isRecording, state.audioLevel, config.minIntensity, config.scaleMultiplier]);

    const audioIndicatorStyle = useMemo(() => {
        if (!state.isRecording) return { opacity: 0, transform: 'scaleY(0.1)' };

        const intensity = Math.max(state.audioLevel, config.minIntensity);
        const opacity = config.opacityBase + (intensity * config.opacityMultiplier);

        // Use a more dramatic scaling for the height
        const heightScale = 0.1 + (intensity * 0.9); // Scale from 10% to 100%

        return {
            opacity: Math.min(opacity, 1),
            transform: `scaleY(${heightScale})`,
            transition: 'none' // Remove transition for instant response
        };
    }, [state.isRecording, state.audioLevel, config.minIntensity, config.opacityBase, config.opacityMultiplier]);

    const statusMessage = useMemo(() => {
        if (state.error) return state.error;
        if (state.isInitializing) return 'Initializing recording...';
        if (state.hasPermission === false) return 'Microphone permission required';
        if (state.hasPermission === null) return 'Click to request microphone access';
        if (state.isRecording) return 'Recording in progress...';
        return 'Ready to record';
    }, [state.error, state.isInitializing, state.hasPermission, state.isRecording]);

    const connectionStatus = useMemo((): 'connected' | 'disconnected' | 'connecting' => {
        if (state.isInitializing) return 'connecting';
        if (state.isConnected) return 'connected';
        return 'disconnected';
    }, [state.isInitializing, state.isConnected]);

    return {
        microphoneStyle,
        audioIndicatorStyle,
        statusMessage,
        connectionStatus
    };
} 