"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "motion/react";
import { Plus, Send, MessageCircle, Loader2, User, Wifi, WifiOff, Smile } from "lucide-react";
import EmojiPicker from "emoji-picker-react";
import Container from "@/components/layout/Container";
import CreatePostGroup from "@/components/groups/CreatePostGroup";
import GroupPostCard from "@/components/groups/GroupPostCard";
import CreateEventModal from "@/components/groups/CreateEventModal";
import EditEventModal from "@/components/groups/EditEventModal";
import EventCard from "@/components/groups/EventCard";
import { getGroupPosts } from "@/actions/groups/get-group-posts";
import { getGroupEvents } from "@/actions/events/get-group-events";
import { getGroupMessages } from "@/actions/chat/get-group-messages";
import { useLiveSocket, ConnectionState } from "@/context/LiveSocketContext";
import { useStore } from "@/store/store";
import Tooltip from "../ui/Tooltip";

export default function GroupPageContent({ group, firstPosts }) {
    const searchParams = useSearchParams();
    const router = useRouter();
    const user = useStore((state) => state.user);

    // Get initial tab from URL or default to "posts"
    const tabFromUrl = searchParams.get("t");
    const validTabs = ["posts", "events", "messages"];
    const initialTab = validTabs.includes(tabFromUrl) ? tabFromUrl : "posts";

    const [activeTab, setActiveTab] = useState(initialTab);
    const [direction, setDirection] = useState(0);
    const [posts, setPosts] = useState(firstPosts || []);
    const [offset, setOffset] = useState(10); // Start after the initial 10 posts
    // Only hasMore if we got a full batch of 10 posts
    const [hasMore, setHasMore] = useState((firstPosts || []).length >= 10);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef(null);

    // Events state
    const [events, setEvents] = useState([]);
    const [eventsOffset, setEventsOffset] = useState(0);
    const [hasMoreEvents, setHasMoreEvents] = useState(true);
    const [loadingEvents, setLoadingEvents] = useState(false);
    const [eventsFetched, setEventsFetched] = useState(false);
    const [isCreateEventOpen, setIsCreateEventOpen] = useState(false);
    const [isEditEventOpen, setIsEditEventOpen] = useState(false);
    const [eventToEdit, setEventToEdit] = useState(null);
    const eventsObserverTarget = useRef(null);

    // Group chat state
    const [messages, setMessages] = useState([]);
    const [messageText, setMessageText] = useState("");
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [messagesFetched, setMessagesFetched] = useState(false);
    const [hasMoreMessages, setHasMoreMessages] = useState(true);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const emojiPickerRef = useRef(null);

    // WebSocket connection
    const {
        connectionState,
        isConnected,
        subscribeToGroup,
        unsubscribeFromGroup,
        addOnGroupMessage,
        removeOnGroupMessage,
        sendGroupMessage
    } = useLiveSocket();

    const handleNewEvent = (newEvent) => {
        setEvents(prev => [newEvent, ...prev]);
    };

    const handleDeleteEvent = (eventId) => {
        setEvents(prev => prev.filter(e => e.event_id !== eventId));
    };

    const handleEditEvent = (event) => {
        setEventToEdit(event);
        setIsEditEventOpen(true);
    };

    const handleEventUpdated = (updatedEvent) => {
        setEvents(prev => prev.map(e =>
            e.event_id === updatedEvent.event_id ? updatedEvent : e
        ));
    };

    // Group chat handlers
    const scrollToBottom = useCallback(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, []);

    const formatMessageTime = (dateString) => {
        if (!dateString) return "";
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return "";
        return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    };

    // Handle incoming group messages from WebSocket
    const handleGroupMessage = useCallback((msg) => {
        // Handle both snake_case (from direct response) and PascalCase (from NATS)
        const groupId = msg.group_id || msg.GroupId;
        if (groupId !== group.group_id) return;

        const isOwnMessage = (msg.Sender?.id || msg.sender?.id) === user?.id;

        setMessages((prev) => {
            if (isOwnMessage) {
                // This is a confirmation of our sent message - replace pending with confirmed
                const messageText = msg.MessageText || msg.message_text;
                const pendingIndex = prev.findIndex(
                    (m) => m._pending && m.MessageText === messageText
                );

                if (pendingIndex !== -1) {
                    // Replace the pending message with the confirmed one
                    const updated = [...prev];
                    updated[pendingIndex] = { ...msg, _pending: false };
                    return updated;
                }
            }

            // Prevent duplicates
            const msgId = msg.Id || msg.id;
            if (prev.some((m) => (m.Id || m.id) === msgId)) return prev;
            return [...prev, msg];
        });
    }, [group.group_id, user?.id]);

    // Fetch group messages
    const fetchMessages = useCallback(async (isInitial = false) => {
        if (isLoadingMessages || (!isInitial && !hasMoreMessages)) return;

        setIsLoadingMessages(true);
        try {
            const boundary = isInitial ? null : messages[0]?.Id;
            const response = await getGroupMessages({
                groupId: group.group_id,
                limit: 50,
                boundary,
                getPrevious: true
            });

            if (response.success && response.data) {
                const newMessages = response.data.Messages || [];
                if (isInitial) {
                    setMessages(newMessages.reverse());
                } else {
                    setMessages(prev => [...newMessages.reverse(), ...prev]);
                }
                setHasMoreMessages(response.data.HaveMore ?? newMessages.length >= 50);
            } else {
                if (isInitial) setMessages([]);
                setHasMoreMessages(false);
            }
            setMessagesFetched(true);
        } catch (error) {
            console.error("Failed to fetch messages:", error);
        } finally {
            setIsLoadingMessages(false);
        }
    }, [isLoadingMessages, hasMoreMessages, messages, group.group_id]);

    // Send group message via WebSocket
    const handleSendMessage = async (e) => {
        e.preventDefault();
        if (!messageText.trim() || !isConnected) return;

        const msgToSend = messageText.trim();
        setMessageText("");

        // Generate a temporary ID to track this optimistic message
        const tempId = `temp-${Date.now()}`;

        // Optimistically add message with pending state (will show with low opacity)
        const optimisticMessage = {
            Id: tempId,
            MessageText: msgToSend,
            Sender: { id: user?.id, username: user?.username, avatar_url: user?.avatar_url },
            GroupId: group.group_id,
            CreatedAt: new Date().toISOString(),
            _pending: true, // Flag for showing low opacity until confirmed
        };
        setMessages((prev) => [...prev, optimisticMessage]);

        try {
            await sendGroupMessage(group.group_id, msgToSend);
            // Server will send the confirmed message back through WebSocket
            // The handleGroupMessage callback will update the message to remove _pending
        } catch (error) {
            console.error("Error sending message:", error);
            // Remove optimistic message and restore text on WebSocket error
            setMessages((prev) => prev.filter((m) => m.Id !== tempId));
            setMessageText(msgToSend);
        }
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

    // Keep handleGroupMessage ref updated
    const handleGroupMessageRef = useRef(handleGroupMessage);
    useEffect(() => {
        handleGroupMessageRef.current = handleGroupMessage;
    }, [handleGroupMessage]);

    // Subscribe to group WebSocket when entering the group page
    useEffect(() => {
        const groupId = group.group_id;
        const messageHandler = (msg) => handleGroupMessageRef.current(msg);

        subscribeToGroup(groupId);
        addOnGroupMessage(messageHandler);

        return () => {
            removeOnGroupMessage(messageHandler);
            unsubscribeFromGroup(groupId);
        };
    }, [group.group_id, subscribeToGroup, unsubscribeFromGroup, addOnGroupMessage, removeOnGroupMessage]);

    // Fetch messages when switching to messages tab
    useEffect(() => {
        if (activeTab === "messages" && !messagesFetched) {
            fetchMessages(true);
        }
    }, [activeTab, messagesFetched, fetchMessages]);

    // Scroll to bottom when new messages arrive
    useEffect(() => {
        scrollToBottom();
    }, [messages, scrollToBottom]);

    // Fetch events when switching to events tab
    const fetchEvents = useCallback(async (isInitial = false) => {
        if (loadingEvents || (!isInitial && !hasMoreEvents)) return;

        setLoadingEvents(true);
        try {
            const currentOffset = isInitial ? 0 : eventsOffset;
            const response = await getGroupEvents({
                groupId: group.group_id,
                limit: 10,
                offset: currentOffset
            });

            if (response.success && response.data?.length > 0) {
                if (isInitial) {
                    setEvents(response.data);
                    setEventsOffset(10);
                } else {
                    setEvents(prev => [...prev, ...response.data]);
                    setEventsOffset(prev => prev + 10);
                }

                if (response.data.length < 10) {
                    setHasMoreEvents(false);
                }
            } else {
                if (isInitial) {
                    setEvents([]);
                }
                setHasMoreEvents(false);
            }
            setEventsFetched(true);
        } catch (error) {
            console.error("Failed to fetch events:", error);
        } finally {
            setLoadingEvents(false);
        }
    }, [eventsOffset, loadingEvents, hasMoreEvents, group.group_id]);

    // Fetch events when tab changes to events
    useEffect(() => {
        if (activeTab === "events" && !eventsFetched) {
            fetchEvents(true);
        }
    }, [activeTab, eventsFetched]);

    // Infinite scroll for events
    useEffect(() => {
        if (activeTab !== "events") return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMoreEvents && !loadingEvents && eventsFetched) {
                    fetchEvents(false);
                }
            },
            { threshold: 0.1 }
        );

        if (eventsObserverTarget.current) {
            observer.observe(eventsObserverTarget.current);
        }

        return () => {
            if (eventsObserverTarget.current) {
                observer.unobserve(eventsObserverTarget.current);
            }
        };
    }, [activeTab, fetchEvents, hasMoreEvents, loadingEvents, eventsFetched]);

    const handleNewPost = (newPost) => {
        setPosts(prev => [newPost, ...prev]);
    }

    const loadMorePosts = useCallback(async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const response = await getGroupPosts({ groupId: group.group_id, limit: 5, offset });

            if (response.success && response.data?.length > 0) {
                setPosts((prevPosts) => [...prevPosts, ...response.data]);
                setOffset((prevOffset) => prevOffset + 5);

                // If we got fewer than 5 posts, we've reached the end
                if (response.data.length < 5) {
                    setHasMore(false);
                }
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Failed to load more posts:", error);
        } finally {
            setLoading(false);
        }
    }, [offset, loading, hasMore, group.group_id]);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !loading) {
                    loadMorePosts();
                }
            },
            { threshold: 0.1 }
        );

        if (observerTarget.current) {
            observer.observe(observerTarget.current);
        }

        return () => {
            if (observerTarget.current) {
                observer.unobserve(observerTarget.current);
            }
        };
    }, [loadMorePosts, hasMore, loading]);

    const tabs = [
        { id: "posts", label: "Posts" },
        { id: "events", label: "Events" },
        { id: "messages", label: "Messages"},
    ];

    const handleTabChange = (tabId) => {
        const currentIndex = tabs.findIndex((t) => t.id === activeTab);
        const newIndex = tabs.findIndex((t) => t.id === tabId);
        setDirection(newIndex > currentIndex ? 1 : -1);
        setActiveTab(tabId);

        // Update URL without full page reload
        const params = new URLSearchParams(searchParams.toString());
        params.set("t", tabId);
        router.replace(`?${params.toString()}`, { scroll: false });
    };

    const slideVariants = {
        enter: (direction) => ({
            x: direction > 0 ? 300 : -300,
            opacity: 0,
        }),
        center: {
            x: 0,
            opacity: 1,
        },
        exit: (direction) => ({
            x: direction > 0 ? -300 : 300,
            opacity: 0,
        }),
    };

    return (
        <div className="w-full">
            {/* Tabs Navigation */}
            <div className="border-b border-(--border) bg-background sticky top-0 z-10">
                <Container>
                    <div className="flex gap-1">
                        {tabs.map((tab) => {
                            const isActive = activeTab === tab.id;
                            return (
                                <button
                                    key={tab.id}
                                    onClick={() => handleTabChange(tab.id)}
                                    className={`flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors relative cursor-pointer ${isActive
                                        ? "text-(--accent)"
                                        : "text-(--muted) hover:text-foreground"
                                        }`}
                                >
                                    <span>{tab.label}</span>
                                    {isActive && (
                                        <motion.span
                                            layoutId="groupTabIndicator"
                                            className="absolute bottom-0 left-0 right-0 h-0.5 bg-(--accent)"
                                            transition={{ type: "spring", stiffness: 500, damping: 30 }}
                                        />
                                    )}
                                </button>
                            );
                        })}
                    </div>
                </Container>
            </div>

            {/* Tab Content */}
            <div className="overflow-hidden">
                <AnimatePresence mode="wait" custom={direction}>
                    <motion.div
                        key={activeTab}
                        custom={direction}
                        variants={slideVariants}
                        initial="enter"
                        animate="center"
                        exit="exit"
                        transition={{ type: "spring", stiffness: 3000, damping: 300 }}
                    >
                        {activeTab === "posts" && (
                            <div>
                                <Container className="pt-6 md:pt-10 mb-6">
                                    <CreatePostGroup onPostCreated={handleNewPost} groupId={group.group_id} />
                                </Container>

                                <div className="section-divider mb-6" />

                                {/* Posts Feed */}
                                <Container className="pt-6 pb-12">
                                    {posts?.length > 0 ? (
                                        <div className="flex flex-col">
                                            <AnimatePresence mode="popLayout">
                                                {posts.map((post, index) => (
                                                    <motion.div
                                                        key={post.post_id + index}
                                                        initial={{ opacity: 0, scale: 0.8 }}
                                                        animate={{ opacity: 1, scale: 1 }}
                                                        exit={{ opacity: 0, scale: 0.95 }}
                                                        transition={{
                                                            duration: 0.3,
                                                            ease: "easeOut"
                                                        }}
                                                        layout
                                                    >
                                                        <GroupPostCard
                                                            post={post}
                                                            onDelete={(postId) => setPosts(prev => prev.filter(p => p.post_id !== postId))}
                                                        />
                                                    </motion.div>
                                                ))}
                                            </AnimatePresence>

                                            {/* Loading indicator */}
                                            {hasMore && (
                                                <div ref={observerTarget} className="flex justify-center py-8">
                                                    {loading && (
                                                        <div className="text-sm text-(--muted)">Loading more posts...</div>
                                                    )}
                                                </div>
                                            )}

                                            {/* End of feed message */}
                                            {!hasMore && posts.length > 0 && (
                                                <div className="text-center py-8 text-xl font-bold text-(--muted)">
                                                    .
                                                </div>
                                            )}
                                        </div>
                                    ) : (
                                        <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                                            <p className="text-muted text-center max-w-md px-4">
                                                Be the first ever to share something!
                                            </p>
                                        </div>
                                    )}
                                </Container>
                            </div>
                        )}

                        {activeTab === "events" && (
                            <div>
                                <Container className="pt-6 md:pt-10 mb-6">
                                    <div className="relative flex items-center">
                                        {/* Centered Title */}
                                        <div className="mx-auto text-center mt-8 mb-6">
                                            <h1 className="feed-title px-4">Events</h1>
                                            <p className="feed-subtitle px-4">
                                                What's happening in your group?
                                            </p>
                                        </div>

                                        {/* Create Event Button (right) */}
                                        <Tooltip content="Create Event">
                                            <button
                                                onClick={() => setIsCreateEventOpen(true)}
                                                className="flex items-center gap-2 bg-(--accent) text-white px-4 py-2 rounded-full text-sm font-medium hover:bg-(--accent-hover) transition-all shadow-lg shadow-black/5 cursor-pointer"
                                            >
                                                <Plus className="w-5 h-5" />
                                                {/* <span>Create Event</span> */}
                                            </button>
                                        </Tooltip>
                                    </div>
                                </Container>


                                <div className="section-divider my-6" />

                                {/* Events List */}
                                <Container className="pb-12 mt-6">
                                    {loadingEvents && events.length === 0 ? (
                                        <div className="flex flex-col items-center justify-center py-20">
                                            <div className="w-8 h-8 border-2 border-(--accent) border-t-transparent rounded-full animate-spin" />
                                            <p className="text-sm text-(--muted) mt-4">Loading events...</p>
                                        </div>
                                    ) : events.length > 0 ? (
                                        <div className="flex flex-col gap-6">
                                            <AnimatePresence mode="popLayout">
                                                {events.map((event) => (
                                                    <motion.div
                                                        key={event.event_id}
                                                        initial={{ opacity: 0, scale: 0.95 }}
                                                        animate={{ opacity: 1, scale: 1 }}
                                                        exit={{ opacity: 0, scale: 0.95 }}
                                                        transition={{ duration: 0.2 }}
                                                        layout
                                                    >
                                                        <EventCard
                                                            event={event}
                                                            onDelete={handleDeleteEvent}
                                                            onEdit={handleEditEvent}
                                                        />
                                                    </motion.div>
                                                ))}
                                            </AnimatePresence>

                                            {/* Loading indicator for infinite scroll */}
                                            {hasMoreEvents && (
                                                <div ref={eventsObserverTarget} className="flex justify-center py-8">
                                                    {loadingEvents && (
                                                        <div className="text-sm text-(--muted)">Loading more events...</div>
                                                    )}
                                                </div>
                                            )}

                                            {/* End of feed message */}
                                            {!hasMoreEvents && events.length > 0 && (
                                                <div className="text-center py-8 text-xl font-bold text-(--muted)">
                                                    .
                                                </div>
                                            )}
                                        </div>
                                    ) : (
                                        <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                                            <p className="text-muted text-center max-w-md px-4">
                                                No events yet. Create the first event for this group!
                                            </p>
                                        </div>
                                    )}
                                </Container>

                                {/* Create Event Modal */}
                                <CreateEventModal
                                    isOpen={isCreateEventOpen}
                                    onClose={() => setIsCreateEventOpen(false)}
                                    onSuccess={handleNewEvent}
                                    groupId={group.group_id}
                                />

                                {/* Edit Event Modal */}
                                <EditEventModal
                                    isOpen={isEditEventOpen}
                                    onClose={() => {
                                        setIsEditEventOpen(false);
                                        setEventToEdit(null);
                                    }}
                                    onSuccess={handleEventUpdated}
                                    event={eventToEdit}
                                />
                            </div>
                        )}

                        {activeTab === "messages" && (
                            <div className="h-[calc(100vh-10rem)] flex flex-col">
                                {/* Chat Header */}
                                <Container className="py-4 border-b border-(--border)">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-3">
                                            <MessageCircle className="w-5 h-5 text-(--accent)" />
                                            <h2 className="font-semibold text-foreground">Group Chat</h2>
                                        </div>
                                        {/* Connection Status */}
                                        <div
                                            className={`flex items-center gap-1.5 px-2 py-1 rounded-full text-xs ${
                                                isConnected
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
                                            <span>
                                                {isConnected
                                                    ? "Live"
                                                    : connectionState === ConnectionState.RECONNECTING
                                                    ? "Reconnecting"
                                                    : "Offline"}
                                            </span>
                                        </div>
                                    </div>
                                </Container>

                                {/* Messages Area */}
                                <div
                                    ref={messagesContainerRef}
                                    className="flex-1 overflow-y-auto"
                                >
                                    <Container className="py-4 space-y-3">
                                        {isLoadingMessages && messages.length === 0 ? (
                                            <div className="flex items-center justify-center py-12">
                                                <Loader2 className="w-6 h-6 text-(--muted) animate-spin" />
                                            </div>
                                        ) : messages.length > 0 ? (
                                            <>
                                                {messages.map((msg, index) => {
                                                    const isMe = msg.Sender?.id === user?.id;
                                                    const isPending = msg._pending;
                                                    return (
                                                        <motion.div
                                                            key={msg.Id || index}
                                                            initial={{ opacity: 0, y: 10 }}
                                                            animate={{ opacity: isPending ? 0.5 : 1, y: 0 }}
                                                            transition={{ duration: 0.2 }}
                                                            className={`flex ${isMe ? "justify-end" : "justify-start"}`}
                                                        >
                                                            <div className={`flex gap-2 max-w-[75%] ${isMe ? "flex-row-reverse" : ""}`}>
                                                                {/* Avatar */}
                                                                {!isMe && (
                                                                    <div className="w-8 h-8 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden border border-(--border) shrink-0">
                                                                        {msg.Sender?.avatar_url ? (
                                                                            <img
                                                                                src={msg.Sender.avatar_url}
                                                                                alt={msg.Sender.username || "User"}
                                                                                className="w-full h-full object-cover"
                                                                            />
                                                                        ) : (
                                                                            <User className="w-4 h-4 text-(--muted)" />
                                                                        )}
                                                                    </div>
                                                                )}
                                                                {/* Message Bubble */}
                                                                <div
                                                                    className={`px-4 py-2.5 rounded-2xl ${
                                                                        isMe
                                                                            ? "bg-(--accent) text-white rounded-br-md"
                                                                            : "bg-(--muted)/10 text-foreground rounded-bl-md"
                                                                    }`}
                                                                >
                                                                    {!isMe && (
                                                                        <p className="text-xs font-medium text-(--accent) mb-1">
                                                                            {msg.Sender?.username || "Unknown"}
                                                                        </p>
                                                                    )}
                                                                    <p className="text-sm whitespace-pre-wrap wrap-break-word">
                                                                        {msg.MessageText}
                                                                    </p>
                                                                    <p
                                                                        className={`text-[10px] mt-1 ${
                                                                            isMe ? "text-white/70" : "text-(--muted)"
                                                                        }`}
                                                                    >
                                                                        {formatMessageTime(msg.CreatedAt)}
                                                                    </p>
                                                                </div>
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
                                                <p className="text-(--muted) text-sm">Start the conversation!</p>
                                            </div>
                                        )}
                                    </Container>
                                </div>

                                {/* Message Input */}
                                <div className="border-t border-(--border) bg-background">
                                    <Container className="py-4">
                                        <form onSubmit={handleSendMessage} className="flex items-center gap-3">
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
                                                className="p-3 bg-(--accent) text-white rounded-full hover:bg-(--accent-hover) transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                                            >
                                                <Send className="w-5 h-5" />
                                            </button>
                                        </form>
                                    </Container>
                                </div>
                            </div>
                        )}
                    </motion.div>
                </AnimatePresence>
            </div>
        </div>
    );
}
