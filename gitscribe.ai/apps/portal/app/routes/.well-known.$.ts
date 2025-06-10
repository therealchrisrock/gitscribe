import type { LoaderFunctionArgs } from "@remix-run/node";

export async function loader({ request, params }: LoaderFunctionArgs) {
    const url = new URL(request.url);

    // Handle Chrome DevTools specific request
    if (url.pathname.includes("com.chrome.devtools.json")) {
        return new Response("{}", {
            status: 200,
            headers: {
                "Content-Type": "application/json",
            },
        });
    }

    // For other .well-known requests, return 404
    throw new Response("Not Found", { status: 404 });
} 