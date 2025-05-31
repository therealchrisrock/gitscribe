// Firebase Auth Service Worker
// This intercepts navigation requests and adds the Firebase Auth token to headers

import { initializeApp } from 'firebase/app';
import { getAuth, getIdToken } from 'firebase/auth';

const firebaseConfig = {
    apiKey: 'AIzaSyB0-YFjFik2ENvJXlfshTUf9yFVfdCvFhY',
    authDomain: 'teammate-5dbc9.firebaseapp.com',
    projectId: 'teammate-5dbc9',
    storageBucket: 'teammate-5dbc9.appspot.com',
    messagingSenderId: '108471893127314145033',
    appId: '1:1086300260990:web:0e8a730175764a8860262f',
};

// Initialize Firebase in service worker
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

self.addEventListener('fetch', async (event) => {
    const { request } = event;

    // Only intercept navigation requests to our domain
    if (
        request.mode === 'navigate' ||
        (request.method === 'GET' && request.headers.get('accept')?.includes('text/html'))
    ) {
        event.respondWith(handleRequest(request));
    }
});

async function handleRequest(request) {
    try {
        // Get current user's ID token
        const user = auth.currentUser;
        let idToken = null;

        if (user) {
            try {
                // Force refresh to ensure token is valid
                idToken = await getIdToken(user, true);
            } catch (error) {
                console.warn('Failed to get ID token:', error);
            }
        }

        // Clone the request and add auth header if we have a token
        const headers = new Headers(request.headers);
        if (idToken) {
            headers.set('Authorization', `Bearer ${idToken}`);
        }

        const modifiedRequest = new Request(request, {
            headers,
        });

        return fetch(modifiedRequest);
    } catch (error) {
        console.error('Service worker error:', error);
        return fetch(request);
    }
} 