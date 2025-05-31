import { createContext, ReactNode, useContext } from "react";
import { getMessagingFunctions } from "@/lib/isomorphic-message";
import { Entrypoint, MessagingFunctions } from "@/lib/message.types";

type ReactEntrypoint = Exclude<Entrypoint, Entrypoint.BACKGROUND>;

interface EntrypointContextValue extends MessagingFunctions { }

const EntrypointContext = createContext<EntrypointContextValue | null>(null);

export function EntrypointProvider({ children, ctx }: { children: ReactNode, ctx: ReactEntrypoint }) {
    const messaging = getMessagingFunctions(ctx);
    return (
        <EntrypointContext.Provider value={{ ...messaging }}>
            {children}
        </EntrypointContext.Provider>
    );
}

export function useMessaging() {
    const context = useContext(EntrypointContext);
    if (!context) {
        throw new Error("useMessaging must be used within an EntrypointProvider");
    }
    return context;
}