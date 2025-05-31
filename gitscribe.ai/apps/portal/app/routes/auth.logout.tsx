import type { ActionFunctionArgs } from "@remix-run/node";
import { destroySession, getSession } from "../lib/session.server";
import { redirect } from "@remix-run/node";

export async function action({ request }: ActionFunctionArgs) {
    const session = await getSession(request.headers.get("Cookie"));

    // Check if this is a fetch request (from JavaScript) or a form submission
    const acceptHeader = request.headers.get("Accept");
    const isJsonRequest = acceptHeader?.includes("application/json");

    if (isJsonRequest) {
        // Return JSON response for client-side handling
        return Response.json(
            { success: true, message: "Logged out successfully" },
            {
                headers: {
                    "Set-Cookie": await destroySession(session),
                },
            }
        );
    } else {
        // Traditional form submission - redirect
        return redirect("/", {
            headers: {
                "Set-Cookie": await destroySession(session),
            },
        });
    }
} 