"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "motion/react";
import { Plus } from "lucide-react";
import Container from "@/components/layout/Container";
import CreatePostGroup from "@/components/groups/CreatePostGroup";
import GroupPostCard from "@/components/groups/GroupPostCard";
import CreateEventModal from "@/components/groups/CreateEventModal";
import EditEventModal from "@/components/groups/EditEventModal";
import EventCard from "@/components/groups/EventCard";
import { getGroupPosts } from "@/actions/groups/get-group-posts";
import { getGroupEvents } from "@/actions/events/get-group-events";
import Tooltip from "../ui/Tooltip";

export default function GroupPageContent({ group, firstPosts }) {
    const searchParams = useSearchParams();
    const router = useRouter();

    // Get initial tab from URL or default to "posts"
    const tabFromUrl = searchParams.get("t");
    const validTabs = ["posts", "events"];
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
                    </motion.div>
                </AnimatePresence>
            </div>
        </div>
    );
}
