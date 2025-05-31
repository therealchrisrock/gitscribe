import type { LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { requireAuth } from "@/lib/auth.server";
import { useAuth } from "@/lib/auth-context";
import { LogoutButton } from "@/components/logout-button";
import { RecordingWidget } from "@/components/recording-widget";

export async function loader({ request }: LoaderFunctionArgs) {
    const user = await requireAuth(request);

    return {
        serverUser: user,
        timestamp: new Date().toISOString(),
    };
}

export default function Dashboard() {
    const { serverUser } = useLoaderData<typeof loader>();
    const { user: clientUser } = useAuth();

    return (
        <div className="container mx-auto p-6">
            <h1 className="text-3xl font-bold mb-6">Dashboard</h1>

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

                {/* Recording Widget */}
                <RecordingWidget />
            </div>
        </div>
    );
} 