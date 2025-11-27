"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import PostCard from "@/components/ui/post-card";

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export default function PostFeed({ posts = [], pageSize = 4 }) {
    const [visiblePosts, setVisiblePosts] = useState(() => posts.slice(0, pageSize));
    const [hasMore, setHasMore] = useState(posts.length > pageSize);
    const [isFetching, setIsFetching] = useState(false);
    const sentinelRef = useRef(null);

    // Reset when posts change
    useEffect(() => {
        setVisiblePosts(posts.slice(0, pageSize));
        setHasMore(posts.length > pageSize);
    }, [pageSize, posts]);

    const loadMore = useCallback(async () => {
        if (isFetching || !hasMore) return;
        setIsFetching(true);
        await delay(500);

        setVisiblePosts((prev) => {
            const nextSlice = posts.slice(prev.length, prev.length + pageSize);
            const combined = [...prev, ...nextSlice];
            setHasMore(combined.length < posts.length);
            return combined;
        });

        setIsFetching(false);
    }, [hasMore, isFetching, pageSize, posts]);

    // Infinite scroll
    useEffect(() => {
        if (!hasMore) return;
        const sentinel = sentinelRef.current;
        if (!sentinel) return;

        const observer = new IntersectionObserver(
            (entries) => {
                const entry = entries[0];
                if (entry.isIntersecting) {
                    loadMore();
                }
            },
            { rootMargin: "240px 0px" }
        );

        observer.observe(sentinel);
        return () => observer.disconnect();
    }, [hasMore, loadMore]);

    if (!posts.length) {
        return (
            <div className="rounded-xl border border-(--muted)/10 p-6 text-(--muted)">
                No posts to show yet.
            </div>
        );
    }

    return (
        <div className="post-feed">
            <div className="flex flex-col">
                {visiblePosts.map((post) => (
                    <PostCard key={post.ID} post={post} />
                ))}
            </div>

            <div className="flex flex-col items-center gap-3 py-6">
                {isFetching && (
                    <span className="text-xs text-(--muted)">Loading more posts...</span>
                )}

                {hasMore && (
                    <button
                        className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-(--foreground) text-(--background) text-sm font-semibold hover:opacity-90 transition-opacity"
                        onClick={loadMore}
                        ref={sentinelRef}
                    >
                        Load more posts
                    </button>
                )}

                {!hasMore && (
                    <span className="text-xs text-(--muted)">You have reached the end.</span>
                )}
            </div>
        </div>
    );
}
