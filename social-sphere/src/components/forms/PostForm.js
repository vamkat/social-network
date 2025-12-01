"use client";

import { useState } from "react";

const API_BASE =
    process.env.NEXT_PUBLIC_API_BASE || "http://localhost:4000"; // dummy

export default function PostForm({ onPostCreated, onCancel, embed = false }) {
   // const [title, setTitle] = useState("");
    const [body, setBody] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");

    async function handleSubmit(e) {
        e.preventDefault();
        if (!body.trim()) return;

        setIsLoading(true);
        setError("");

        try {
            const res = await fetch(`${API_BASE}/api/posts`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    //post_title: title,
                    post_body: body,
                }),
            });

            if (!res.ok) {
                throw new Error("Failed to create post");
            }

            const data = await res.json();

            const newPost = {
                id: data.id,
                //post_title: data.post_title,
                post_body: data.post_body,
                post_creator: data.post_creator,
                group_id: data.group_id ?? null,
                visibility: data.visibility ?? "public",
                image_id: data.image_id ?? null,
                created_at: data.created_at,
            };

            if (onPostCreated) {
                onPostCreated(newPost);
            }

            //setTitle("");
            setBody("");
        } catch (err) {
            console.error(err);
            setError("Could not publish your post. Try again.");
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form
            onSubmit={handleSubmit}
            className={
                embed
                    ? "space-y-4"
                    : "rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4 space-y-4"
            }
        >
            <div className="flex-1 space-y-3">
                {/* <div>
                    <input
                        type="text"
                        name="post_title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30"
                        placeholder="Give your post a title..."
                        required
                    />
                </div> */}

                <div>
                    <textarea
                        name="post_body"
                        value={body}
                        onChange={(e) => setBody(e.target.value)}
                        rows={3}
                        className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30 resize-none"
                        placeholder="Share something with your circle..."
                        required
                    />
                </div>

                {error && (
                    <p className="text-xs text-red-500 animate-fade-in">{error}</p>
                )}

                <div className="flex flex-wrap items-center gap-3 text-xs text-slate-500">
                    {/* Just UI placeholders for now */}
                    <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                        📷 <span>Add photo</span>
                    </div>
                    {/* <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                        📊 <span>Create poll</span>
                    </div>
                    <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                        📅 <span>Start event</span>
                    </div> */}
                    <div className="ml-auto flex items-center gap-2">
                        <button
                            type="button"
                            onClick={() => {
                                //setTitle("");
                                setBody("");
                                if (onCancel) {
                                    onCancel();
                                } else if (onPostCreated) {
                                    onPostCreated(null);
                                }
                            }}
                            className="btn btn-secondary px-4 py-1.5 text-xs disabled:opacity-0 disabled:cursor-not-allowed"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={isLoading || !body.trim()}
                            className="btn btn-primary px-4 py-1.5 text-xs disabled:opacity-100 disabled:cursor-not-allowed"
                        >
                            {isLoading ? "Posting..." : "Post"}
                        </button>
                    </div>
                </div>
            </div>
        </form>
    );
}
