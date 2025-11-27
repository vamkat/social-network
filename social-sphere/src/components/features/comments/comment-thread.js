"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { countCommentsForPost, fetchPaginatedComments } from "@/mock-data/comments";

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export default function CommentThread({ postId, batchSize = 3, totalCount, skipCommentId }) {
    const [comments, setComments] = useState([]);
    const [cursor, setCursor] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const [prefetching, setPrefetching] = useState(false);
    const [prefetchedPage, setPrefetchedPage] = useState(null);
    const [total, setTotal] = useState(totalCount ?? 0);

    const sentinelRef = useRef(null);

    const displayedCount = comments.length;
    const resolvedTotal = useMemo(() => {
        const base = total || countCommentsForPost(postId);
        return skipCommentId ? Math.max(0, base - 1) : base;
    }, [postId, skipCommentId, total]);
    const hasVisibleComments = comments.length > 0;

    // Reset when post changes
    useEffect(() => {
        let mounted = true;
        setComments([]);
        setCursor(0);
        setHasMore(true);
        setPrefetchedPage(null);
        setTotal(totalCount ?? 0);

        const bootstrap = async () => {
            setLoading(true);
            const page = await fetchPaginatedComments(postId, 0, batchSize);
            if (!mounted) return;
            setComments(skipCommentId ? page.comments.filter((c) => c.ID !== skipCommentId) : page.comments);
            setCursor(page.nextCursor);
            setHasMore(page.hasMore);
            setTotal(page.total);
            setLoading(false);
        };

        bootstrap();

        return () => {
            mounted = false;
        };
    }, [batchSize, postId, totalCount]);

    const appendUnique = useCallback((incoming) => {
        setComments((prev) => {
            const seen = new Set(prev.map((c) => c.ID));
            const filtered = incoming.filter((c) => !seen.has(c.ID));
            return [...prev, ...filtered];
        });
    }, []);

    const loadFromPrefetch = useCallback(() => {
        if (!prefetchedPage) return false;
        if (!prefetchedPage.comments.length) {
            setHasMore(false);
            setPrefetchedPage(null);
            return true;
        }
        appendUnique(prefetchedPage.comments);
        setCursor(prefetchedPage.nextCursor);
        setHasMore(prefetchedPage.hasMore);
        setPrefetchedPage(null);
        return true;
    }, [appendUnique, prefetchedPage]);

    const loadNext = useCallback(async () => {
        if (loading) return;
        if (loadFromPrefetch()) return;
        if (!hasMore) return;

        setLoading(true);
        const page = await fetchPaginatedComments(postId, cursor, batchSize);
        if (!page.comments.length) {
            setHasMore(false);
            setLoading(false);
            return;
        }
        const filtered = skipCommentId ? page.comments.filter((c) => c.ID !== skipCommentId) : page.comments;
        if (filtered.length === 0) {
            setCursor(page.nextCursor);
            setHasMore(page.hasMore);
            setLoading(false);
            return;
        }
        appendUnique(filtered);
        setCursor(page.nextCursor);
        setHasMore(page.hasMore);
        setTotal(page.total);
        setLoading(false);
    }, [appendUnique, batchSize, cursor, hasMore, loadFromPrefetch, loading, postId, skipCommentId]);

    const prefetchNext = useCallback(async () => {
        if (!hasMore || prefetching || prefetchedPage || loading) return;
        setPrefetching(true);
        const page = await fetchPaginatedComments(postId, cursor, batchSize);
        await delay(50); // tiny buffer to avoid layout thrash
        setPrefetchedPage(page);
        setPrefetching(false);
    }, [batchSize, cursor, hasMore, loading, postId, prefetchedPage, prefetching]);

    // Auto-prefetch the next chunk once we render the current one
    useEffect(() => {
        prefetchNext();
    }, [prefetchNext]);

    // Infinite scroll trigger
    useEffect(() => {
        if (!hasMore) return;
        const sentinel = sentinelRef.current;
        if (!sentinel) return;

        const observer = new IntersectionObserver(
            (entries) => {
                const entry = entries[0];
                if (entry.isIntersecting) {
                    loadNext();
                }
            },
            { rootMargin: "160px 0px" }
        );

        observer.observe(sentinel);
        return () => observer.disconnect();
    }, [hasMore, loadNext]);

    return (
        <div className="comment-thread">
            <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-semibold text-(--foreground)">Comments</span>
                <span className="text-xs text-(--muted)">
                    Showing {Math.min(displayedCount, resolvedTotal)} of {resolvedTotal}
                </span>
            </div>

            {!hasVisibleComments && !loading && (
                <p className="text-sm text-(--muted)">No comments yet.</p>
            )}

            <div className="space-y-3">
                {comments.map((comment) => (
                    <div key={comment.ID} className="comment-row">
                        <div className="comment-avatar">
                            <img src={comment.BasicUserInfo.Avatar} alt={comment.BasicUserInfo.Username} />
                        </div>
                        <div className="comment-body">
                            <div className="comment-meta">
                                <span className="font-semibold text-sm">@{comment.BasicUserInfo.Username}</span>
                                <span className="text-xs text-(--muted)">{comment.CreatedAt}</span>
                            </div>
                            <p className="comment-content">{comment.Content}</p>
                        </div>
                    </div>
                ))}

                {(loading || prefetching) && (
                    <div className="comment-loader">
                        <span className="loader-dot" />
                        <span>{loading ? "Loading comments..." : "Buffering next 3..."}</span>
                    </div>
                )}

                {hasMore && (
                    <button
                        className="comment-load-more"
                        onClick={loadNext}
                        ref={sentinelRef}
                        aria-label="Load more comments"
                    >
                        Load more
                    </button>
                )}
            </div>
        </div>
    );
}
