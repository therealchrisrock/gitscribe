import React from "react";
import { RecordingWidget } from "./recording-widget";

/**
 * Examples of how to use the RecordingWidget component in different configurations
 */

// Example 1: Default configuration (uses environment-based settings)
export const DefaultRecordingWidget = () => (
    <RecordingWidget className="max-w-md mx-auto" />
);

// Example 2: Force test mode (useful for development/demos)
export const TestModeRecordingWidget = () => (
    <RecordingWidget
        className="max-w-md mx-auto"
        config={{ testMode: true }}
    />
);

// Example 3: Production mode with custom WebSocket URL
export const ProductionRecordingWidget = () => (
    <RecordingWidget
        className="max-w-md mx-auto"
        config={{
            testMode: false,
            websocketUrl: "wss://api.gitscribe.ai/ws/audio"
        }}
    />
);

// Example 4: Development mode with local server
export const DevelopmentRecordingWidget = () => (
    <RecordingWidget
        className="max-w-md mx-auto"
        config={{
            testMode: true,
            websocketUrl: "ws://localhost:8080/ws/audio"
        }}
    />
);

// Example 5: Multiple widgets with different configurations
export const MultipleRecordingWidgets = () => (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
        <div>
            <h3 className="text-lg font-semibold mb-4">Test Mode</h3>
            <TestModeRecordingWidget />
        </div>
        <div>
            <h3 className="text-lg font-semibold mb-4">Production Mode</h3>
            <ProductionRecordingWidget />
        </div>
    </div>
);

/**
 * Usage in your app:
 * 
 * ```tsx
 * import { RecordingWidget } from "./components/recording-widget";
 * 
 * // Simple usage with defaults
 * <RecordingWidget />
 * 
 * // Custom configuration
 * <RecordingWidget 
 *   config={{
 *     testMode: false,
 *     websocketUrl: "wss://your-server.com/ws/audio"
 *   }}
 *   className="my-custom-styles"
 * />
 * ```
 */ 