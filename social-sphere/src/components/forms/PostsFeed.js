"use client";

import { useEffect, useRef, useState } from "react";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE || "http://localhost:4000"; // dummy

export default function PostsFeed({ initialPosts = [] }) {
  const [posts, setPosts] = useState(initialPosts);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const sentinelRef = useRef(null);

  async function fetchPosts(nextPage) {
    setIsLoading(true);
    setError("");

    try {
      const res = await fetch(`${API_BASE}/api/posts?page=${nextPage}&limit=5`);

      if (!res.ok) {
        throw new Error("Failed to fetch posts");
      }

      const data = await res.json();
      const newPosts = data.posts ?? data;

      // Sort newest first by created_at
      newPosts.sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );

      setPosts((prev) => {
        const all = [...prev, ...newPosts];
        const map = new Map();
        all.forEach((p) => map.set(p.id, p));
        return Array.from(map.values());
      });

      if (newPosts.length === 0 || newPosts.length < 5) {
        setHasMore(false);
      } else {
        setPage(nextPage);
      }
    } catch (err) {
      console.error(err);
      setError("Could not load more posts.");
    } finally {
      setIsLoading(false);
    }
  }

  // Initial load
  useEffect(() => {
    if (posts.length === 0) {
      fetchPosts(1);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Infinite scroll
  useEffect(() => {
    if (!hasMore || isLoading) return;

    const sentinel = sentinelRef.current;
    if (!sentinel) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting && !isLoading && hasMore) {
          fetchPosts(page + 1);
        }
      },
      {
        root: null,
        rootMargin: "0px",
        threshold: 1.0,
      }
    );

    observer.observe(sentinel);
    return () => observer.disconnect();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, hasMore, isLoading]);

  return (
    <div className="space-y-4">
      <div className="text-xs uppercase tracking-wide text-slate-500 px-1">
        For you
      </div>

      {posts.map((post) => (
        <article
          key={post.id}
          className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4 space-y-3"
        >
          {/* Header */}
          <header className="flex gap-3 items-start">
            <div className="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-[10px]">
              {(post.post_title?.[0] || "U").toUpperCase()}
            </div>
            <div className="flex-1">
              <div className="flex items-center justify-between gap-2">
                <div>
                  <div className="text-sm font-semibold">
                    {post.post_title}
                  </div>
                  <div className="flex items-center gap-2 text-xs text-slate-500">
                    <span>
                      {new Date(post.created_at).toLocaleString("en-GB", {
                        day: "2-digit",
                        month: "short",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </span>
                    <span>·</span>
                    <span className="capitalize text-xs">
                      {post.visibility || "public"}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </header>

          {/* Body */}
          <div className="text-sm text-slate-700 dark:text-slate-200 whitespace-pre-wrap">
            {post.post_body}
          </div>

          {/* Footer actions placeholder */}
          <footer className="flex items-center gap-6 text-xs text-slate-500 pt-2">
            <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
              ❤️ <span>Like</span>
            </button>
            <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
              💬 <span>Comment</span>
            </button>
            <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
              🔁 <span>Share</span>
            </button>
          </footer>
        </article>
      ))}

      {error && (
        <div className="text-xs text-red-500 px-1 animate-fade-in">
          {error}
        </div>
      )}

      <div ref={sentinelRef} className="h-10 flex items-center justify-center">
        {isLoading && (
          <span className="text-xs text-slate-500">Loading more…</span>
        )}
        {!hasMore && posts.length > 0 && (
          <span className="text-xs text-slate-400">You are all caught up.</span>
        )}
      </div>
    </div>
  );
}
