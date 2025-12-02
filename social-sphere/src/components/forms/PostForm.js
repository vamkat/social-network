"use client";

import { useEffect, useMemo, useState } from "react";

const API_BASE =
    process.env.NEXT_PUBLIC_API_BASE || "http://localhost:4000"; // dummy

const DEFAULT_VISIBILITY_OPTIONS = [
    {
        value: "public",
        label: "Public",
        helper: "Visible in both Public and Friends feeds.",
    },
    {
        value: "friends",
        label: "Friends",
        helper: "Visible only to Friends feed.",
    },
    {
        value: "custom",
        label: "Select members",
        helper: "Choose specific members who can view this post.",
    },
    {
        value: "group",
        label: "Group",
        helper: "Visible only inside groups.",
    },
];

export default function PostForm({
    onPostCreated,
    onCancel,
    embed = false,
    defaultVisibility = "public",
    visibilityOptions,
    groupId = null,
    showVisibilityPicker = true,
    audienceMarker,
}) {
   // const [title, setTitle] = useState("");
    const [body, setBody] = useState("");
    const [visibility, setVisibility] = useState(defaultVisibility);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");

    const options = useMemo(() => {
        if (visibilityOptions?.length) return visibilityOptions;
        return DEFAULT_VISIBILITY_OPTIONS;
    }, [visibilityOptions]);

    useEffect(() => {
        // Keep the current selection in sync with the default the parent passes.
        setVisibility(defaultVisibility);
    }, [defaultVisibility]);

    const selectedOption = options.find((option) => option.value === visibility) ?? options[0];
    const hasOptions = Array.isArray(options) && options.length > 0;

    // NOTE: API call temporarily disabled; this stub only resets the form and emits a mock post.
    function handleSubmit(e) {
        e.preventDefault();
        if (!body.trim()) return;
        setError("");

        const newPost = {
            id: Date.now().toString(),
            //post_title: title,
            post_body: body,
            post_creator: "me",
            group_id: visibility === "group" ? groupId : null,
            visibility,
            image_id: null,
            created_at: new Date().toISOString(),
        };

        if (onPostCreated) {
            onPostCreated(newPost);
        }

        //setTitle("");
        setBody("");
        setVisibility(defaultVisibility);
        setIsLoading(false);
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
                    {showVisibilityPicker && hasOptions && (
                        <div className="flex items-center gap-2 text-(--muted)">
                            <span className="font-semibold text-(--foreground)">Post to</span>
                            <select
                                value={visibility}
                                onChange={(e) => setVisibility(e.target.value)}
                                title={selectedOption?.helper}
                                className="rounded-lg border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-2 py-1 text-xs text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30"
                            >
                                {options.map((option) => (
                                    <option key={option.value} value={option.value}>
                                        {option.label}
                                    </option>
                                ))}
                            </select>
                        </div>
                    )}
                    {!showVisibilityPicker && audienceMarker && (
                        <span className="px-3 py-1 rounded-full border border-(--muted)/30 text-(--muted) bg-(--muted)/10">
                            {audienceMarker}
                        </span>
                    )}
                    <div className="ml-auto flex items-center gap-2">
                        <button
                            type="button"
                            onClick={() => {
                                //setTitle("");
                                setBody("");
                                setVisibility(defaultVisibility);
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

                {visibility === "custom" && (
                    <p className="text-[11px] text-(--muted)">
                        Member picker coming soon. Your post will target selected members once configured.
                    </p>
                )}
            </div>
        </form>
    );
}
