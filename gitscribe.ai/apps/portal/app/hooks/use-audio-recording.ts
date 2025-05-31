import { useState, useRef, useCallback, useEffect } from 'react';
import type { RecordingConfig } from '../lib/recording-config';
import { createAudioHandler, type AudioHandler, type AudioPlaybackData } from '../lib/audio-handlers';

export interface AudioRecordingState {
    isRecording: boolean;
    hasPermission: boolean | null;
    error: string | null;
    audioLevel: number;
    isConnected: boolean;
    isInitializing: boolean;
    audioPlayback: AudioPlaybackData | null;
}

export interface AudioRecordingActions {
    startRecording: () => Promise<void>;
    stopRecording: () => Promise<void>;
    requestPermission: () => Promise<void>;
    clearAudioPlayback: () => void;
}

export interface UseAudioRecordingReturn {
    state: AudioRecordingState;
    actions: AudioRecordingActions;
}

export function useAudioRecording(config: RecordingConfig): UseAudioRecordingReturn {
    // State
    const [isRecording, setIsRecording] = useState(false);
    const [hasPermission, setHasPermission] = useState<boolean | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [audioLevel, setAudioLevel] = useState(0);
    const [isConnected, setIsConnected] = useState(false);
    const [isInitializing, setIsInitializing] = useState(false);
    const [audioPlayback, setAudioPlayback] = useState<AudioPlaybackData | null>(null);

    // Refs
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);
    const streamRef = useRef<MediaStream | null>(null);
    const audioContextRef = useRef<AudioContext | null>(null);
    const analyserRef = useRef<AnalyserNode | null>(null);
    const audioHandlerRef = useRef<AudioHandler | null>(null);
    const animationFrameRef = useRef<number | null>(null);
    const isRecordingRef = useRef<boolean>(false); // Track recording state in ref

    // Audio level analysis
    const analyzeAudioLevel = useCallback(() => {
        if (!analyserRef.current || !isRecordingRef.current) {
            return;
        }

        const bufferLength = analyserRef.current.frequencyBinCount;
        const dataArray = new Uint8Array(bufferLength);
        analyserRef.current.getByteFrequencyData(dataArray);

        // Calculate RMS (Root Mean Square) for audio level
        const inputBuffer = Array.from(dataArray).map(x => x / 255);
        const rms = Math.sqrt(inputBuffer.reduce((sum, x) => sum + x * x, 0) / inputBuffer.length);

        // Debug logging every 30 frames (~0.5 seconds at 60fps)
        if (Math.random() < 0.033) {
            console.log('🎵 Audio levels:', {
                raw: rms.toFixed(3),
                maxFreq: Math.max(...dataArray),
                avgFreq: (dataArray.reduce((a, b) => a + b, 0) / bufferLength).toFixed(1),
                isRecording: isRecordingRef.current
            });
        }

        setAudioLevel(rms);

        // Continue the animation loop while recording
        if (isRecordingRef.current) {
            animationFrameRef.current = requestAnimationFrame(analyzeAudioLevel);
        }
    }, []);

    // Setup audio analysis
    const setupAudioAnalysis = useCallback((stream: MediaStream) => {
        try {
            audioContextRef.current = new AudioContext();
            analyserRef.current = audioContextRef.current.createAnalyser();

            // Configure analyser for better responsiveness
            analyserRef.current.fftSize = 512; // Increased for better frequency resolution
            analyserRef.current.smoothingTimeConstant = 0.1; // Much lower for instant response (was 0.3)
            analyserRef.current.minDecibels = -90;
            analyserRef.current.maxDecibels = -10;

            const source = audioContextRef.current.createMediaStreamSource(stream);
            source.connect(analyserRef.current);

            // Don't start analyzing here - wait for recording to actually start
            console.log('🎵 Audio analysis setup complete');
        } catch (err) {
            console.error('Error setting up audio analysis:', err);
        }
    }, []);

    // Get supported MIME type
    const getSupportedMimeType = useCallback((): string => {
        const types = [
            'audio/webm;codecs=opus',
            'audio/webm',
            'audio/mp4',
            'audio/ogg;codecs=opus',
            'audio/wav'
        ];

        for (const type of types) {
            if (MediaRecorder.isTypeSupported(type)) {
                return type;
            }
        }

        return ''; // Fallback to default
    }, []);

    // Request microphone permission
    const requestPermission = useCallback(async () => {
        try {
            setError(null);
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: config.audioConstraints
            });

            setHasPermission(true);

            // Clean up the test stream
            stream.getTracks().forEach(track => track.stop());
        } catch (err) {
            console.error('Permission denied:', err);
            setHasPermission(false);
            setError('Microphone permission denied. Please allow access to record audio.');
        }
    }, [config.audioConstraints]);

    // Start recording
    const startRecording = useCallback(async () => {
        try {
            setError(null);
            setIsInitializing(true);
            setAudioPlayback(null); // Clear any previous playback

            // Check permission first
            if (hasPermission === null) {
                await requestPermission();
            }

            if (hasPermission === false) {
                throw new Error('Microphone permission required');
            }

            // Get media stream
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: config.audioConstraints
            });
            streamRef.current = stream;

            // Setup audio analysis for visual feedback
            setupAudioAnalysis(stream);

            // Initialize audio handler
            audioHandlerRef.current = createAudioHandler(config.testMode, config.websocketUrl);
            await audioHandlerRef.current.initialize();
            setIsConnected(audioHandlerRef.current.isConnected());

            // In test mode, set up the audio ready callback
            if (config.testMode && audioHandlerRef.current.setOnAudioReady) {
                audioHandlerRef.current.setOnAudioReady((audioData) => {
                    console.log('🎬 Audio ready, setting playback data');
                    setAudioPlayback(audioData);
                });
            }

            // Setup MediaRecorder
            const mimeType = getSupportedMimeType();
            const options = mimeType ? { mimeType } : undefined;

            mediaRecorderRef.current = new MediaRecorder(stream, options);

            // Handle audio data
            mediaRecorderRef.current.ondataavailable = async (event) => {
                if (event.data.size > 0 && audioHandlerRef.current) {
                    const timestamp = new Date().toISOString();
                    try {
                        await audioHandlerRef.current.handleAudioChunk(event.data, timestamp);
                    } catch (err) {
                        console.error('Error handling audio chunk:', err);
                        setError('Failed to process audio data');
                    }
                }
            };

            mediaRecorderRef.current.onstart = () => {
                console.log(`🎙️ Recording started (${config.testMode ? 'TEST MODE' : 'PRODUCTION MODE'})`);
                setIsRecording(true);
                isRecordingRef.current = true; // Update ref immediately
                setIsInitializing(false);

                // Start audio level analysis
                if (analyserRef.current) {
                    console.log('🎵 Starting audio level analysis');
                    analyzeAudioLevel();
                }
            };

            mediaRecorderRef.current.onstop = async () => {
                console.log('🛑 Recording stopped');
                setIsRecording(false);
                isRecordingRef.current = false; // Update ref immediately

                // Finalize the audio handler
                if (audioHandlerRef.current) {
                    try {
                        await audioHandlerRef.current.finalize();
                        console.log('✅ Audio handler finalized successfully');
                    } catch (err) {
                        console.error('❌ Error finalizing audio handler:', err);
                    }

                    // Clean up audio handler after finalization
                    audioHandlerRef.current.cleanup();
                    audioHandlerRef.current = null;
                }

                setAudioLevel(0);
                setIsConnected(false);
            };

            mediaRecorderRef.current.onerror = (event) => {
                console.error('MediaRecorder error:', event);
                setError('Recording failed. Please try again.');
                setIsRecording(false);
                isRecordingRef.current = false; // Update ref immediately
                setIsInitializing(false);
            };

            // Start recording with time slice for regular data events
            mediaRecorderRef.current.start(config.timeSlice);

        } catch (err) {
            console.error('Error starting recording:', err);
            setError(err instanceof Error ? err.message : 'Failed to start recording');
            setIsRecording(false);
            setIsInitializing(false);
        }
    }, [
        hasPermission,
        requestPermission,
        config.audioConstraints,
        config.testMode,
        config.websocketUrl,
        config.timeSlice,
        setupAudioAnalysis,
        getSupportedMimeType
    ]);

    // Stop recording
    const stopRecording = useCallback(async () => {
        try {
            // Stop recording state immediately
            isRecordingRef.current = false;

            // Stop MediaRecorder
            if (mediaRecorderRef.current && mediaRecorderRef.current.state === 'recording') {
                mediaRecorderRef.current.stop();
            }

            // Stop audio analysis
            if (animationFrameRef.current) {
                cancelAnimationFrame(animationFrameRef.current);
                animationFrameRef.current = null;
            }

            // Clean up audio context
            if (audioContextRef.current && audioContextRef.current.state !== 'closed') {
                await audioContextRef.current.close();
                audioContextRef.current = null;
            }

            // Stop media stream
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => track.stop());
                streamRef.current = null;
            }

            // NOTE: audioHandler cleanup moved to onstop callback after finalize()

        } catch (err) {
            console.error('Error stopping recording:', err);
            setError('Error stopping recording');
        }
    }, [config.testMode]);

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            if (animationFrameRef.current) {
                cancelAnimationFrame(animationFrameRef.current);
            }
            if (audioContextRef.current && audioContextRef.current.state !== 'closed') {
                audioContextRef.current.close();
            }
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => track.stop());
            }
            if (audioHandlerRef.current) {
                audioHandlerRef.current.cleanup();
            }
        };
    }, []);

    return {
        state: {
            isRecording,
            hasPermission,
            error,
            audioLevel,
            isConnected,
            isInitializing,
            audioPlayback
        },
        actions: {
            startRecording,
            stopRecording,
            requestPermission,
            clearAudioPlayback: () => {
                setAudioPlayback(null);
            }
        }
    };
} 