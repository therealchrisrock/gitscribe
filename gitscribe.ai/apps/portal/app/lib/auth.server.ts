import { redirect } from "@remix-run/node";
import { getAuth } from "firebase-admin/auth";
import { initializeApp, getApps, cert } from "firebase-admin/app";
import { getSession, commitSession, destroySession } from "./session.server";

// Initialize Firebase Admin SDK
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

    console.log('Firebase Admin SDK Config:', {
        projectId: projectId ? 'SET' : 'MISSING',
        clientEmail: clientEmail ? 'SET' : 'MISSING',
        privateKey: processedPrivateKey ? 'SET (processed)' : 'MISSING'
    });

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

export async function createUserSession(
    request: Request,
    idToken: string,
    redirectTo: string = "/"
) {
    try {
        // Create session cookie (lasts 5 days, no refresh needed)
        const expiresIn = 60 * 60 * 24 * 5 * 1000; // 5 days
        const sessionCookie = await auth.createSessionCookie(idToken, {
            expiresIn,
        });

        const session = await getSession(request.headers.get("Cookie"));
        session.set("sessionCookie", sessionCookie);

        return redirect(redirectTo, {
            headers: {
                "Set-Cookie": await commitSession(session),
            },
        });
    } catch (error) {
        throw new Error("Invalid token");
    }
}

export async function requireAuth(request: Request) {
    const session = await getSession(request.headers.get("Cookie"));
    const sessionCookie = session.get("sessionCookie");

    if (!sessionCookie) {
        throw redirect("/login");
    }

    try {
        const decodedClaims = await auth.verifySessionCookie(
            sessionCookie,
            true // checkRevoked
        );

        return {
            uid: decodedClaims.uid,
            email: decodedClaims.email,
            displayName: decodedClaims.name,
            claims: decodedClaims,
        };
    } catch (error) {
        // Session expired/invalid, clear it
        const session = await getSession(request.headers.get("Cookie"));
        throw redirect("/login", {
            headers: {
                "Set-Cookie": await destroySession(session),
            },
        });
    }
}

export async function getOptionalUser(request: Request) {
    try {
        return await requireAuth(request);
    } catch {
        return null;
    }
}

export async function logout(request: Request) {
    const session = await getSession(request.headers.get("Cookie"));

    return redirect("/", {
        headers: {
            "Set-Cookie": await destroySession(session),
        },
    });
} 