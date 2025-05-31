import { ReactNode } from 'react';
import { redirect } from '@remix-run/node';
import { useAuth } from '@/lib/auth-context';

interface ProtectedRouteProps {
    children: ReactNode;
    fallback?: ReactNode;
}

export function ProtectedRoute({ children, fallback }: ProtectedRouteProps) {
    const { user, loading } = useAuth();

    // Show loading state while checking authentication
    if (loading) {
        return (
            fallback || (
                <div className="flex h-screen items-center justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                </div>
            )
        );
    }

    // Redirect to login if not authenticated
    if (!user) {
        // In Remix, client-side redirects should be handled differently
        // This component should primarily be used for UI protection
        // Server-side protection should be handled in loaders
        window.location.href = '/login';
        return null;
    }

    return <>{children}</>;
} 