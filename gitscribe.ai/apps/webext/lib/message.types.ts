import { SendMessageFn } from "@/lib/isomorphic-message";
import { OpenStreamFn } from "@/lib/isomorphic-message";
import { OnOpenStreamChannelFn } from "@/lib/isomorphic-message";
import { OnMessageFn } from "@/lib/isomorphic-message";

export interface MessagingFunctions {
    sendMessage: SendMessageFn;
    onMessage: OnMessageFn;
    openStream: OpenStreamFn;
    onOpenStreamChannel: OnOpenStreamChannelFn;
}
export enum Entrypoint {
    CONTENT_SCRIPT = "content-script",
    POPUP = "popup",
    BACKGROUND = "background",
    SIDE_PANEL = "side-panel"
}
