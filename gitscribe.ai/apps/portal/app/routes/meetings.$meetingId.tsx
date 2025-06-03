import type { LoaderFunctionArgs } from "@remix-run/node";
import { json } from "@remix-run/node";
import { useLoaderData, Link } from "@remix-run/react";
import { requireAuth } from "@/lib/auth.server";

interface Meeting {
    id: string;
    title: string;
    type: string;
    status: string;
    start_time: string;
    end_time?: string;
    meeting_url: string;
    recording_path?: string;
    created_at: string;
    updated_at: string;
}

export async function loader({ request, params }: LoaderFunctionArgs) {
    const user = await requireAuth(request);
    const { meetingId } = params;

    if (!meetingId) {
        throw new Response("Meeting ID not found", { status: 404 });
    }

    // TODO: Fetch actual meeting data from backend
    // For now, return mock data based on meetingId
    const mockMeetings: Record<string, Meeting> = {
        "meeting_1": {
            id: "meeting_1",
            title: "Team Standup",
            type: "zoom",
            status: "completed",
            start_time: "2024-01-15T09:00:00Z",
            end_time: "2024-01-15T09:30:00Z",
            meeting_url: "https://zoom.us/j/123456789",
            recording_path: "/recordings/meeting_1.mp4",
            created_at: "2024-01-15T08:55:00Z",
            updated_at: "2024-01-15T09:30:00Z"
        },
        "meeting_2": {
            id: "meeting_2",
            title: "Product Review",
            type: "google_meet",
            status: "completed",
            start_time: "2024-01-14T14:00:00Z",
            end_time: "2024-01-14T15:00:00Z",
            meeting_url: "https://meet.google.com/abc-def-ghi",
            recording_path: undefined,
            created_at: "2024-01-14T13:55:00Z",
            updated_at: "2024-01-14T15:00:00Z"
        },
        "meeting_3": {
            id: "meeting_3",
            title: "Client Demo",
            type: "microsoft_teams",
            status: "completed",
            start_time: "2024-01-12T16:00:00Z",
            end_time: "2024-01-12T17:30:00Z",
            meeting_url: "https://teams.microsoft.com/l/meetup-join/xyz",
            recording_path: "/recordings/meeting_3.mp4",
            created_at: "2024-01-12T15:55:00Z",
            updated_at: "2024-01-12T17:30:00Z"
        },
        "meeting_4": {
            id: "meeting_4",
            title: "Weekly Sync",
            type: "generic",
            status: "scheduled",
            start_time: "2024-01-16T10:00:00Z",
            end_time: undefined,
            meeting_url: "https://example.com/meeting/123",
            recording_path: undefined,
            created_at: "2024-01-16T09:00:00Z",
            updated_at: "2024-01-16T09:00:00Z"
        }
    };

    const meeting = mockMeetings[meetingId];

    if (!meeting) {
        throw new Response("Meeting not found", { status: 404 });
    }

    return json({ meeting });
}

export default function MeetingDetails() {
    const { meeting } = useLoaderData<typeof loader>();

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            timeZoneName: 'short',
        });
    };

    const formatDuration = (startTime: string, endTime?: string) => {
        if (!endTime) return 'Ongoing';

        const start = new Date(startTime);
        const end = new Date(endTime);
        const durationMs = end.getTime() - start.getTime();
        const hours = Math.floor(durationMs / (1000 * 60 * 60));
        const minutes = Math.floor((durationMs % (1000 * 60 * 60)) / (1000 * 60));

        if (hours > 0) {
            return `${hours}h ${minutes}m`;
        }
        return `${minutes}m`;
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'completed':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'in_progress':
                return 'bg-blue-100 text-blue-800 border-blue-200';
            case 'scheduled':
                return 'bg-yellow-100 text-yellow-800 border-yellow-200';
            case 'failed':
                return 'bg-red-100 text-red-800 border-red-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
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

    return (
        <div className="container mx-auto p-6">
            {/* Breadcrumb */}
            <nav className="mb-6">
                <Link to="/dashboard" className="text-blue-600 hover:text-blue-800">
                    ← Back to Dashboard
                </Link>
            </nav>

            {/* Meeting Header */}
            <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
                <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                        <span className="text-3xl">{getTypeIcon(meeting.type)}</span>
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">{meeting.title}</h1>
                            <p className="text-gray-600 mt-1">Meeting ID: {meeting.id}</p>
                        </div>
                    </div>
                    <span className={`px-3 py-1 text-sm font-medium rounded-lg border ${getStatusColor(meeting.status)}`}>
                        {meeting.status.replace('_', ' ').toUpperCase()}
                    </span>
                </div>

                {/* Meeting Details Grid */}
                <div className="grid md:grid-cols-2 gap-6">
                    <div className="space-y-4">
                        <div>
                            <h3 className="text-sm font-medium text-gray-700 mb-1">Start Time</h3>
                            <p className="text-gray-900">{formatDate(meeting.start_time)}</p>
                        </div>

                        {meeting.end_time && (
                            <div>
                                <h3 className="text-sm font-medium text-gray-700 mb-1">End Time</h3>
                                <p className="text-gray-900">{formatDate(meeting.end_time)}</p>
                            </div>
                        )}

                        <div>
                            <h3 className="text-sm font-medium text-gray-700 mb-1">Duration</h3>
                            <p className="text-gray-900">{formatDuration(meeting.start_time, meeting.end_time)}</p>
                        </div>
                    </div>

                    <div className="space-y-4">
                        <div>
                            <h3 className="text-sm font-medium text-gray-700 mb-1">Meeting Type</h3>
                            <p className="text-gray-900 capitalize">{meeting.type.replace('_', ' ')}</p>
                        </div>

                        <div>
                            <h3 className="text-sm font-medium text-gray-700 mb-1">Meeting URL</h3>
                            <a
                                href={meeting.meeting_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-blue-600 hover:text-blue-800 break-all"
                            >
                                {meeting.meeting_url}
                            </a>
                        </div>

                        {meeting.recording_path && (
                            <div>
                                <h3 className="text-sm font-medium text-gray-700 mb-1">Recording</h3>
                                <div className="flex items-center gap-2">
                                    <span className="text-green-600">🎬</span>
                                    <span className="text-gray-900">Available</span>
                                    <button className="text-sm bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700">
                                        View Recording
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* Coming Soon Sections */}
            <div className="grid md:grid-cols-2 gap-6">
                {/* Transcript Section */}
                <div className="bg-white rounded-lg shadow-lg p-6">
                    <h2 className="text-xl font-semibold mb-4">Meeting Transcript</h2>
                    <div className="text-center py-8 text-gray-500">
                        <div className="text-4xl mb-2">📝</div>
                        <p>Transcript will be available here</p>
                        <p className="text-sm mt-1">Coming soon...</p>
                    </div>
                </div>

                {/* Action Items Section */}
                <div className="bg-white rounded-lg shadow-lg p-6">
                    <h2 className="text-xl font-semibold mb-4">Action Items</h2>
                    <div className="text-center py-8 text-gray-500">
                        <div className="text-4xl mb-2">✅</div>
                        <p>AI-generated action items will appear here</p>
                        <p className="text-sm mt-1">Coming soon...</p>
                    </div>
                </div>
            </div>
        </div>
    );
} 