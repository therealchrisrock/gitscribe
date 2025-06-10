import type { LoaderFunctionArgs } from "@remix-run/node";

export async function loader({ request }: LoaderFunctionArgs) {
    // Return empty JSON response for Chrome DevTools
    return new Response("{}", {
        status: 200,
        headers: {
            "Content-Type": "application/json",
        },
    });
} 