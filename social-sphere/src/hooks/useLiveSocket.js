"use client";

import { useEffect, useRef, useCallback, useState } from "react";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8082";

/**
 * WebSocket connection states
 */
export const ConnectionState = {
    CONNECTING: "connecting",
    CONNECTED: "connected",
    DISCONNECTED: "disconnected",
    RECONNECTING: "reconnecting",
};

/**
 * Custom hook for WebSocket connection to the live service
 * Handles real-time messaging for private and group chats
 *
 * @param {Object} options
 * @param {Function} options.onPrivateMessage - Callback when receiving a private message
 * @param {Function} options.onGroupMessage - Callback when receiving a group message
 * @param {Function} options.onError - Callback for error handling
 * @param {boolean} options.autoConnect - Whether to connect automatically (default: true)
 */
export function useLiveSocket({
    onPrivateMessage,
    onGroupMessage,
    onError,
    autoConnect = true,
} = {}) {
    const wsRef = useRef(null);
    const reconnectTimeoutRef = useRef(null);
    const reconnectAttemptsRef = useRef(0);
    const currentGroupSubRef = useRef(null);

    const [connectionState, setConnectionState] = useState(ConnectionState.DISCONNECTED);

    const onPrivateMessageRef = useRef(onPrivateMessage);
    const onGroupMessageRef = useRef(onGroupMessage);
    const onErrorRef = useRef(onError);

    // Keep refs updated
    useEffect(() => {
        onPrivateMessageRef.current = onPrivateMessage;
        onGroupMessageRef.current = onGroupMessage;
        onErrorRef.current = onError;
    }, [onPrivateMessage, onGroupMessage, onError]);

    /**
     * Connect to WebSocket
     */
    const connect = useCallback(() => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            return;
        }

        // Clean up existing connection
        if (wsRef.current) {
            wsRef.current.close();
        }

        setConnectionState(reconnectAttemptsRef.current > 0
            ? ConnectionState.RECONNECTING
            : ConnectionState.CONNECTING
        );

        try {
            const ws = new WebSocket(`${WS_URL}/live`);
            wsRef.current = ws;

            ws.onopen = () => {
                console.log("[WebSocket] Connected to live service");
                setConnectionState(ConnectionState.CONNECTED);
                reconnectAttemptsRef.current = 0;

                // Re-subscribe to previous group if any
                if (currentGroupSubRef.current) {
                    ws.send(`sub:${currentGroupSubRef.current}`);
                }
            };

            ws.onmessage = (event) => {
                try {
                    // Server sends JSON array of messages (batched)
                    const messages = JSON.parse(event.data);

                    if (!Array.isArray(messages)) {
                        console.warn("[WebSocket] Expected array, got:", typeof messages);
                        return;
                    }

                    for (const msg of messages) {
                        // Determine message type based on structure
                        if (msg.group_id) {
                            // Group message
                            onGroupMessageRef.current?.(msg);
                        } else if (msg.conversation_id) {
                            // Private message
                            onPrivateMessageRef.current?.(msg);
                        } else {
                            console.log("[WebSocket] Unknown message type:", msg);
                        }
                    }
                } catch (err) {
                    console.error("[WebSocket] Failed to parse message:", err, event.data);
                }
            };

            ws.onerror = (error) => {
                console.error("[WebSocket] Error:", error);
                onErrorRef.current?.(error);
            };

            ws.onclose = (event) => {
                console.log("[WebSocket] Disconnected:", event.code, event.reason);
                setConnectionState(ConnectionState.DISCONNECTED);

                // Attempt reconnection with exponential backoff
                if (event.code !== 1000) { // 1000 = normal closure
                    const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
                    reconnectAttemptsRef.current++;

                    console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current})`);

                    reconnectTimeoutRef.current = setTimeout(() => {
                        connect();
                    }, delay);
                }
            };
        } catch (err) {
            console.error("[WebSocket] Failed to create connection:", err);
            setConnectionState(ConnectionState.DISCONNECTED);
            onErrorRef.current?.(err);
        }
    }, []);

    /**
     * Disconnect from WebSocket
     */
    const disconnect = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }

        reconnectAttemptsRef.current = 0;
        currentGroupSubRef.current = null;

        if (wsRef.current) {
            wsRef.current.close(1000, "User initiated disconnect");
            wsRef.current = null;
        }
    }, []);

    /**
     * Subscribe to a group chat channel
     * @param {string} groupId - The group ID to subscribe to
     */
    const subscribeToGroup = useCallback((groupId) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            // Unsubscribe from previous group first
            if (currentGroupSubRef.current && currentGroupSubRef.current !== groupId) {
                wsRef.current.send(`unsub:${currentGroupSubRef.current}`);
            }
            wsRef.current.send(`sub:${groupId}`);
            currentGroupSubRef.current = groupId;
            console.log("[WebSocket] Subscribed to group:", groupId);
        } else {
            // Store for later when connected
            currentGroupSubRef.current = groupId;
            console.log("[WebSocket] Will subscribe to group when connected:", groupId);
        }
    }, []);

    /**
     * Unsubscribe from a group chat channel
     * @param {string} groupId - The group ID to unsubscribe from (optional, defaults to current)
     */
    const unsubscribeFromGroup = useCallback((groupId) => {
        const targetGroup = groupId || currentGroupSubRef.current;
        if (wsRef.current?.readyState === WebSocket.OPEN && targetGroup) {
            wsRef.current.send(`unsub:${targetGroup}`);
            console.log("[WebSocket] Unsubscribed from group:", targetGroup);
        }
        if (currentGroupSubRef.current === targetGroup) {
            currentGroupSubRef.current = null;
        }
    }, []);

    /**
     * Send a private message through WebSocket
     * @param {Object} message - Message object with interlocutor_id and message_text
     * @returns {boolean} - Whether the message was sent
     */
    const sendPrivateMessage = useCallback((message) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            const payload = JSON.stringify({
                interlocutor_id: message.interlocutor_id,
                message_text: message.message_text,
            });
            wsRef.current.send(`ch:${payload}`);
            console.log("[WebSocket] Sent private message");
            return true;
        }
        console.warn("[WebSocket] Cannot send message - not connected");
        return false;
    }, []);

    /**
     * Check if currently connected
     */
    const isConnected = connectionState === ConnectionState.CONNECTED;

    // Auto-connect on mount
    useEffect(() => {
        if (autoConnect) {
            connect();
        }

        return () => {
            disconnect();
        };
    }, [autoConnect, connect, disconnect]);

    return {
        connectionState,
        isConnected,
        connect,
        disconnect,
        subscribeToGroup,
        unsubscribeFromGroup,
        sendPrivateMessage,
    };
}
