"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import Container from "@/components/layout/Container";
import CreatePostGroup from "@/components/groups/CreatePostGroup";
import PostCard from "@/components/ui/PostCard";
import { getGroupPosts } from "@/actions/groups/get-group-posts";

export default function GroupPageContent({ group, firstPosts }) {
    const [activeTab, setActiveTab] = useState("posts");
    const [direction, setDirection] = useState(0);
    const [posts, setPosts] = useState(firstPosts || []);
    const [offset, setOffset] = useState(10); // Start after the initial 10 posts
    // Only hasMore if we got a full batch of 10 posts
    const [hasMore, setHasMore] = useState((firstPosts || []).length >= 10);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef(null);

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
                                                        <PostCard
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
                                                <div className="text-center py-8 text-sm text-(--muted)">
                                                    You've reached the end of the feed
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
                                {/* Events content will go here */}
                                <div className="text-center text-(--muted) py-12">
                                    <p>Events tab - Create events and events list will be added here</p>
                                </div>
                            </div>
                        )}
                    </motion.div>
                </AnimatePresence>
            </div>
        </div>
    );
}
