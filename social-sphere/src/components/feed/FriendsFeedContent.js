"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";
import Container from "@/components/layout/Container";
import { getFriendsPosts } from "@/actions/posts/get-friends-posts";

export default function FriendsFeedContent({ initialPosts }) {
    const [posts, setPosts] = useState(initialPosts || []);
    const [offset, setOffset] = useState(10); // Start after the initial 10 posts
    // Only hasMore if we got a full batch of 10 posts
    const [hasMore, setHasMore] = useState((initialPosts || []).length >= 10);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef(null);

    const loadMorePosts = useCallback(async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const newPosts = await getFriendsPosts({ limit: 5, offset });

            if (newPosts && newPosts.length > 0) {
                setPosts((prevPosts) => [...prevPosts, ...newPosts]);
                setOffset((prevOffset) => prevOffset + 5);

                // If we got fewer than 5 posts, we've reached the end
                if (newPosts.length < 5) {
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
    }, [offset, loading, hasMore]);

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

    return (
        <div className="w-full">
            {/* Create Post Section */}
            <Container className="pt-6 md:pt-10">
                <CreatePost />
            </Container>

            {/* Feed Header */}
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title px-4">Friends Feed</h1>
                <p className="text-center feed-subtitle px-4">What's happening in your sphere?</p>
            </div>

            <div className="section-divider mb-6" />

            {/* Posts Feed */}
            <Container className="pt-6 pb-12">
                {posts?.length > 0 ? (
                    <div className="flex flex-col">
                        {posts.map((post, index) => (
                            <PostCard key={`${post.post_id}-${index}`} post={post} />
                        ))}

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
                            Your friends haven't shared anything yet. <br></br> Why not be the first?
                        </p>
                    </div>
                )}
            </Container>
        </div>
    );
}
