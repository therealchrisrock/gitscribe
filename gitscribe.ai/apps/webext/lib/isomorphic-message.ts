import {
    onMessage as onMessageContentScript,
    sendMessage as sendMessageContentScript,
    openStream as openStreamContentScript,
    onOpenStreamChannel as onOpenStreamChannelContentScript
} from "webext-bridge/content-script";

import {
    onMessage as onMessagePopup,
    sendMessage as sendMessagePopup,
    openStream as openStreamPopup,
    onOpenStreamChannel as onOpenStreamChannelPopup
} from "webext-bridge/popup";

import {
    onMessage as onMessageBackground,
    sendMessage as sendMessageBackground,
    openStream as openStreamBackground,
    onOpenStreamChannel as onOpenStreamChannelBackground
} from "webext-bridge/background";

import {
    Entrypoint,
    MessagingFunctions
} from "@/lib/message.types";


export type SendMessageFn = typeof sendMessageContentScript | typeof sendMessagePopup | typeof sendMessageBackground;
export type OnMessageFn = typeof onMessageContentScript | typeof onMessagePopup | typeof onMessageBackground;
export type OpenStreamFn = typeof openStreamContentScript | typeof openStreamPopup | typeof openStreamBackground;
export type OnOpenStreamChannelFn = typeof onOpenStreamChannelContentScript | typeof onOpenStreamChannelPopup | typeof onOpenStreamChannelBackground;

export const getMessagingFunctions = (ctx: Entrypoint): MessagingFunctions => {
    switch (ctx) {
        case Entrypoint.CONTENT_SCRIPT:
            return {
                sendMessage: sendMessageContentScript,
                onMessage: onMessageContentScript,
                openStream: openStreamContentScript,
                onOpenStreamChannel: onOpenStreamChannelContentScript
            };
        case Entrypoint.POPUP:
            return {
                sendMessage: sendMessagePopup,
                onMessage: onMessagePopup,
                openStream: openStreamPopup,
                onOpenStreamChannel: onOpenStreamChannelPopup
            };
        case Entrypoint.BACKGROUND:
            return {
                sendMessage: sendMessageBackground,
                onMessage: onMessageBackground,
                openStream: openStreamBackground,
                onOpenStreamChannel: onOpenStreamChannelBackground
            };
        case Entrypoint.SIDE_PANEL:
            return {
                sendMessage: sendMessageContentScript,
                onMessage: onMessageContentScript,
                openStream: openStreamContentScript,
                onOpenStreamChannel: onOpenStreamChannelContentScript
            };
        default:
            throw new Error(`Invalid entrypoint context: ${ctx}`);
    }
};