"use client";

import { LogoutButton } from "@/components/LogoutButton";
import { ProfileHeader } from "@/components/profile/ProfileHeader";

export default function ProfileContent({ result }) {
    // Handle error state
    if (!result.success) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen gap-4">
                <div className="text-red-500 text-lg font-medium">
                    {result.error || "Failed to load profile"}
                </div>
                <button
                    onClick={() => window.location.reload()}
                    className="px-4 py-2 bg-(--accent) text-white rounded-lg hover:opacity-90 transition-opacity"
                >
                    Try Again
                </button>
            </div>
        );
    }

    // Handle no user
    if (!result.user) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="text-(--muted) text-lg">User not found</div>
            </div>
        );
    }

    // Render profile
    return (
        <div className="w-full py-3 animate-in fade-in duration-500">
            <div className="max-w-full mx-auto px-22 pr-22">
                <ProfileHeader user={result.user} />
                <LogoutButton />
            </div>
        </div>
    );
}