import { createContext, useContext } from 'react';
import type { AudioRecordingState, AudioRecordingActions } from '../hooks/use-audio-recording';

interface AudioLevelContextValue {
    state: AudioRecordingState;
    actions: AudioRecordingActions;
}

const AudioLevelContext = createContext<AudioLevelContextValue | undefined>(undefined);

export const AudioLevelProvider = AudioLevelContext.Provider;

export function useAudioLevel() {
    const ctx = useContext(AudioLevelContext);
    if (!ctx) throw new Error('useAudioLevel must be used within AudioLevelProvider');
    return ctx;
} 