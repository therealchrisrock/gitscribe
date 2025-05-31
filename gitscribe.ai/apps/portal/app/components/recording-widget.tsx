"use client";

import React from 'react';
import { useAudioRecording } from '../hooks/use-audio-recording';
import { useAudioVisualization } from '../hooks/use-audio-visualization';
import { getRecordingConfig, type RecordingConfig } from '../lib/recording-config';
import { AudioLevelProvider } from './AudioLevelContext';
import { AudioLevelIndicator } from './AudioLevelIndicator';
import { AudioPlaybackWidget } from './audio-playback-widget';

// SVG Icons
const MicIcon = () => (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M12 1C10.34 1 9 2.34 9 4V12C9 13.66 10.34 15 12 15C13.66 15 15 13.66 15 12V4C15 2.34 13.66 1 12 1Z" fill="currentColor" />
        <path d="M19 10V12C19 16.42 15.42 20 11 20H13C17.42 20 21 16.42 21 12V10H19Z" fill="currentColor" />
        <path d="M5 10V12C5 16.42 8.58 20 13 20H11C6.58 20 3 16.42 3 12V10H5Z" fill="currentColor" />
        <path d="M12 22V20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
);

const MicOffIcon = () => (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M16.5 12C16.5 12.66 16.35 13.3 16.08 13.88L17.23 15.03C17.73 14.22 18 13.16 18 12H16.5Z" fill="currentColor" />
        <path d="M12 15.5C11.34 15.5 10.7 15.35 10.12 15.08L8.97 16.23C9.78 16.73 10.84 17 12 17C13.66 17 15 15.66 15 14V13.41L13.59 12H12V15.5Z" fill="currentColor" />
        <path d="M19.07 4.93L4.93 19.07L6.34 20.48L20.48 6.34L19.07 4.93Z" fill="currentColor" />
        <path d="M9 9V4C9 2.34 10.34 1 12 1C13.66 1 15 2.34 15 4V9L9 9Z" fill="currentColor" />
    </svg>
);

const StopIcon = () => (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="6" y="6" width="12" height="12" rx="2" fill="currentColor" />
    </svg>
);

interface RecordingWidgetProps {
    config?: Partial<RecordingConfig>;
    className?: string;
}

export function RecordingWidget({ config = {}, className = '' }: RecordingWidgetProps) {
    // Get environment-aware configuration
    const defaultConfig = getRecordingConfig();
    const finalConfig = { ...defaultConfig, ...config };

    const { state, actions } = useAudioRecording(finalConfig);

    const visualization = useAudioVisualization(state, finalConfig.visualization);

    const handleToggleRecording = async () => {
        if (state.isRecording) {
            await actions.stopRecording();
        } else {
            if (state.hasPermission === null) {
                await actions.requestPermission();
            }
            if (state.hasPermission !== false) {
                await actions.startRecording();
            }
        }
    };

    const getButtonIcon = () => {
        if (state.isRecording) return <StopIcon />;
        if (state.hasPermission === false) return <MicOffIcon />;
        return <MicIcon />;
    };

    const getButtonColor = () => {
        if (state.error) return 'bg-red-500 hover:bg-red-600';
        if (state.isRecording) return 'bg-red-500 hover:bg-red-600';
        if (state.hasPermission === false) return 'bg-gray-400 hover:bg-gray-500';
        return 'bg-blue-500 hover:bg-blue-600';
    };

    const getConnectionIndicator = () => {
        const colors = {
            connected: 'bg-green-500',
            connecting: 'bg-yellow-500',
            disconnected: 'bg-gray-400'
        };

        return (
            <div className={`w-2 h-2 rounded-full ${colors[visualization.connectionStatus]}`} />
        );
    };

    return (
        <AudioLevelProvider value={{ state, actions }}>
            <div className={`bg-white rounded-lg shadow-lg p-6 ${className}`}>
                <div className="flex flex-col items-center space-y-4">
                    {/* Header */}
                    <div className="text-center">
                        <h3 className="text-lg font-semibold text-gray-900">Audio Recording</h3>
                        <p className="text-sm text-gray-600 mt-1">{visualization.statusMessage}</p>
                    </div>

                    {/* Recording Button with Visual Feedback */}
                    <div className="relative">
                        {/* Pulsing ring when recording */}
                        {state.isRecording && (
                            <div className="absolute inset-0 rounded-full bg-red-400 animate-ping opacity-75" />
                        )}

                        {/* Main button */}
                        <button
                            onClick={handleToggleRecording}
                            disabled={state.isInitializing}
                            className={`
                                relative w-16 h-16 rounded-full text-white flex items-center justify-center cursor-pointer
                                transition-all duration-200 ease-in-out
                                focus:outline-none focus:ring-4 focus:ring-blue-300
                                disabled:opacity-50 disabled:cursor-not-allowed
                                ${getButtonColor()}
                            `}
                            style={visualization.microphoneStyle}
                        >
                            {state.isInitializing ? (
                                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-white mx-auto" />
                            ) : (
                                getButtonIcon()
                            )}
                        </button>
                    </div>

                    {/* Audio Level Indicator */}
                    <AudioLevelIndicator layout="bars" />

                    {/* Debug Info (only show during recording) */}
                    {state.isRecording && (
                        <div className="w-full max-w-sm p-3 bg-gray-50 rounded-lg border text-xs font-mono">
                            <div className="grid grid-cols-2 gap-2 text-gray-600">
                                <div>Raw Level: <span className="font-bold text-blue-600">{state.audioLevel.toFixed(3)}</span></div>
                                <div>Percentage: <span className="font-bold text-green-600">{Math.round(state.audioLevel * 100)}%</span></div>
                                <div>Status: <span className="font-bold text-purple-600">{state.isRecording ? 'Recording' : 'Stopped'}</span></div>
                                <div>Connected: <span className="font-bold text-orange-600">{state.isConnected ? 'Yes' : 'No'}</span></div>
                                <div>Floor: <span className="font-bold text-gray-600">{(finalConfig.visualization.audioFloor * 100).toFixed(1)}%</span></div>
                                <div>Ceiling: <span className="font-bold text-gray-600">{(finalConfig.visualization.audioCeiling * 100).toFixed(1)}%</span></div>
                            </div>
                        </div>
                    )}

                    {/* Connection Status */}
                    <div className="flex items-center space-x-2 text-sm text-gray-600">
                        {getConnectionIndicator()}
                        <span>
                            {visualization.connectionStatus === 'connected' && 'Connected'}
                            {visualization.connectionStatus === 'connecting' && 'Connecting...'}
                            {visualization.connectionStatus === 'disconnected' && 'Disconnected'}
                        </span>
                    </div>

                    {/* Error Display */}
                    {state.error && (
                        <div className="w-full p-3 bg-red-50 border border-red-200 rounded-md">
                            <p className="text-sm text-red-700">{state.error}</p>
                        </div>
                    )}

                    {/* Permission Request */}
                    {state.hasPermission === false && (
                        <div className="w-full p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                            <p className="text-sm text-yellow-700 mb-2">
                                Microphone access is required to record audio.
                            </p>
                            <button
                                onClick={actions.requestPermission}
                                className="text-sm bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
                            >
                                Grant Permission
                            </button>
                        </div>
                    )}
                </div>
            </div>

            {/* Audio Playback Widget - shows when test recording is complete */}
            {state.audioPlayback && (
                <AudioPlaybackWidget
                    audioData={state.audioPlayback}
                    onClose={actions.clearAudioPlayback}
                />
            )}
        </AudioLevelProvider>
    );
} 