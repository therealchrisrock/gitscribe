import { createCookieSessionStorage } from "@remix-run/node";

export const { getSession, commitSession, destroySession } =
    createCookieSessionStorage({
        cookie: {
            name: "__session",
            httpOnly: true,
            maxAge: 60 * 60 * 24 * 5, // 5 days
            path: "/",
            sameSite: "lax",
            secrets: [process.env.SESSION_SECRET!],
            secure: process.env.NODE_ENV === "production",
        },
    }); 