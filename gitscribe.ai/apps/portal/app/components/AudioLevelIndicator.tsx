import { useAudioLevel } from './AudioLevelContext';
import { AudioIndicator, AudioIndicatorProvider } from './audio-ind';

export function AudioLevelIndicator({ layout = "default" }: { layout?: "default" | "bars" }) {
    const { state } = useAudioLevel();

    if (layout === "bars") {
        return (<></>
            // <AudioIndicatorProvider>
            //     <AudioIndicator color="white" />
            // </AudioIndicatorProvider>
        );
        // return (
        //     <div className="w-full max-w-sm flex flex-col items-center">
        //         <VolumeBars level={state.audioLevel} />
        //         {state.isRecording && (
        //             <div className="text-xs font-mono text-gray-600 mt-2">
        //                 {Math.round(state.audioLevel * 100)}%
        //             </div>
        //         )}
        //     </div>
        // );
    }

    return (
        <div className="w-full max-w-sm">
            <div className="flex items-center space-x-3">
                <span className="text-sm font-medium text-gray-700 min-w-[80px]">Audio Level</span>
                <div className="flex-1 relative">
                    {/* Background bar */}
                    <div className="w-full bg-gray-200 rounded-full h-4 overflow-hidden border border-gray-300">
                        {/* Active level bar */}
                        <div
                            className="h-full bg-gradient-to-r from-green-400 via-yellow-400 to-red-500 rounded-full transition-all duration-0 ease-linear shadow-sm"
                            style={{
                                width: `${state.audioLevel * 100}%`,
                                minWidth: state.isRecording ? '2px' : '0px'
                            }}
                        />
                    </div>
                    {/* Level percentage text */}
                    {state.isRecording && (
                        <div className="absolute -top-6 left-0 text-xs font-mono text-gray-600">
                            {Math.round(state.audioLevel * 100)}%
                        </div>
                    )}
                </div>
            </div>

            {/* Audio level bars visualization */}
            {state.isRecording && (
                <div className="flex items-end justify-center space-x-1 mt-3 h-8">
                    {Array.from({ length: 12 }, (_, i) => {
                        const barHeight = Math.max(0.1, state.audioLevel - (i * 0.08));
                        const isActive = barHeight > 0.1;
                        return (
                            <div
                                key={i}
                                className={`w-2 rounded-t transition-none ${isActive
                                    ? i < 4 ? 'bg-green-500'
                                        : i < 8 ? 'bg-yellow-500'
                                            : 'bg-red-500'
                                    : 'bg-gray-300'
                                    }`}
                                style={{
                                    height: `${Math.max(2, barHeight * 32)}px`,
                                    opacity: isActive ? 0.8 + (barHeight * 0.2) : 0.3
                                }}
                            />
                        );
                    })}
                </div>
            )}
        </div>
    );
} 