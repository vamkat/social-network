"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import PostCard from "@/components/ui/post-card";

export default function FeedList({ initialPosts, fetchPosts }) {
    const [posts, setPosts] = useState(initialPosts);
    const [offset, setOffset] = useState(initialPosts.length);
    const [hasMore, setHasMore] = useState(initialPosts.length >= 5);
    const [loading, setLoading] = useState(false);
    const observer = useRef();

    const lastPostElementRef = useCallback(node => {
        if (loading) return;
        if (observer.current) observer.current.disconnect();
        observer.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting && hasMore) {
                loadMorePosts();
            }
        });
        if (node) observer.current.observe(node);
    }, [loading, hasMore]);

    const loadMorePosts = async () => {
        setLoading(true);
        try {
            const limit = 5;
            const newPosts = await fetchPosts(offset, limit);

            if (newPosts.length < limit) {
                setHasMore(false);
            }

            if (newPosts.length > 0) {
                setPosts(prevPosts => [...prevPosts, ...newPosts]);
                setOffset(prevOffset => prevOffset + newPosts.length);
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Failed to fetch posts:", error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex flex-col gap-4">
            {posts.map((post, index) => {
                if (posts.length === index + 1) {
                    return (
                        <div ref={lastPostElementRef} key={`${post.ID}-${index}`}>
                            <PostCard post={post} />
                        </div>
                    );
                } else {
                    return <PostCard key={`${post.ID}-${index}`} post={post} />;
                }
            })}

            {loading && (
                <div className="flex justify-center p-4">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                </div>
            )}

            {!hasMore && posts.length > 0 && (
                <div className="text-center p-4 text-gray-500">
                    You're up to date.
                </div>
            )}
        </div>
    );
}
