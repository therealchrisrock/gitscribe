import type { LoaderFunctionArgs } from "@remix-run/node";
import { json } from "@remix-run/node";
import { requireAuth } from "@/lib/auth.server";

export async function loader({ request }: LoaderFunctionArgs) {
    // Verify authentication
    const user = await requireAuth(request);

    // Extract userId from query params
    const url = new URL(request.url);
    const userId = url.searchParams.get('userId');

    if (!userId) {
        return json({ error: 'User ID is required' }, { status: 400 });
    }

    // Verify the user is requesting their own meetings
    if (userId !== user.uid) {
        return json({ error: 'Unauthorized' }, { status: 403 });
    }

    try {
        // Get the API base URL from environment or default to localhost for development
        const apiBaseUrl = process.env.API_BASE_URL || 'http://localhost:8080';

        // For now, we'll return mock data until the backend API is ready
        // TODO: Replace with actual API call to ${apiBaseUrl}/api/meetings?user_id=${userId}

        const mockMeetings = [
            {
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
            {
                id: "meeting_2",
                title: "Product Review",
                type: "google_meet",
                status: "completed",
                start_time: "2024-01-14T14:00:00Z",
                end_time: "2024-01-14T15:00:00Z",
                meeting_url: "https://meet.google.com/abc-def-ghi",
                recording_path: null,
                created_at: "2024-01-14T13:55:00Z",
                updated_at: "2024-01-14T15:00:00Z"
            },
            {
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
            {
                id: "meeting_4",
                title: "Weekly Sync",
                type: "generic",
                status: "scheduled",
                start_time: "2024-01-16T10:00:00Z",
                end_time: null,
                meeting_url: "https://example.com/meeting/123",
                recording_path: null,
                created_at: "2024-01-16T09:00:00Z",
                updated_at: "2024-01-16T09:00:00Z"
            }
        ];

        return json({
            meetings: mockMeetings,
            total: mockMeetings.length
        });

        // When backend is ready, use this instead:
        /*
        const backendResponse = await fetch(`${apiBaseUrl}/api/meetings?user_id=${userId}`, {
            headers: {
                'Content-Type': 'application/json',
                // Add any authentication headers needed for backend
            },
        });
        
        if (!backendResponse.ok) {
            throw new Error(`Backend API error: ${backendResponse.statusText}`);
        }
        
        const data = await backendResponse.json();
        return json(data);
        */

    } catch (error) {
        console.error('Error fetching meetings:', error);
        return json(
            { error: 'Failed to fetch meetings' },
            { status: 500 }
        );
    }
} 