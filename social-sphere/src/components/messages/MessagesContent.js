"use client";

import { useEffect, useState, useRef } from "react";
import { getConv } from "@/actions/chat/get-conv";
import { getMessages } from "@/actions/chat/get-messages";
import { sendMsg } from "@/actions/chat/send-msg";
import { useStore } from "@/store/store";
import { User, Send, MessageCircle, Loader2, ChevronLeft } from "lucide-react";
import { motion } from "motion/react";

export default function MessagesContent({ initialConversations = [] }) {
    const user = useStore((state) => state.user);
    const [conversations, setConversations] = useState(initialConversations);
    const [selectedConv, setSelectedConv] = useState(null);
    const [messages, setMessages] = useState([]);
    const [isLoadingConversations, setIsLoadingConversations] = useState(false);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [isSending, setIsSending] = useState(false);
    const [messageText, setMessageText] = useState("");
    const [showMobileChat, setShowMobileChat] = useState(false);
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);

    // Scroll to bottom of messages
    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

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
            const result = await getConv({ first: true, limit: 50 });
            if (result.success && result.data) {
                setConversations(result.data);
            }
        } catch (error) {
            console.error("Error loading conversations:", error);
        } finally {
            setIsLoadingConversations(false);
        }
    };

    // Load messages for selected conversation
    const loadMessages = async (interlocutorId) => {
        setIsLoadingMessages(true);
        try {
            const result = await getMessages({ interlocutorId, limit: 50 });
            console.log(result);
            if (result.success && result.data) {
                // Messages come in reverse order (newest first), so reverse them
                setMessages(result.data.Messages?.reverse() || []);
            }
        } catch (error) {
            console.error("Error loading messages:", error);
        } finally {
            setIsLoadingMessages(false);
        }
    };

    // Handle conversation selection
    const handleSelectConversation = (conv) => {
        setSelectedConv(conv);
        setShowMobileChat(true);
        loadMessages(conv.Interlocutor?.id);
    };

    // Handle send message
    const handleSendMessage = async (e) => {
        e.preventDefault();
        if (!messageText.trim() || !selectedConv || isSending) return;

        setIsSending(true);
        const msgToSend = messageText.trim();
        setMessageText("");

        try {
            const result = await sendMsg({
                interlocutor: selectedConv.Interlocutor.id,
                msg: msgToSend
            });

            if (result.success) {
                // Add the new message to the list
                const newMessage = {
                    id: result.id || Date.now().toString(),
                    message_text: msgToSend,
                    sender: { id: user?.id },
                    created_at: new Date().toISOString()
                };
                setMessages((prev) => [...prev, newMessage]);

                // Update conversation's last message
                setConversations((prev) =>
                    prev.map((c) =>
                        c.ConversationId === selectedConv.ConversationId
                            ? {
                                  ...c,
                                  LastMessage: { ...c.LastMessage, message_text: msgToSend, sender: { id: user?.id } },
                                  UpdatedAt: new Date().toISOString()
                              }
                            : c
                    )
                );
            }
        } catch (error) {
            console.error("Error sending message:", error);
            // Restore the message if sending failed
            setMessageText(msgToSend);
        } finally {
            setIsSending(false);
        }
    };

    // Handle back button on mobile
    const handleBackToList = () => {
        setShowMobileChat(false);
    };

    // Load conversations on mount if not provided
    useEffect(() => {
        if (initialConversations.length === 0) {
            loadConversations();
        }
    }, []);

    return (
        <div className="h-[calc(100vh-5rem)] flex bg-background">
            {/* Left Sidebar - Conversations List */}
            <div
                className={`w-full md:w-80 lg:w-96 border-r border-(--border) flex flex-col ${
                    showMobileChat ? "hidden md:flex" : "flex"
                }`}
            >
                {/* Header */}
                <div className="p-4 border-b border-(--border)">
                    <h1 className="text-xl font-bold text-foreground">Messages</h1>
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
                                    key={conv.ConversationId}
                                    onClick={() => handleSelectConversation(conv)}
                                    className={`w-full flex items-start gap-3 px-4 py-3 transition-colors cursor-pointer text-left border-b border-(--border)/50 ${
                                        selectedConv?.ConversationId === conv.ConversationId
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
                                                className={`text-sm truncate ${
                                                    hasUnreadMessages(conv)
                                                        ? "font-semibold text-foreground"
                                                        : "font-medium text-foreground"
                                                }`}
                                            >
                                                {conv.Interlocutor?.username || "Unknown User"}
                                            </p>
                                            <span className="text-xs text-(--muted) shrink-0">
                                                {formatRelativeTime(conv.UpdatedAt)}
                                            </span>
                                        </div>
                                        <p
                                            className={`text-sm mt-0.5 truncate ${
                                                hasUnreadMessages(conv) ? "text-foreground" : "text-(--muted)"
                                            }`}
                                        >
                                            {conv.LastMessage?.sender?.id === user?.id ? "You: " : ""}
                                            {truncateMessage(conv.LastMessage?.message_text)}
                                        </p>
                                    </div>
                                </button>
                            ))}
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
                            className="flex-1 overflow-y-auto p-4 space-y-3"
                        >
                            {isLoadingMessages ? (
                                <div className="flex items-center justify-center py-12">
                                    <Loader2 className="w-6 h-6 text-(--muted) animate-spin" />
                                </div>
                            ) : messages.length > 0 ? (
                                <>
                                    {messages.map((msg, index) => {
                                        const isMe = msg.sender?.id === user?.id;
                                        return (
                                            <motion.div
                                                key={msg.id || index}
                                                initial={{ opacity: 0, y: 10 }}
                                                animate={{ opacity: 1, y: 0 }}
                                                transition={{ duration: 0.2 }}
                                                className={`flex ${isMe ? "justify-end" : "justify-start"}`}
                                            >
                                                <div
                                                    className={`max-w-[75%] px-4 py-2.5 rounded-2xl ${
                                                        isMe
                                                            ? "bg-(--accent) text-white rounded-br-md"
                                                            : "bg-(--muted)/10 text-foreground rounded-bl-md"
                                                    }`}
                                                >
                                                    <p className="text-sm whitespace-pre-wrap wrap-break-word">
                                                        {msg.message_text}
                                                    </p>
                                                    <p
                                                        className={`text-[10px] mt-1 ${
                                                            isMe ? "text-white/70" : "text-(--muted)"
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
                                <input
                                    type="text"
                                    value={messageText}
                                    onChange={(e) => setMessageText(e.target.value)}
                                    placeholder="Type a message..."
                                    className="flex-1 px-4 py-3 border border-(--border) rounded-full text-sm bg-(--muted)/5 text-foreground placeholder-(--muted) hover:border-foreground focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all"
                                    disabled={isSending}
                                />
                                <button
                                    type="submit"
                                    disabled={!messageText.trim() || isSending}
                                    className="p-3 bg-(--accent) text-white rounded-full hover:bg-(--accent-hover) transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    {isSending ? (
                                        <Loader2 className="w-5 h-5 animate-spin" />
                                    ) : (
                                        <Send className="w-5 h-5" />
                                    )}
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
