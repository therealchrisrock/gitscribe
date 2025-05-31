import React, { useState } from 'react';
import Draggable from 'react-draggable';
import './widget.css';

export default function App() {
    console.log('Widget Gitscribe')
    const [isExpanded, setIsExpanded] = useState(false);
    const [isRecording, setIsRecording] = useState(false);

    const toggleExpanded = () => setIsExpanded(!isExpanded);
    const toggleRecording = () => setIsRecording(!isRecording);

    return (
        <Draggable
            handle=".drag-handle"
            bounds="parent"
            defaultPosition={{ x: typeof window !== 'undefined' ? window.innerWidth - 80 : 50, y: 50 }}
        >
            <div className="widget-container">
                {/* Main Widget */}
                <div className={`widget ${isExpanded ? 'expanded' : ''}`}>
                    {/* Drag Handle */}
                    <div className="drag-handle">
                        <div className="dots">
                            <span className="dot"></span>
                            <span className="dot"></span>
                            <span className="dot"></span>
                            <span className="dot"></span>
                            <span className="dot"></span>
                            <span className="dot"></span>
                        </div>
                    </div>

                    {/* Logo/Brand Section */}
                    <div className="logo-section">
                        <div className="logo">
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                                <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor" />
                                <path d="M2 17L12 22L22 17" stroke="currentColor" strokeWidth="2" />
                                <path d="M2 12L12 17L22 12" stroke="currentColor" strokeWidth="2" />
                            </svg>
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div className="actions">
                        {/* Recording Toggle */}
                        <button
                            className={`action-btn ${isRecording ? 'recording' : ''}`}
                            onClick={toggleRecording}
                            title={isRecording ? 'Stop Recording' : 'Start Recording'}
                        >
                            <div className={`record-indicator ${isRecording ? 'recording' : ''}`}></div>
                        </button>

                        {/* Share Button */}
                        <button className="action-btn" title="Share">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                                <path d="M18 16.08C17.24 16.08 16.56 16.38 16.04 16.85L8.91 12.7C8.96 12.47 9 12.24 9 12S8.96 11.53 8.91 11.3L15.96 7.19C16.5 7.69 17.21 8 18 8C19.66 8 21 6.66 21 5S19.66 2 18 2 15 3.34 15 5C15 5.24 15.04 5.47 15.09 5.7L8.04 9.81C7.5 9.31 6.79 9 6 9C4.34 9 3 10.34 3 12S4.34 15 6 15C6.79 15 7.5 14.69 8.04 14.19L15.16 18.34C15.11 18.55 15.08 18.77 15.08 19C15.08 20.61 16.39 21.92 18 21.92S20.92 20.61 20.92 19C20.92 17.39 19.61 16.08 18 16.08Z" fill="currentColor" />
                            </svg>
                        </button>

                        {/* Target/Focus Button */}
                        <button className="action-btn" title="Focus">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                                <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" />
                                <circle cx="12" cy="12" r="3" fill="currentColor" />
                            </svg>
                        </button>

                        {/* Chat/Messages Button */}
                        <button className="action-btn" title="Messages">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                                <path d="M21 15C21 15.5304 20.7893 16.0391 20.4142 16.4142C20.0391 16.7893 19.5304 17 19 17H7L3 21V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H19C19.5304 3 20.0391 3.21071 20.4142 3.58579C20.7893 3.96086 21 4.46957 21 5V15Z" stroke="currentColor" strokeWidth="2" />
                            </svg>
                        </button>

                        {/* Database/Storage Button */}
                        <button className="action-btn" title="Storage">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                                <ellipse cx="12" cy="5" rx="9" ry="3" stroke="currentColor" strokeWidth="2" />
                                <path d="M21 12C21 13.66 16.97 15 12 15S3 13.66 3 12" stroke="currentColor" strokeWidth="2" />
                                <path d="M3 5V19C3 20.66 7.03 22 12 22S21 20.66 21 19V5" stroke="currentColor" strokeWidth="2" />
                            </svg>
                        </button>

                        {/* Settings Button */}
                        <button className="action-btn" title="Settings" onClick={toggleExpanded}>
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                                <path d="M12.22 2H11.78L11.11 4.15C10.68 4.32 10.26 4.53 9.88 4.78L7.85 3.85L7.41 4.29L8.34 6.32C8.09 6.7 7.88 7.12 7.71 7.55L5.56 8.22V8.66L7.71 9.33C7.88 9.76 8.09 10.18 8.34 10.56L7.41 12.59L7.85 13.03L9.88 12.1C10.26 12.35 10.68 12.56 11.11 12.73L11.78 14.88H12.22L12.89 12.73C13.32 12.56 13.74 12.35 14.12 12.1L16.15 13.03L16.59 12.59L15.66 10.56C15.91 10.18 16.12 9.76 16.29 9.33L18.44 8.66V8.22L16.29 7.55C16.12 7.12 15.91 6.7 15.66 6.32L16.59 4.29L16.15 3.85L14.12 4.78C13.74 4.53 13.32 4.32 12.89 4.15L12.22 2ZM12 10.5C11.17 10.5 10.5 9.83 10.5 9S11.17 7.5 12 7.5S13.5 8.17 13.5 9S12.83 10.5 12 10.5Z" fill="currentColor" />
                            </svg>
                        </button>
                    </div>

                    {/* Collapse/Expand Indicator */}
                    <div className={`collapse-indicator ${isExpanded ? 'expanded' : ''}`}>
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                            <path d="M7 14L12 9L17 14" stroke="currentColor" strokeWidth="2" />
                        </svg>
                    </div>
                </div>
            </div>
        </Draggable>
    );
}