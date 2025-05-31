import { createRoot } from "react-dom/client";
import { StrictMode } from 'react';
import { defineContentScript, createShadowRootUi } from '#imports';
import { EntrypointProvider } from '@/components/providers/entrypoint-provider';
import { Entrypoint } from "@/lib/message.types";
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '@/lib/query-client';
import App from "@/entrypoints/widget/App";

export default defineContentScript({
    matches: ["*://*/*", '*://*.google.com/*'],
    cssInjectionMode: "ui",
    async main(ctx) {
        const ui = await createShadowRootUi(ctx, {
            name: "gitscribe-widget",
            inheritStyles: false,
            position: 'inline',
            anchor: 'body',
            onMount: (container) => {
                const wrapper = document.createElement("div");
                container.append(wrapper);
                const root = createRoot(wrapper);
                root.render(
                    <StrictMode>
                        <QueryClientProvider client={queryClient}>
                            <EntrypointProvider ctx={Entrypoint.CONTENT_SCRIPT}>
                                <App />
                            </EntrypointProvider>
                        </QueryClientProvider>
                    </StrictMode>,
                );
                return { root, wrapper };
            },
            onRemove: (elements) => {
                elements?.root.unmount();
                elements?.wrapper.remove();
            },
        });
        await ui.mount();
    },
});
