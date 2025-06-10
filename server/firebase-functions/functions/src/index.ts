/**
 * Import function triggers from their respective submodules:
 *
 * import {onCall} from "firebase-functions/v2/https";
 * import {onDocumentWritten} from "firebase-functions/v2/firestore";
 *
 * See a full list of supported triggers at https://firebase.google.com/docs/functions
 */

import { onCall } from "firebase-functions/v2/https";
import * as functions from "firebase-functions/v1";
import * as logger from "firebase-functions/logger";
import { defineString } from "firebase-functions/params";
import { initializeApp } from "firebase-admin/app";
import { UserRecord } from "firebase-admin/auth";

// Initialize Firebase Admin
initializeApp();

// Start writing functions
// https://firebase.google.com/docs/functions/typescript

// export const helloWorld = onRequest((request, response) => {
//   logger.info("Hello logs!", {structuredData: true});
//   response.send("Hello from Firebase!");
// });

/**
 * Firebase Functions for GitScribe
 * Handles user registration by calling our Go server's RPC endpoint
 */

// Configuration from environment
const SERVER_URL = defineString("SERVER_URL", { default: "http://localhost:8080" });

/**
 * User data interface for registration
 */
interface UserData {
    id: string;
    name: string;
    email: string;
}

/**
 * Triggered automatically when a new user is created in Firebase Auth
 * Registers the user in our Go server using the /register endpoint
 * 
 * This uses v1 functions syntax because auth triggers are only available in v1
 */
export const createUserInDatabase = functions.auth.user().onCreate(async (user: UserRecord) => {
    logger.info("New user created", { uid: user.uid });

    try {
        // Prepare user data for registration
        const userData: UserData = {
            id: user.uid,
            name: user.displayName || extractNameFromEmail(user.email || ""),
            email: user.email || "",
        };

        logger.info("Registering user", { userData });

        // Call our Go server's registration endpoint
        const response = await fetch(`${SERVER_URL.value()}/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "User-Agent": "Firebase-Functions/1.0",
            },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP ${response.status}: ${errorText}`);
        }

        const responseData = await response.json();
        logger.info("Registration success", { response: responseData });

        return { success: true, userId: user.uid };
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : "Unknown error";
        logger.error("Registration failed", {
            userId: user.uid,
            email: user.email,
            error: errorMessage,
        });

        // Don't throw - we don't want to fail Firebase user creation
        return { success: false, userId: user.uid, error: errorMessage };
    }
});

/**
 * Extract a name from email if displayName is not available
 */
function extractNameFromEmail(email: string): string {
    if (!email) return "Unknown User";

    const localPart = email.split("@")[0];
    return localPart
        .replace(/[._-]/g, " ")
        .split(" ")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        .join(" ");
}

/**
 * Health check function for testing (v2 function)
 */
export const healthCheck = onCall(async (request) => {
    logger.info("Health check called", { data: request.data });
    return {
        status: "healthy",
        timestamp: new Date().toISOString(),
        serverUrl: SERVER_URL.value(),
    };
});
