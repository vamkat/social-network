"use client";

import { useState } from "react";
import { togglePrivacy } from "@/services/profile/profile-actions";
import { Lock, Globe } from "lucide-react";

export default function PrivacyForm({ user }) {
    const [isPrivate, setIsPrivate] = useState(!user?.publicProf);
    const [isLoading, setIsLoading] = useState(false);
    const [message, setMessage] = useState(null);

    async function handleToggle() {
        setIsLoading(true);
        setMessage(null);

        try {
            // Optimistic update
            const newState = !isPrivate;
            setIsPrivate(newState);

            const result = await togglePrivacy(user.ID);

            if (result.success) {
                setMessage({ type: "success", text: "Privacy settings updated" });
            } else {
                // Revert on failure
                setIsPrivate(!newState);
                setMessage({ type: "error", text: "Failed to update privacy" });
            }
        } catch (error) {
            setIsPrivate(!isPrivate); // Revert
            setMessage({ type: "error", text: "Something went wrong" });
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div>
                <h3 className="text-lg font-semibold">Account Privacy</h3>
                <p className="text-sm text-(--muted)">Control who can see your profile and posts.</p>
            </div>

            <div className="p-6 rounded-2xl border border-(--muted)/20 bg-(--muted)/5">
                <div className="flex items-start justify-between gap-4">
                    <div className="flex gap-4">
                        <div className={`p-3 rounded-xl ${isPrivate ? 'bg-red-500/10 text-red-600' : 'bg-green-500/10 text-green-600'}`}>
                            {isPrivate ? <Lock className="w-6 h-6" /> : <Globe className="w-6 h-6" />}
                        </div>
                        <div>
                            <h4 className="font-medium text-lg">
                                {isPrivate ? "Private Account" : "Public Account"}
                            </h4>
                            <p className="text-sm text-(--muted) mt-1 max-w-md leading-relaxed">
                                {isPrivate
                                    ? "Only people you approve can see your photos and videos. Your existing followers won't be affected."
                                    : "Anyone on or off the platform can see your photos and videos."
                                }
                            </p>
                        </div>
                    </div>

                    <button
                        onClick={handleToggle}
                        disabled={isLoading}
                        className={`relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-(--foreground)/20 focus:ring-offset-2 ${isPrivate ? 'bg-(--foreground)' : 'bg-(--muted)/30'
                            }`}
                    >
                        <span
                            className={`${isPrivate ? 'translate-x-6' : 'translate-x-1'
                                } inline-block h-5 w-5 transform rounded-full bg-(--background) transition-transform duration-200`}
                        />
                    </button>
                </div>
            </div>

            {message && (
                <div className={`p-4 rounded-xl text-sm ${message.type === 'success' ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                    {message.text}
                </div>
            )}
        </div>
    );
}
