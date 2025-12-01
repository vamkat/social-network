"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const API_BASE =
    process.env.NEXT_PUBLIC_API_BASE || "http://localhost:4000"; // dummy

export default function GroupForm({ onCreated, onCancel, availableMembers = [] }) {
    const router = useRouter();
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [privacy, setPrivacy] = useState("public");
    const [imageUrl, setImageUrl] = useState("");
    const [invites, setInvites] = useState("");
    const [showInvitePicker, setShowInvitePicker] = useState(false);
    const [inviteQuery, setInviteQuery] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const parsedInvites = parseInvites(invites);
    const filteredMembers = (availableMembers || []).filter((member) => {
        const search = inviteQuery.toLowerCase();
        return (
            member.email?.toLowerCase().includes(search) ||
            member.name?.toLowerCase().includes(search) ||
            member.username?.toLowerCase().includes(search)
        );
    });

    function parseInvites(raw) {
        return raw
            .split(/[,\n]/)
            .map((item) => item.trim())
            .filter(Boolean);
    }

    function addInvite(email) {
        const list = new Set(parseInvites(invites));
        list.add(email);
        setInvites(Array.from(list).join(", "));
    }

    function removeInvite(email) {
        const list = parseInvites(invites).filter((item) => item !== email);
        setInvites(list.join(", "));
    }

    async function handleSubmit(e) {
        e.preventDefault();
        if (!title.trim() || !description.trim()) return;

        setIsLoading(true);
        setError("");
        setSuccess("");

        try {
            const inviteList = parseInvites(invites);

            const payload = {
                title: title.trim(),
                description: description.trim(),
                visibility: privacy,
                image: imageUrl.trim() || null,
                invites: inviteList,
            };

            const res = await fetch(`${API_BASE}/api/groups`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });

            if (!res.ok) {
                throw new Error("Failed to create group");
            }

            const data = await res.json();
            const newGroup = {
                id: data.id,
                title: data.title,
                description: data.description,
                visibility: data.visibility ?? privacy,
                image: data.image ?? payload.image,
                invites: payload.invites,
            };

            setSuccess("Group created successfully.");
            setTitle("");
            setDescription("");
            setImageUrl("");
            setInvites("");
            setInviteQuery("");
            setShowInvitePicker(false);

            if (onCreated) {
                onCreated(newGroup);
            }
        } catch (err) {
            console.error(err);
            setError("Could not create the group. Try again.");
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form
            onSubmit={handleSubmit}
            className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-6 space-y-5"
        >
            <div className="space-y-2">
                <label className="text-sm font-semibold text-(--foreground)">Group title</label>
                <input
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Name your group"
                    className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30"
                    required
                />
            </div>

            <div className="space-y-2">
                <label className="text-sm font-semibold text-(--foreground)">Description</label>
                <textarea
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows={4}
                    placeholder="Tell members what this group is about"
                    className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30 resize-none"
                    required
                />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                    <label className="text-sm font-semibold text-(--foreground)">Visibility</label>
                    <div className="flex gap-2">
                        {["public", "private"].map((option) => (
                            <button
                                key={option}
                                type="button"
                                onClick={() => setPrivacy(option)}
                                className={`px-4 py-2 rounded-full text-xs font-semibold border transition-colors ${privacy === option
                                    ? "bg-(--foreground) text-(--background) border-(--foreground)"
                                    : "border-slate-200 dark:border-slate-800 text-(--muted) hover:text-(--foreground)"
                                    }`}
                            >
                                {option.charAt(0).toUpperCase() + option.slice(1)}
                            </button>
                        ))}
                    </div>
                    <p className="text-xs text-(--muted)">
                        Public: anyone can find and join. Private: approval required.
                    </p>
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-semibold text-(--foreground)">Cover image URL (optional)</label>
                    <input
                        type="url"
                        value={imageUrl}
                        onChange={(e) => setImageUrl(e.target.value)}
                        placeholder="https://example.com/cover.jpg"
                        className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30"
                    />
                </div>
            </div>

            <div className="space-y-2">
                <label className="text-sm font-semibold text-(--foreground)">Invite members (optional)</label>
                <textarea
                    value={invites}
                    onChange={(e) => setInvites(e.target.value)}
                    rows={2}
                    placeholder="Enter emails separated by comma or new line"
                    className="w-full rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-700 dark:text-slate-100 outline-none focus:ring-2 focus:ring-purple-500/30 resize-none"
                />
                <p className="text-xs text-(--muted)">
                    We will send invitations to these emails after the group is created.
                </p>

                <div className="flex flex-wrap items-center gap-2">
                    {parsedInvites.map((email) => (
                        <span
                            key={email}
                            className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-(--muted)/10 text-xs"
                        >
                            {email}
                            <button
                                type="button"
                                onClick={() => removeInvite(email)}
                                className="text-(--muted) hover:text-(--foreground)"
                            >
                                ✕
                            </button>
                        </span>
                    ))}
                    <button
                        type="button"
                        onClick={() => setShowInvitePicker((v) => !v)}
                        className="text-xs font-semibold px-3 py-1.5 rounded-full border border-(--muted)/30 hover:border-(--foreground) transition-colors"
                    >
                        Browse members
                    </button>
                </div>

                {showInvitePicker && (
                    <div className="mt-2 rounded-xl border border-(--muted)/20 bg-(--background) shadow-sm p-3 space-y-3">
                        <input
                            type="text"
                            value={inviteQuery}
                            onChange={(e) => setInviteQuery(e.target.value)}
                            placeholder="Search by name or email"
                            className="w-full rounded-lg border border-(--muted)/20 px-3 py-2 text-sm bg-transparent focus:outline-none focus:ring-2 focus:ring-(--foreground)/20"
                        />
                        <div className="max-h-40 overflow-y-auto space-y-2 pr-1">
                            {filteredMembers.length === 0 && (
                                <p className="text-xs text-(--muted)">No matches.</p>
                            )}
                            {filteredMembers.map((member) => (
                                <button
                                    key={member.email || member.username}
                                    type="button"
                                    onClick={() => addInvite(member.email)}
                                    className="w-full flex items-center justify-between text-left px-3 py-2 rounded-lg hover:bg-(--muted)/10"
                                >
                                    <div>
                                        <div className="text-sm font-semibold">{member.name || member.username}</div>
                                        <div className="text-xs text-(--muted)">{member.email}</div>
                                    </div>
                                    <span className="text-xs text-(--muted)">Add</span>
                                </button>
                            ))}
                        </div>
                    </div>
                )}
            </div>

            {error && <p className="text-xs text-red-500 animate-fade-in">{error}</p>}
            {success && <p className="text-xs text-emerald-600 animate-fade-in">{success}</p>}

            <div className="flex items-center justify-end gap-3 text-xs">
                <button
                    type="button"
                    onClick={() => {
                        setTitle("");
                        setDescription("");
                        setImageUrl("");
                        setInvites("");
                        setSuccess("");
                        setError("");
                        if (onCancel) {
                            onCancel();
                        } else {
                            router.back();
                        }
                    }}
                    className="btn btn-secondary px-4 py-2 border border-(--muted)/100"
                >
                    Cancel
                </button>
                <button
                    type="submit"
                    disabled={isLoading || !title.trim() || !description.trim()}
                    className="btn btn-primary px-5 py-2 text-xs disabled:opacity-100 disabled:cursor-not-allowed"
                >
                    {isLoading ? "Creating..." : "Create group"}
                </button>
            </div>
        </form>
    );
}
