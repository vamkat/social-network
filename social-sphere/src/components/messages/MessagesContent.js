"use client";

import { useEffect, useState, useRef, useCallback, useMemo } from "react";
import { useRouter } from "next/navigation";
import { getConv } from "@/actions/chat/get-conv";
import { getMessages } from "@/actions/chat/get-messages";
import { getConvByID } from "@/actions/chat/get-conv-by-id";
import { useStore } from "@/store/store";
import { User, Send, MessageCircle, Loader2, ChevronLeft, Wifi, WifiOff, Smile } from "lucide-react";
import { motion } from "motion/react";
import { useLiveSocket, ConnectionState } from "@/context/LiveSocketContext";
import EmojiPicker from "emoji-picker-react";
import { useMsgReceiver } from "@/store/store";

export default function MessagesContent({
    initialConversations = [],
    initialSelectedId = null,
    initialMessages = [],
    firstMessage = false,
}) {
    const router = useRouter();
    const user = useStore((state) => state.user);
    const decrementUnreadCount = useStore((state) => state.decrementUnreadCount);
    const [messages, setMessages] = useState(initialMessages);
    const [isLoadingConversations, setIsLoadingConversations] = useState(false);
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [hasMoreConvs, setHasMoreConvs] = useState(() => initialConversations.length >= 15);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [isLoadingMoreMessages, setIsLoadingMoreMessages] = useState(false);
    const [hasMoreMessages, setHasMoreMessages] = useState(() => initialMessages.length >= 20);
    const [messageText, setMessageText] = useState("");
    const [showMobileChat] = useState(!!initialSelectedId);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const selectedConvRef = useRef(null);
    const isLoadingMoreRef = useRef(false);
    const emojiPickerRef = useRef(null);
    const receiver = useMsgReceiver((state) => state.msgReceiver);
    const clearMsgReceiver = useMsgReceiver((state) => state.clearMsgReceiver);
    const [conversations, setConversations] = useState(() => {
        if (firstMessage && receiver) {
            const newConv = {
                Interlocutor: {
                    id: receiver.id,
                    username: receiver.username,
                    avatar_url: receiver.avatar_url
                }
            };
            return [newConv];
        }
        return initialConversations;
    });

    // Only sync initialConversations if NOT in firstMessage mode
    useEffect(() => {
        if (!firstMessage && initialConversations.length > 0) {
            setConversations(initialConversations);
        }
    }, [initialConversations, firstMessage]);

    //Clear receiver after adding conversation
    useEffect(() => {
        if (firstMessage && receiver) {
            clearMsgReceiver();
        }
    }, [firstMessage, receiver, clearMsgReceiver]);

    // Find selected conversation from ID
    const selectedConv = useMemo(() => {
        if (!initialSelectedId) return null;
        const selected = conversations.find(
            (conv) => conv.Interlocutor?.id === initialSelectedId
        ) || null;
        if (selected) {
            selected.UnreadCount = 0;
        }
        return selected;
    }, [initialSelectedId, conversations]);

    // Keep selectedConv ref in sync for WebSocket callback
    useEffect(() => {
        selectedConvRef.current = selectedConv;
    }, [selectedConv]);

    // WebSocket connection from context
    const { connectionState, isConnected, addOnPrivateMessage, removeOnPrivateMessage, sendPrivateMessage } = useLiveSocket();

    // Handle incoming private messages from WebSocket
    const handlePrivateMessage = useCallback(async (msg) => {
        console.log("[Chat] Received private message:", msg);

        const senderId = msg.sender?.id;
        const isOwnMessage = senderId === user?.id;

        // Add message to the current conversation if it matches
        const currentConv = selectedConvRef.current;
        console.log("Current: ", currentConv)
        if (currentConv) {
            const interlocutorId = currentConv.Interlocutor?.id;

            if (isOwnMessage) {
                // This is a confirmation of our sent message - replace pending with confirmed
                setMessages((prev) => {
                    // Find the pending message with matching text (most recent)
                    const pendingIndex = prev.findIndex(
                        (m) => m._pending && m.message_text === msg.message_text
                    );

                    if (pendingIndex !== -1) {
                        // Replace the pending message with the confirmed one
                        const updated = [...prev];
                        updated[pendingIndex] = { ...msg, _pending: false };
                        return updated;
                    }

                    // If no pending message found (edge case), add if not duplicate
                    if (prev.some((m) => m.id === msg.id)) return prev;
                    return [...prev, msg];
                });
            } else if (senderId === interlocutorId) {
                // Message from the other person in this conversation
                setMessages((prev) => {
                    // Prevent duplicates
                    if (prev.some((m) => m.id === msg.id)) return prev;
                    return [...prev, msg];
                });
            }
        }

        // Track if we need to fetch new conversation data
        let isNewConversation = false;

        // Update conversation list with new message preview (only for incoming messages)
        // Skip for own messages since handleSendMessage already updates the conversation
        if (!isOwnMessage) {
            setConversations((prev) => {
                // Check if conversation exists
                const existingIndex = prev.findIndex((conv) => conv.Interlocutor?.id === senderId);

                if (existingIndex !== -1) {
                    console.log("conversation exists")
                    // Update existing conversation
                    const updated = prev.map((conv, idx) => {
                        if (idx === existingIndex) {
                            return {
                                ...conv,
                                LastMessage: {
                                    message_text: msg.message_text,
                                    sender: msg.sender,
                                },
                                UpdatedAt: msg.created_at,
                                // Increment unread count for incoming messages
                                UnreadCount: (conv.UnreadCount || 0) + 1,
                            };
                        }
                        return conv;
                    });
                    return updated.sort((a, b) => new Date(b.UpdatedAt) - new Date(a.UpdatedAt));
                } else {
                    // New conversation - mark for fetching full data
                    isNewConversation = true;
                    console.log("New")
                    return prev;
                }
            });
        }

        // If new conversation from someone else, fetch full data from server
        if (isNewConversation && !isOwnMessage) {
            console.log("Fetching new conversation with:", { senderId, convId: msg.conversation_id });

            const result = await getConvByID({
                interlocutorId: senderId,
                convId: msg.conversation_id,
            });

            if (result.success && result.data) {
                setConversations((prev) => {
                    // Check if already added (race condition prevention)
                    const alreadyExists = prev.some((conv) => conv.Interlocutor?.id === senderId);
                    console.log("Already exists in conversations:", alreadyExists);
                    if (alreadyExists) {
                        return prev;
                    }
                    // Add the new conversation with full data and set UnreadCount
                    const newConv = {
                        ...result.data,
                        UnreadCount: 1,
                    };
                    console.log("Adding new conversation:", newConv);
                    return [newConv, ...prev];
                });
            }
        }
    }, [user?.id]);

    // Register message handler when on messages page
    useEffect(() => {
        addOnPrivateMessage(handlePrivateMessage);
        return () => removeOnPrivateMessage(handlePrivateMessage);
    }, [addOnPrivateMessage, removeOnPrivateMessage, handlePrivateMessage]);

    // Scroll to bottom of messages
    const scrollToBottom = (instant = false) => {
        if (instant) {
            messagesContainerRef.current?.scrollTo(0, messagesContainerRef.current.scrollHeight);
        } else {
            messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
        }
    };

    // Scroll to bottom only on first load of conversation
    useEffect(() => {
        if (messages.length > 0) {
            scrollToBottom(true);
        }
    }, [initialSelectedId]);

    // const handleClickMsg = async () => {
    //     if (!messages.length || !selectedConv) return;

    //     // Only mark as read if there are unread messages
    //     if (!selectedConv.UnreadCount || selectedConv.UnreadCount === 0) return;

    //     const lastMsg = messages[messages.length - 1];

    //     const res = await markAsRead({ convID: lastMsg.conversation_id, lastMsgID: lastMsg.id });
    //     if (!res.success) {
    //         return;
    //     }

    //     // Update the unread count for this conversation
    //     setConversations((prev) =>
    //         prev.map((conv) =>
    //             conv.ConversationId === lastMsg.conversation_id ||
    //                 conv.Interlocutor?.id === selectedConv?.Interlocutor?.id
    //                 ? { ...conv, UnreadCount: 0 }
    //                 : conv
    //         )
    //     );
    // };

    // Format relative time (compact)
    const formatRelativeTime = (dateString) => {
        const date = new Date(dateString);
        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return "now";
        if (diffMins < 60) return `${diffMins}m`;
        if (diffHours < 24) return `${diffHours}h`;
        if (diffDays < 7) return `${diffDays}d`;
        return date.toLocaleDateString();
    };

    // Format message time
    const formatMessageTime = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    };

    // Truncate message text
    const truncateMessage = (text, maxLength = 40) => {
        if (!text) return "";
        return text.length > maxLength ? text.substring(0, maxLength) + "..." : text;
    };

    // Check if conversation has unread messages for current user
    const hasUnreadMessages = (conv) => {
        return conv.UnreadCount > 0 && conv.LastMessage?.sender?.id !== user?.id;
    };

    // Load conversations
    const loadConversations = async () => {
        setIsLoadingConversations(true);
        try {
            const result = await getConv({ first: true, limit: 15 });
            if (result.success && result.data) {
                setConversations(result.data);
                setHasMoreConvs(result.data.length >= 15);
            }
        } catch (error) {
            console.error("Error loading conversations:", error);
        } finally {
            setIsLoadingConversations(false);
        }
    };

    // Load more conversations (pagination)
    const loadMoreConversations = async () => {
        if (isLoadingMore || !hasMoreConvs || conversations.length === 0) return;

        setIsLoadingMore(true);
        try {
            // Conversations are displayed newest-first, so the oldest is at the end
            // We need to find the minimum UpdatedAt to get conversations before it
            const oldestConv = conversations[conversations.length - 1];
            console.log("OLDEST: ", oldestConv);
            const beforeDate = oldestConv.UpdatedAt;
            console.log("BeforeDate: ", beforeDate);
            const result = await getConv({ first: false, beforeDate, limit: 5 });
            if (result.success && result.data) {
                setConversations((prev) => [...prev, ...result.data]);
                setHasMoreConvs(result.data.length >= 5);
            }
        } catch (error) {
            console.error("Error loading more conversations:", error);
        } finally {
            setIsLoadingMore(false);
        }
    };

    // Load messages for selected conversation
    const loadMessages = async (interlocutorId) => {
        setIsLoadingMessages(true);
        try {
            const result = await getMessages({ interlocutorId, limit: 20 });
            console.log(result);
            if (result.success && result.data) {
                // Messages come in reverse order (newest first), so reverse them
                const msgs = result.data.Messages?.reverse() || [];
                setMessages(msgs);
                setHasMoreMessages(msgs.length >= 20);
            }
        } catch (error) {
            console.error("Error loading messages:", error);
        } finally {
            setIsLoadingMessages(false);
        }
    };

    // Load more messages (pagination)
    const loadMoreMessages = useCallback(async () => {
        if (isLoadingMoreMessages || !hasMoreMessages || messages.length === 0 || !selectedConv) return;

        isLoadingMoreRef.current = true;
        setIsLoadingMoreMessages(true);

        const container = messagesContainerRef.current;
        const prevScrollHeight = container?.scrollHeight || 0;

        try {
            // Messages are displayed oldest-first, so messages[0] is the oldest
            const oldestMsg = messages[0];
            const boundary = oldestMsg.id;

            const result = await getMessages({
                interlocutorId: selectedConv.Interlocutor.id,
                boundary,
                limit: 10,
            });

            if (result.success && result.data) {
                const olderMsgs = result.data.Messages?.reverse() || [];
                setMessages((prev) => [...olderMsgs, ...prev]);
                setHasMoreMessages(olderMsgs.length >= 10);

                // Preserve scroll position after prepending
                requestAnimationFrame(() => {
                    if (container) {
                        const newScrollHeight = container.scrollHeight;
                        container.scrollTop = newScrollHeight - prevScrollHeight;
                    }
                });
            }
        } catch (error) {
            console.error("Error loading more messages:", error);
        } finally {
            setIsLoadingMoreMessages(false);
            isLoadingMoreRef.current = false;
        }
    }, [isLoadingMoreMessages, hasMoreMessages, messages, selectedConv]);

    // Handle scroll for infinite loading (called via onScroll prop)
    const handleMessagesScroll = (e) => {
        // Use ref (synchronous) to prevent multiple calls
        if (isLoadingMoreRef.current || !hasMoreMessages) return;

        const container = e.target;
        if (container.scrollTop < 100) {
            loadMoreMessages();
        }
    };

    // Handle conversation selection - navigate to /messages/[id]
    const handleSelectConversation = (conv) => {
        // Decrement unread count if this conversation has unread messages from someone else
        if (conv.UnreadCount > 0 && conv.LastMessage?.sender?.id !== user?.id) {
            decrementUnreadCount();
        }
        const id = conv.Interlocutor?.id;
        router.push(`/messages/${id}`);
    };

    // Handle send message - uses WebSocket
    const handleSendMessage = async (e) => {
        e.preventDefault();
        if (!messageText.trim() || !selectedConv || !isConnected) return;

        const msgToSend = messageText.trim();
        setMessageText("");

        // Generate a temporary ID to track this optimistic message
        const tempId = `temp-${Date.now()}`;

        // Optimistically add message to UI with pending state (will show with low opacity)
        const optimisticMessage = {
            id: tempId,
            message_text: msgToSend,
            sender: { id: user?.id },
            created_at: new Date().toISOString(),
            _pending: true, // Flag for showing low opacity until confirmed
        };
        setMessages((prev) => [...prev, optimisticMessage]);

        // Update conversation's last message preview immediately
        setConversations((prev) =>
            prev.map((c) =>
                c.ConversationId === selectedConv.ConversationId ||
                    c.Interlocutor?.id === selectedConv.Interlocutor?.id
                    ? {
                        ...c,
                        LastMessage: {
                            ...c.LastMessage,
                            message_text: msgToSend,
                            sender: { id: user?.id },
                        },
                        UpdatedAt: new Date().toISOString(),
                    }
                    : c
            ).sort((a, b) => new Date(b.UpdatedAt) - new Date(a.UpdatedAt))
        );

        try {
            await sendPrivateMessage(selectedConv.Interlocutor.id, msgToSend);
            // Server will send the confirmed message back through WebSocket
            // The handlePrivateMessage callback will update the message to remove _pending
        } catch (error) {
            console.error("Error sending message:", error);
            // Remove optimistic message and restore text on WebSocket error
            setMessages((prev) => prev.filter((m) => m.id !== tempId));
            setMessageText(msgToSend);
        }
    };

    // Handle back button on mobile - navigate to /messages
    const handleBackToList = () => {
        router.push("/messages");
    };

    // Handle emoji selection
    const onEmojiClick = (emojiData) => {
        setMessageText((prev) => prev + emojiData.emoji);
        setShowEmojiPicker(false);
    };

    // Close emoji picker when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (emojiPickerRef.current && !emojiPickerRef.current.contains(event.target)) {
                setShowEmojiPicker(false);
            }
        };

        if (showEmojiPicker) {
            document.addEventListener("mousedown", handleClickOutside);
        }

        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [showEmojiPicker]);

    // Load conversations on mount if not provided (skip for firstMessage mode)
    useEffect(() => {
        if (initialConversations.length === 0 && !firstMessage) {
            loadConversations();
        }
    }, []);

    return (
        <div className="h-[calc(100vh-5rem)] flex bg-background">
            {/* Left Sidebar - Conversations List */}
            <div
                className={`w-full md:w-80 lg:w-96 border-r border-(--border) flex flex-col ${showMobileChat ? "hidden md:flex" : "flex"
                    }`}
            >
                {/* Header */}
                <div className="p-4 border-b border-(--border) flex items-center justify-between">
                    <h1 className="text-xl font-bold text-foreground">Messages</h1>
                    {/* Connection Status Indicator */}
                    <div
                        className={`flex items-center gap-1.5 px-2 py-1 rounded-full text-xs ${isConnected
                            ? "bg-green-500/10 text-green-600"
                            : connectionState === ConnectionState.CONNECTING ||
                                connectionState === ConnectionState.RECONNECTING
                                ? "bg-yellow-500/10 text-yellow-600"
                                : "bg-red-500/10 text-red-500"
                            }`}
                        title={
                            isConnected
                                ? "Connected - Real-time updates active"
                                : connectionState === ConnectionState.RECONNECTING
                                    ? "Reconnecting..."
                                    : "Disconnected - Messages may be delayed"
                        }
                    >
                        {isConnected ? (
                            <Wifi className="w-3.5 h-3.5" />
                        ) : connectionState === ConnectionState.CONNECTING ||
                            connectionState === ConnectionState.RECONNECTING ? (
                            <Loader2 className="w-3.5 h-3.5 animate-spin" />
                        ) : (
                            <WifiOff className="w-3.5 h-3.5" />
                        )}
                        <span className="hidden sm:inline">
                            {isConnected ? "Live" : connectionState === ConnectionState.RECONNECTING ? "Reconnecting" : "Offline"}
                        </span>
                    </div>
                </div>

                {/* Conversations List */}
                <div className="flex-1 overflow-y-auto">
                    {isLoadingConversations ? (
                        <div className="flex items-center justify-center py-12">
                            <Loader2 className="w-6 h-6 text-(--muted) animate-spin" />
                        </div>
                    ) : conversations.length > 0 ? (
                        <div>
                            {conversations.map((conv) => (
                                <button
                                    key={conv.Interlocutor.id}
                                    onClick={() => handleSelectConversation(conv)}
                                    className={`w-full flex items-start gap-3 px-4 py-3 transition-colors cursor-pointer text-left border-b border-(--border)/50 ${selectedConv?.ConversationId === conv.ConversationId
                                        ? "bg-(--accent)/10"
                                        : "hover:bg-(--muted)/5"
                                        }`}
                                >
                                    {/* Avatar */}
                                    <div className="relative shrink-0">
                                        <div className="w-12 h-12 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden border border-(--border)">
                                            {conv.Interlocutor?.avatar_url ? (
                                                <img
                                                    src={conv.Interlocutor.avatar_url}
                                                    alt={conv.Interlocutor.username || "User"}
                                                    className="w-full h-full object-cover"
                                                />
                                            ) : (
                                                <User className="w-6 h-6 text-(--muted)" />
                                            )}
                                        </div>
                                        {hasUnreadMessages(conv) && (
                                            <span className="absolute -top-1 -right-1 min-w-5 h-5 px-1.5 text-[10px] font-bold text-white bg-red-500 rounded-full flex items-center justify-center border-2 border-background">
                                                {conv.UnreadCount}
                                            </span>
                                        )}
                                    </div>

                                    {/* Content */}
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center justify-between gap-2">
                                            <p
                                                className={`text-sm truncate ${hasUnreadMessages(conv)
                                                    ? "font-semibold text-foreground"
                                                    : "font-medium text-foreground"
                                                    }`}
                                            >
                                                {conv.Interlocutor?.username || "Unknown User"}
                                            </p>
                                            {conv?.UpdatedAt ? (
                                                <span className="text-xs text-(--muted) shrink-0">
                                                    {formatRelativeTime(conv.UpdatedAt)}
                                                </span>
                                            ) : <></>}

                                        </div>
                                        <p
                                            className={`text-sm mt-0.5 truncate ${hasUnreadMessages(conv) ? "text-foreground" : "text-(--muted)"
                                                }`}
                                        >
                                            {conv.LastMessage?.sender?.id === user?.id ? "You: " : ""}
                                            {truncateMessage(conv.LastMessage?.message_text)}
                                        </p>
                                    </div>
                                </button>
                            ))}
                            {/* Load More Button */}
                            {hasMoreConvs && (
                                <button
                                    onClick={loadMoreConversations}
                                    disabled={isLoadingMore}
                                    className="w-full py-3 text-sm text-(--accent) hover:bg-(--muted)/5 transition-colors disabled:opacity-50"
                                >
                                    {isLoadingMore ? (
                                        <Loader2 className="w-4 h-4 mx-auto animate-spin" />
                                    ) : (
                                        "Load more"
                                    )}
                                </button>
                            )}
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center py-12 px-4">
                            <MessageCircle className="w-12 h-12 text-(--muted) mb-3 opacity-30" />
                            <p className="text-(--muted) text-center">No conversations yet</p>
                            <p className="text-(--muted) text-sm text-center mt-1">
                                Start chatting with someone!
                            </p>
                        </div>
                    )}
                </div>
            </div>

            {/* Main Chat Area */}
            <div
                className={`flex-1 flex flex-col ${showMobileChat ? "flex" : "hidden md:flex"}`}
            >
                {selectedConv ? (
                    <>
                        {/* Chat Header */}
                        <div className="p-4 border-b border-(--border) flex items-center gap-3">
                            {/* Back button for mobile */}
                            <button
                                onClick={handleBackToList}
                                className="md:hidden p-2 -ml-2 rounded-full hover:bg-(--muted)/10 transition-colors"
                            >
                                <ChevronLeft className="w-5 h-5 text-(--muted)" />
                            </button>

                            <div className="w-10 h-10 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden border border-(--border)">
                                {selectedConv.Interlocutor?.avatar_url ? (
                                    <img
                                        src={selectedConv.Interlocutor.avatar_url}
                                        alt={selectedConv.Interlocutor.username || "User"}
                                        className="w-full h-full object-cover"
                                    />
                                ) : (
                                    <User className="w-5 h-5 text-(--muted)" />
                                )}
                            </div>
                            <div>
                                <p className="font-semibold text-foreground">
                                    {selectedConv.Interlocutor?.username || "Unknown User"}
                                </p>
                            </div>
                        </div>

                        {/* Messages */}
                        <div
                            ref={messagesContainerRef}
                            onScroll={handleMessagesScroll}
                            className="flex-1 overflow-y-auto p-4 space-y-3"
                        >
                            {isLoadingMessages ? (
                                <div className="flex items-center justify-center py-12">
                                    <Loader2 className="w-6 h-6 text-(--muted) animate-spin" />
                                </div>
                            ) : messages.length > 0 ? (
                                <>
                                    {/* Loading indicator for infinite scroll */}
                                    {isLoadingMoreMessages && (
                                        <div className="flex justify-center py-2 mb-3">
                                            <Loader2 className="w-4 h-4 text-(--muted) animate-spin" />
                                        </div>
                                    )}
                                    {messages.map((msg, index) => {
                                        const isMe = msg.sender?.id === user?.id;
                                        const isPending = msg._pending;
                                        return (
                                            <motion.div
                                                key={msg.id || index}
                                                initial={{ opacity: 0, y: 10 }}
                                                animate={{ opacity: isPending ? 0.5 : 1, y: 0 }}
                                                transition={{ duration: 0.2 }}
                                                className={`flex ${isMe ? "justify-end" : "justify-start"}`}
                                            >
                                                <div
                                                    className={`max-w-[75%] px-4 py-2.5 rounded-2xl ${isMe
                                                        ? "bg-(--accent) text-white rounded-br-md"
                                                        : "bg-(--muted)/10 text-foreground rounded-bl-md"
                                                        }`}
                                                >
                                                    <p className="text-sm whitespace-pre-wrap wrap-break-word">
                                                        {msg.message_text}
                                                    </p>
                                                    <p
                                                        className={`text-[10px] mt-1 ${isMe ? "text-white/70" : "text-(--muted)"
                                                            }`}
                                                    >
                                                        {formatMessageTime(msg.created_at)}
                                                    </p>
                                                </div>
                                            </motion.div>
                                        );
                                    })}
                                    <div ref={messagesEndRef} />
                                </>
                            ) : (
                                <div className="flex flex-col items-center justify-center py-12">
                                    <MessageCircle className="w-12 h-12 text-(--muted) mb-3 opacity-30" />
                                    <p className="text-(--muted)">No messages yet</p>
                                    <p className="text-(--muted) text-sm">Say hello!</p>
                                </div>
                            )}
                        </div>

                        {/* Message Input */}
                        <form
                            onSubmit={handleSendMessage}
                            className="p-4 border-t border-(--border)"
                        >
                            <div className="flex items-center gap-3">
                                {/* Emoji Picker */}
                                <div className="relative" ref={emojiPickerRef}>
                                    <button
                                        type="button"
                                        onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                                        className="p-3 text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-all"
                                    >
                                        <Smile className="w-5 h-5" />
                                    </button>
                                    {showEmojiPicker && (
                                        <div className="absolute bottom-14 left-0 z-50">
                                            <EmojiPicker
                                                onEmojiClick={onEmojiClick}
                                                width={320}
                                                height={400}
                                                previewConfig={{ showPreview: false }}
                                            />
                                        </div>
                                    )}
                                </div>
                                <input
                                    type="text"
                                    value={messageText}
                                    onChange={(e) => setMessageText(e.target.value)}
                                    placeholder="Type a message..."
                                    className="flex-1 px-4 py-3 border border-(--border) rounded-full text-sm bg-(--muted)/5 text-foreground placeholder-(--muted) hover:border-foreground focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all"
                                />
                                <button
                                    type="submit"
                                    disabled={!messageText.trim() || !isConnected}
                                    className="p-3 bg-(--accent) text-white rounded-full hover:bg-(--accent-hover) transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    <Send className="w-5 h-5" />
                                </button>
                            </div>
                        </form>
                    </>
                ) : (
                    /* No conversation selected */
                    <div className="flex-1 flex flex-col items-center justify-center px-4">
                        <div className="w-20 h-20 rounded-full bg-(--muted)/10 flex items-center justify-center mb-4">
                            <MessageCircle className="w-10 h-10 text-(--muted)" />
                        </div>
                        <h2 className="text-xl font-semibold text-foreground mb-2">Your Messages</h2>
                        <p className="text-(--muted) text-center max-w-sm">
                            Select a conversation from the list to start chatting
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}
