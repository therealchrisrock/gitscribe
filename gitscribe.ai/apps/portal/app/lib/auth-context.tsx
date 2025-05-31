import {
    createContext,
    useContext,
    useEffect,
    useState,
    ReactNode,
    useCallback,
    useRef,
} from "react";
import { getFirebaseAuth } from "./firebase.client";
import { onAuthStateChanged, signOut, getIdToken, type User, type Auth } from "firebase/auth";
import { useFetcher } from "@remix-run/react";

interface AuthContextType {
    user: User | null;
    loading: boolean;
    refreshToken: () => Promise<void>;
    getAuthInstance: () => Auth;
}

const AuthContext = createContext<AuthContextType>({
    user: null,
    loading: true,
    refreshToken: async () => { },
    getAuthInstance: () => {
        if (typeof window === "undefined") {
            console.warn("getAuthInstance called on server in placeholder AuthProvider");
            return getFirebaseAuth();
        }
        return getFirebaseAuth();
    },
});

function ClientOnlyAuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const fetcher = useFetcher();
    const syncInProgress = useRef(false);
    const refreshInterval = useRef<NodeJS.Timeout | null>(null);
    const lastSyncTime = useRef<number>(0);

    const authInstance = getFirebaseAuth();

    const refreshToken = useCallback(async () => {
        if (!user || syncInProgress.current || fetcher.state !== "idle") return;

        // Prevent too frequent refreshes (minimum 30 seconds between refreshes)
        const now = Date.now();
        if (now - lastSyncTime.current < 30000) return;

        try {
            syncInProgress.current = true;
            lastSyncTime.current = now;
            const idToken = await user.getIdToken(true); // Force refresh

            fetcher.submit(
                { idToken },
                { method: "post", action: "/auth/sync-session" }
            );
        } catch (error) {
            console.error("Token refresh failed:", error);
        } finally {
            // Reset sync flag after a delay to allow fetcher to complete
            setTimeout(() => {
                syncInProgress.current = false;
            }, 1000);
        }
    }, [user, fetcher]);

    const syncSession = useCallback(async (firebaseUser: User) => {
        if (syncInProgress.current || fetcher.state !== "idle") return;

        // Prevent too frequent syncs
        const now = Date.now();
        if (now - lastSyncTime.current < 5000) return;

        try {
            syncInProgress.current = true;
            lastSyncTime.current = now;
            const idToken = await firebaseUser.getIdToken();
            fetcher.submit(
                { idToken },
                { method: "post", action: "/auth/sync-session" }
            );
        } catch (error) {
            console.error("Session sync failed:", error);
        } finally {
            // Reset sync flag after a delay to allow fetcher to complete
            setTimeout(() => {
                syncInProgress.current = false;
            }, 1000);
        }
    }, [fetcher]);

    const clearSession = useCallback(async () => {
        if (syncInProgress.current) return;

        try {
            syncInProgress.current = true;

            await fetch("/auth/logout", {
                method: "POST",
                headers: {
                    "Accept": "application/json",
                    "Content-Type": "application/json",
                },
            });
        } catch (error) {
            console.error("Failed to clear server session:", error);
        } finally {
            // Reset sync flag after a delay
            setTimeout(() => {
                syncInProgress.current = false;
            }, 1000);
        }
    }, []);

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(authInstance, async (firebaseUser) => {
            setUser(firebaseUser);
            setLoading(false);

            if (firebaseUser) {
                // Sync fresh token with server session
                await syncSession(firebaseUser);
            } else {
                // User logged out, clear server session
                clearSession();
            }
        });

        return unsubscribe;
    }, [syncSession, clearSession, authInstance]);

    // Auto-refresh tokens before they expire
    useEffect(() => {
        // Clear existing interval
        if (refreshInterval.current) {
            clearInterval(refreshInterval.current);
            refreshInterval.current = null;
        }

        if (!user) return;

        refreshInterval.current = setInterval(async () => {
            await refreshToken();
        }, 50 * 60 * 1000); // Refresh every 50 minutes

        return () => {
            if (refreshInterval.current) {
                clearInterval(refreshInterval.current);
                refreshInterval.current = null;
            }
        };
    }, [user, refreshToken]);

    return (
        <AuthContext.Provider value={{ user, loading, refreshToken, getAuthInstance: () => authInstance }}>
            {children}
        </AuthContext.Provider>
    );
}

export function AuthProvider({ children }: { children: ReactNode }) {
    const [isClient, setIsClient] = useState(false);

    useEffect(() => {
        setIsClient(true);
    }, []);

    if (!isClient) {
        return (
            <AuthContext.Provider value={{
                user: null,
                loading: true,
                refreshToken: async () => { },
                getAuthInstance: () => {
                    if (typeof window === "undefined") {
                        console.warn("getAuthInstance called on server in placeholder AuthProvider");
                        return getFirebaseAuth();
                    }
                    return getFirebaseAuth();
                }
            }}>
                {children}
            </AuthContext.Provider>
        );
    }

    // On the client, use the full Firebase auth provider
    return <ClientOnlyAuthProvider>{children}</ClientOnlyAuthProvider>;
}

export const useAuth = () => useContext(AuthContext); 