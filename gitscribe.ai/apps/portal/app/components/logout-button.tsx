import { useNavigate } from "@remix-run/react";
import { Button } from "@workspace/ui/components/button";
import { useState, forwardRef } from "react";
import React from "react";
import { signOut } from "firebase/auth";
import { useAuth } from "@/lib/auth-context"; // Import useAuth


type LogoutButtonProps = Omit<React.ComponentProps<typeof Button>, 'onClick'> & {
    children?: string;
};

export const LogoutButton = forwardRef<HTMLButtonElement, LogoutButtonProps>(
    ({ variant = "destructive", children = "Logout", ...props }, ref) => {
        const navigate = useNavigate();
        const [isLoggingOut, setIsLoggingOut] = useState(false);
        const { getAuthInstance } = useAuth(); // Get auth instance from context

        const handleLogout = async () => {
            if (isLoggingOut) return;

            setIsLoggingOut(true);
            const authInstance = getAuthInstance(); // Call to get the instance

            try {
                // First, sign out from Firebase on the client
                await signOut(authInstance); // Use the instance from context

                // Then, destroy the server session
                const response = await fetch("/auth/logout", {
                    method: "POST",
                    headers: {
                        "Accept": "application/json",
                        "Content-Type": "application/json",
                    },
                });

                if (response.ok) {
                    // Navigate to home page after successful logout
                    navigate("/");
                } else {
                    console.error("Server logout failed");
                    // Still navigate even if server logout fails
                    navigate("/");
                }
            } catch (error) {
                console.error("Logout failed:", error);
                // Navigate anyway to prevent user from being stuck
                navigate("/");
            } finally {
                setIsLoggingOut(false);
            }
        };

        return React.createElement(Button, {
            ref,
            variant,
            disabled: isLoggingOut,
            onClick: handleLogout,
            ...props,
        }, isLoggingOut ? "Logging out..." : children);
    }
);

LogoutButton.displayName = "LogoutButton"; 