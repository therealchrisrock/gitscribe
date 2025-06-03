import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useLoaderData,
} from "@remix-run/react";
import type { LinksFunction } from "@remix-run/node";
import { AuthProvider } from "@/lib/auth-context";
import { Providers } from "@/components/providers";

import "./tailwind.css";
import { Toaster } from "@workspace/ui/components/sonner";

export const links: LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
];
const fallbackMeta = {
  title: "GitScribe - AI-Powered Code Documentation",
  description: "Automatically generate and maintain high-quality code documentation using AI. GitScribe helps teams keep documentation in sync with their codebase."
};

export async function loader() {
  return {
    meta: fallbackMeta
  };
}

export function Layout({ children }: { children: React.ReactNode }) {
  const { meta } = useLoaderData<typeof loader>();
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="robots" content="index,follow" />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={meta.title} />
        <meta property="og:description" content={meta.description} />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content={meta.title} />
        <meta name="twitter:description" content={meta.description} />
        <Meta />
        <Links />
      </head>
      <body>
        <Providers>
          <AuthProvider>
            {children}
          </AuthProvider>
          <Toaster />
        </Providers>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return Outlet({}) as React.ReactElement;
}
