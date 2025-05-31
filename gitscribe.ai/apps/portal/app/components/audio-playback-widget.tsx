import React, { useRef, useEffect } from 'react';
import type { AudioPlaybackData } from '../lib/audio-handlers';

interface AudioPlaybackWidgetProps {
    audioData: AudioPlaybackData;
    onClose: () => void;
}

export function AudioPlaybackWidget({ audioData, onClose }: AudioPlaybackWidgetProps) {
    const audioRef = useRef<HTMLAudioElement>(null);

    useEffect(() => {
        // Auto-play the audio when component mounts
        if (audioRef.current) {
            audioRef.current.play().catch(console.error);
        }
    }, []);

    const handleDownload = () => {
        const a = document.createElement('a');
        a.href = audioData.audioUrl;
        a.download = `test-recording-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.webm`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
    };

    const formatFileSize = (bytes: number) => {
        if (bytes < 1024) return `${bytes} B`;
        if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
        return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    };

    return (
        <div className="bg-white rounded-lg shadow-lg border-2 border-green-200 p-6 mt-4">
            <div className="flex justify-between items-start mb-4">
                <div>
                    <h3 className="text-lg font-semibold text-green-700 flex items-center">
                        🎵 Test Recording Ready
                        <span className="ml-2 text-sm bg-green-100 text-green-600 px-2 py-1 rounded-full">
                            Test Mode
                        </span>
                    </h3>
                    <p className="text-sm text-gray-600 mt-1">✅ Recording completed successfully!</p>
                </div>
                <button
                    onClick={onClose}
                    className="text-gray-400 hover:text-gray-600 text-xl font-bold"
                    title="Close"
                >
                    ×
                </button>
            </div>

            {/* Recording Statistics */}
            <div className="bg-gray-50 rounded-lg p-4 mb-4">
                <h4 className="text-sm font-semibold text-gray-700 mb-3">Recording Statistics</h4>
                <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="text-center">
                        <div className="text-lg font-bold text-blue-600">{audioData.chunkCount}</div>
                        <div className="text-gray-600">Audio Chunks</div>
                    </div>
                    <div className="text-center">
                        <div className="text-lg font-bold text-blue-600">{formatFileSize(audioData.totalSize)}</div>
                        <div className="text-gray-600">Total Size</div>
                    </div>
                    <div className="text-center">
                        <div className="text-lg font-bold text-blue-600">{audioData.mimeType.split('/')[1]?.split(';')[0] || 'Unknown'}</div>
                        <div className="text-gray-600">Audio Format</div>
                    </div>
                    <div className="text-center">
                        <div className="text-lg font-bold text-blue-600">{formatFileSize(audioData.totalSize / audioData.chunkCount)}</div>
                        <div className="text-gray-600">Avg Chunk Size</div>
                    </div>
                </div>
            </div>

            {/* Audio Player */}
            <div className="mb-4">
                <h4 className="text-sm font-semibold text-gray-700 mb-2">Audio Playback</h4>
                <audio
                    ref={audioRef}
                    controls
                    className="w-full"
                    preload="auto"
                >
                    <source src={audioData.audioUrl} type={audioData.mimeType} />
                    Your browser does not support the audio element.
                </audio>
            </div>

            {/* Action Buttons */}
            <div className="flex justify-center space-x-3">
                <button
                    onClick={handleDownload}
                    className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center"
                >
                    📥 Download Recording
                </button>
                <button
                    onClick={onClose}
                    className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg text-sm font-medium"
                >
                    ✅ Close
                </button>
            </div>
        </div>
    );
} 