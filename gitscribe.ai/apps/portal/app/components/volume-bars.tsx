import { useEffect, useRef, useState } from "react";

const BAR_COUNT = 20;
const SMOOTHING_ALPHA = 0.2; // lower = slower animation

type VolumeVisualizerProps = {
    level: number;
};

function normalizeVolume(volume: number): number {
    const safe = Math.max(volume, 0.001);
    const db = 20 * Math.log10(safe);
    return Math.min(1, Math.max(0, (db + 40) / 40));
}

export function VolumeBars({ level }: VolumeVisualizerProps) {
    // RMS (Root Mean Square) normalization
    function normalizeRMS(volume: number): number {
        // Convert to RMS value (square root of squared value)
        const rms = Math.sqrt(Math.abs(volume));
        // Scale to 0-1 range with some headroom
        return Math.min(1, rms * 2);
    }
    const [smoothedLevel, setSmoothedLevel] = useState(0);

    const last = useRef(0);

    useEffect(() => {
        // Apply a power curve to boost mid/high values
        const curved = Math.pow(Math.max(0, Math.min(1, level)), 0.5); // sqrt curve
        const next = SMOOTHING_ALPHA * curved + (1 - SMOOTHING_ALPHA) * last.current;
        last.current = next;
        setSmoothedLevel(next);
    }, [level]);

    return (
        <div className="flex items-center gap-[3px] h-full  px-2 py-1">
            {Array.from({ length: BAR_COUNT }).map((_, i) => {
                const threshold = (i + 1) / BAR_COUNT;
                const isActive = smoothedLevel >= threshold;
                return (
                    <div
                        key={i}
                        className={`w-[10px] h-[10px] rounded-full border border-black  transition-colors duration-100`}
                        style={{
                            backgroundColor: isActive ? "#000" : "",
                            opacity: isActive ? 1 : 0.7,
                        }}
                    />
                );
            })}
        </div>
    );
}