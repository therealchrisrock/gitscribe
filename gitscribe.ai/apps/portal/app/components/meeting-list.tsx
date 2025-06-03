import React from 'react';
import { Link } from '@remix-run/react';

interface Meeting {
    id: string;
    title: string;
    type: string;
    status: string;
    start_time: string;
    end_time?: string;
    meeting_url: string;
    recording_path?: string | null;
    created_at: string;
    updated_at: string;
}

interface MeetingsData {
    meetings: Meeting[];
    total: number;
}

interface MeetingListProps {
    meetingsData: MeetingsData;
    className?: string;
    onMeetingSelect?: (meeting: Meeting) => void;
    currentMeeting?: Meeting | null;
}

export function MeetingList({ meetingsData, className = '', onMeetingSelect, currentMeeting }: MeetingListProps) {
    const { meetings } = meetingsData;

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const formatDuration = (startTime: string, endTime?: string) => {
        if (!endTime) return 'Ongoing';

        const start = new Date(startTime);
        const end = new Date(endTime);
        const durationMs = end.getTime() - start.getTime();
        const minutes = Math.floor(durationMs / (1000 * 60));

        if (minutes < 60) {
            return `${minutes}m`;
        } else {
            const hours = Math.floor(minutes / 60);
            const remainingMinutes = minutes % 60;
            return `${hours}h ${remainingMinutes}m`;
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'completed':
                return 'bg-green-100 text-green-800';
            case 'in_progress':
                return 'bg-blue-100 text-blue-800';
            case 'scheduled':
                return 'bg-yellow-100 text-yellow-800';
            case 'failed':
                return 'bg-red-100 text-red-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getTypeIcon = (type: string) => {
        switch (type) {
            case 'zoom':
                return '📹';
            case 'google_meet':
                return '🎥';
            case 'microsoft_teams':
                return '💼';
            default:
                return '🎙️';
        }
    };

    const handleMeetingClick = (e: React.MouseEvent, meeting: Meeting) => {
        if (onMeetingSelect) {
            e.preventDefault();
            onMeetingSelect(meeting);
        }
    };

    const isCurrentMeeting = (meeting: Meeting) => {
        return currentMeeting?.id === meeting.id;
    };

    return (
        <div className={`bg-card p-6 rounded-lg border ${className}`}>
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold">Past Meetings</h2>
                <div className="flex items-center gap-4">
                    <span className="text-sm text-gray-500">
                        {meetings.length} meeting{meetings.length !== 1 ? 's' : ''}
                    </span>
                    {onMeetingSelect && (
                        <span className="text-xs text-blue-600">
                            Click to use for recording
                        </span>
                    )}
                </div>
            </div>

            {meetings.length === 0 ? (
                <div className="text-center py-8">
                    <div className="text-gray-400 text-4xl mb-2">🎙️</div>
                    <p className="text-gray-600">No meetings found</p>
                    <p className="text-sm text-gray-500 mt-1">
                        Your meeting history will appear here
                    </p>
                </div>
            ) : (
                <div className="space-y-3">
                    {meetings.map((meeting) => (
                        <div key={meeting.id} className="relative">
                            <Link
                                to={`/meetings/${meeting.id}`}
                                onClick={(e) => handleMeetingClick(e, meeting)}
                                className={`
                                    block p-4 rounded-lg border transition-colors
                                    ${isCurrentMeeting(meeting)
                                        ? 'bg-blue-50 border-blue-200 ring-2 ring-blue-500 ring-opacity-50'
                                        : 'bg-gray-50 hover:bg-gray-100 border-gray-200'
                                    }
                                `}
                            >
                                <div className="flex items-start justify-between">
                                    <div className="flex-1">
                                        <div className="flex items-center gap-2 mb-1">
                                            <span className="text-lg">{getTypeIcon(meeting.type)}</span>
                                            <h3 className="font-medium text-gray-900 truncate">
                                                {meeting.title}
                                            </h3>
                                            <span className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(meeting.status)}`}>
                                                {meeting.status.replace('_', ' ')}
                                            </span>
                                            {isCurrentMeeting(meeting) && (
                                                <span className="px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800">
                                                    Active
                                                </span>
                                            )}
                                        </div>

                                        <div className="flex items-center gap-4 text-sm text-gray-600">
                                            <span>📅 {formatDate(meeting.start_time)}</span>
                                            <span>⏱️ {formatDuration(meeting.start_time, meeting.end_time)}</span>
                                            {meeting.recording_path && (
                                                <span className="text-green-600">🎬 Recording</span>
                                            )}
                                        </div>
                                    </div>

                                    <div className="ml-4 text-gray-400">
                                        {onMeetingSelect ? (
                                            <div className="text-blue-500">
                                                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                                                    <path d="M8 12a4 4 0 100-8 4 4 0 000 8zm0 1A5 5 0 108 3a5 5 0 000 10z" />
                                                    <path d="M8 8a.5.5 0 01-.5-.5V5.707L6.354 6.854a.5.5 0 11-.708-.708l2-2a.5.5 0 01.708 0l2 2a.5.5 0 01-.708.708L8.5 5.707V7.5A.5.5 0 018 8z" />
                                                </svg>
                                            </div>
                                        ) : (
                                            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                                                <path fillRule="evenodd" d="M6.22 3.22a.75.75 0 011.06 0l4.25 4.25a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06-1.06L9.94 8 6.22 4.28a.75.75 0 010-1.06z" />
                                            </svg>
                                        )}
                                    </div>
                                </div>
                            </Link>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
} 