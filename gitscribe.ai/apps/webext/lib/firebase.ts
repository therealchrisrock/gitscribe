// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
    apiKey: "AIzaSyB0-YFjFik2ENvJXlfshTUf9yFVfdCvFhY",
    authDomain: "teammate-5dbc9.firebaseapp.com",
    projectId: "teammate-5dbc9",
    storageBucket: "teammate-5dbc9.firebasestorage.app",
    messagingSenderId: "1086300260990",
    appId: "1:1086300260990:web:85451894b20a9d0760262f",
    measurementId: "G-DYG9VE9Y7Y"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const analytics = getAnalytics(app);

export { app, analytics, auth };