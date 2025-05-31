import { initializeApp, getApps, type FirebaseApp } from "firebase/app";
import { getAuth, type Auth } from "firebase/auth";
import { getFirestore, type Firestore } from "firebase/firestore";

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

let app: FirebaseApp;
let auth: Auth;
let db: Firestore;

function initializeFirebase() {
    if (!getApps().length) {
        app = initializeApp(firebaseConfig);
    } else {
        app = getApps()[0];
    }
    return app;
}

export function getFirebaseApp() {
    if (!app) {
        initializeFirebase();
    }
    return app;
}

export function getFirebaseAuth() {
    if (!auth) {
        const currentApp = getFirebaseApp();
        auth = getAuth(currentApp);
    }
    return auth;
}

export function getFirebaseDb() {
    if (!db) {
        const currentApp = getFirebaseApp();
        db = getFirestore(currentApp);
    }
    return db;
} 