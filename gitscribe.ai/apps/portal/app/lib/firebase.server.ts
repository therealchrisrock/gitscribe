import { initializeServerApp, deleteApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

export async function getServerAuth(request: Request) {
    // Extract ID token from Authorization header
    const authHeader = request.headers.get("Authorization");
    const authIdToken = authHeader?.startsWith("Bearer ")
        ? authHeader.split("Bearer ")[1]
        : undefined;

    // Initialize FirebaseServerApp with the ID token
    const serverApp = initializeServerApp(firebaseConfig, {
        authIdToken,
        // Auto-cleanup when request is garbage collected
        releaseOnDeref: request,
    });

    const serverAuth = getAuth(serverApp);

    // Wait for auth to be ready
    await serverAuth.authStateReady();

    return {
        auth: serverAuth,
        user: serverAuth.currentUser,
        cleanup: () => deleteApp(serverApp),
    };
} 