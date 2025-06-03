import type { LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData, Await } from "@remix-run/react";
import { Suspense, useState } from "react";
import { requireAuth } from "@/lib/auth.server";
import { useAuth } from "@/lib/auth-context";
import { LogoutButton } from "@/components/logout-button";
import { RecordingWidget } from "@/components/recording-widget";
import { MeetingList } from "@/components/meeting-list";
import { CreateMeetingButton } from "@/components/create-meeting-button";

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

async function fetchMeetings(userId: string) {
    // Get the API base URL from environment or default to localhost for development
    const apiBaseUrl = process.env.API_BASE_URL || 'http://localhost:8080';

    // For now, we'll return mock data until the backend API is ready
    // TODO: Replace with actual API call to ${apiBaseUrl}/api/meetings?user_id=${userId}

    // Simulate network delay to show loading state
    await new Promise(resolve => setTimeout(resolve, 1000));

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

    return {
        meetings: mockMeetings,
        total: mockMeetings.length
    };

    // When backend is ready, use this instead:
    /*
    try {
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
        return data;
    } catch (error) {
        console.error('Error fetching meetings:', error);
        throw new Error('Failed to fetch meetings');
    }
    */
}

export async function loader({ request }: LoaderFunctionArgs) {
    const user = await requireAuth(request);

    // Return immediately available data and defer the meetings fetch
    return {
        serverUser: user,
        timestamp: new Date().toISOString(),
        meetingsPromise: fetchMeetings(user.uid), // This will be streamed
    }
}

export default function Dashboard() {
    const { serverUser, meetingsPromise } = useLoaderData<typeof loader>();
    const { user: clientUser } = useAuth();
    const [currentMeeting, setCurrentMeeting] = useState<Meeting | null>(null);

    const handleMeetingCreated = (meeting: Meeting) => {
        // Set the current meeting for recording
        setCurrentMeeting(meeting);
    };

    const handleSelectMeeting = (meeting: Meeting) => {
        // Allow selecting an existing meeting for recording
        setCurrentMeeting(meeting);
    };

    return (
        <div className="container mx-auto p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Dashboard</h1>
                <CreateMeetingButton onMeetingCreated={handleMeetingCreated} />
            </div>

            <div className="grid gap-6 md:grid-cols-2">
                {/* User Info Card */}
                <div className="bg-card p-6 rounded-lg border">
                    <h2 className="text-xl font-semibold mb-4">Welcome back!</h2>
                    <div className="space-y-2">
                        <p><strong>Server User Email:</strong> {serverUser.email}</p>
                        <p><strong>Server User UID:</strong> {serverUser.uid}</p>
                        {serverUser.displayName && (
                            <p><strong>Server Display Name:</strong> {serverUser.displayName}</p>
                        )}
                        <hr className="my-4" />
                        <p><strong>Client User Email:</strong> {clientUser?.email || "Loading..."}</p>
                        <p><strong>Client User UID:</strong> {clientUser?.uid || "Loading..."}</p>
                        {currentMeeting && (
                            <>
                                <hr className="my-4" />
                                <p><strong>Active Meeting:</strong> {currentMeeting.title}</p>
                                <p><strong>Meeting ID:</strong> {currentMeeting.id}</p>
                                <button
                                    onClick={() => setCurrentMeeting(null)}
                                    className="text-sm text-red-600 hover:text-red-800"
                                >
                                    Clear Active Meeting
                                </button>
                            </>
                        )}
                    </div>
                    <div className="mt-6 space-x-4">
                        <LogoutButton />

                        {/* Example of using LogoutButton with different variants and sizes */}
                        {/* <LogoutButton variant="outline" size="sm">Sign Out</LogoutButton> */}
                        {/* <LogoutButton variant="ghost" className="text-red-600">Exit</LogoutButton> */}

                        {/* Client-side Firebase operations work automatically */}
                        <button
                            onClick={async () => {
                                if (clientUser) {
                                    const token = await clientUser.getIdToken();
                                    console.log("Fresh token:", token);
                                    alert("Fresh token logged to console");
                                }
                            }}
                            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                        >
                            Get Fresh Token
                        </button>
                    </div>
                </div>

                {/* Recording Widget - Now conditional on having an active meeting */}
                <RecordingWidget meetingId={currentMeeting?.id} />
            </div>

            {/* Meeting List - Full width below the grid with deferred loading */}
            <div className="mt-6">
                <Suspense fallback={<MeetingListSkeleton />}>
                    <Await resolve={meetingsPromise}>
                        {(meetingsData) => (
                            <MeetingList
                                meetingsData={meetingsData}
                                onMeetingSelect={handleSelectMeeting}
                                currentMeeting={currentMeeting}
                            />
                        )}
                    </Await>
                </Suspense>
            </div>
        </div>
    );
}

// Loading skeleton component for the meeting list
function MeetingListSkeleton() {
    return (
        <div className="bg-card p-6 rounded-lg border">
            <h2 className="text-xl font-semibold mb-4">Past Meetings</h2>
            <div className="space-y-3">
                {[...Array(3)].map((_, i) => (
                    <div key={i} className="animate-pulse">
                        <div className="h-16 bg-gray-200 rounded-lg"></div>
                    </div>
                ))}
            </div>
        </div>
    );
} 