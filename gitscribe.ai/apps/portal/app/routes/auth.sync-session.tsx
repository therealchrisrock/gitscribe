import type { ActionFunctionArgs } from "@remix-run/node";
import { getSession, commitSession } from "@/lib/session.server";
import { getAuth } from "firebase-admin/auth";
import { initializeApp, getApps, cert } from "firebase-admin/app";

export async function action({ request }: ActionFunctionArgs) {
    // Initialize Firebase Admin SDK (server-side only)
    if (!getApps().length) {
        const projectId = process.env.FIREBASE_PROJECT_ID;
        const clientEmail = process.env.FIREBASE_CLIENT_EMAIL;
        const rawPrivateKeyFromEnv = process.env.FIREBASE_PRIVATE_KEY;

        let processedPrivateKey;
        if (rawPrivateKeyFromEnv) {
            processedPrivateKey = rawPrivateKeyFromEnv
                .replace(/^["\']|["\']$/g, '') // Remove surrounding single or double quotes
                .replace(/\\n/g, '\n');      // Convert literal \n to actual newlines
        }

        if (!projectId || !clientEmail || !processedPrivateKey) {
            throw new Error('Missing or invalid Firebase Admin SDK environment variables');
        }

        initializeApp({
            credential: cert({
                projectId,
                clientEmail,
                privateKey: processedPrivateKey,
            }),
        });
    }

    const auth = getAuth();

    // Handle potential aborted requests gracefully
    if (request.signal?.aborted) {
        return Response.json({ error: "Request aborted" }, { status: 499 });
    }

    const formData = await request.formData();
    const idToken = formData.get("idToken");

    if (typeof idToken !== "string") {
        return Response.json({ error: "Invalid token" }, { status: 400 });
    }

    try {
        // Create session cookie (lasts 5 days, no refresh needed)
        const expiresIn = 60 * 60 * 24 * 5 * 1000; // 5 days
        const sessionCookie = await auth.createSessionCookie(idToken, {
            expiresIn,
        });

        const session = await getSession(request.headers.get("Cookie"));
        session.set("sessionCookie", sessionCookie);

        // Return success response instead of redirecting
        return Response.json(
            { success: true },
            {
                headers: {
                    "Set-Cookie": await commitSession(session),
                },
            }
        );
    } catch (error) {
        console.error("Session sync failed:", error);
        return Response.json({ error: "Session sync failed" }, { status: 400 });
    }
} 