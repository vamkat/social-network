"use client";

import { use, useState, useEffect } from "react";
import ProfileHeader from "@/components/features/profile/profile-header";
import { Lock } from "lucide-react";
import { fetchUserProfile } from "@/actions/profile-actions";

export default function ProfilePage({ params }) {
    const { username } = use(params);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState(null);

    // Data Fetching
    useEffect(() => {
        const loadUser = async () => {
            try {
                const data = await fetchUserProfile(username);
                setUser(data);
            } catch (error) {
                console.error("Failed to fetch user:", error);
            } finally {
                setLoading(false);
            }
        };

        loadUser();
    }, [username]);

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <div className="w-8 h-8 border-4 border-(--foreground) border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    // Check if profile is private and viewer is not following (and not owner)
    const isOwnProfile = username === "ychaniot"; // Mock check
    const isPrivateView = !user.publicProf && !user.isFollower && !isOwnProfile;

    return (
        <div className="animate-in fade-in duration-500">
            <ProfileHeader user={user} isOwnProfile={isOwnProfile} />

            {isPrivateView ? (
                <div className="flex flex-col items-center justify-center py-24 text-center bg-(--muted)/5 rounded-2xl border border-(--muted)/10">
                    <div className="w-16 h-16 rounded-full bg-(--muted)/10 flex items-center justify-center mb-4">
                        <Lock className="w-8 h-8 text-(--muted)" />
                    </div>
                    <h2 className="text-xl font-bold mb-2">This profile is private</h2>
                    <p className="text-(--muted) max-w-md">
                        Follow this account to see their photos and videos.
                    </p>
                </div>
            ) : (
                <>
                    <div className="space-y-6">
                        <h3 className="text-lg font-bold">Posts</h3>
                        {/* Placeholder for Posts Feed */}
                        {[...Array(3)].map((_, i) => (
                            <div key={i} className="p-6 rounded-2xl bg-(--muted)/5 border border-(--muted)/10 space-y-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-full bg-(--muted)/20" />
                                    <div>
                                        <div className="font-medium">{user.firstName} {user.lastName}</div>
                                        <div className="text-sm text-(--muted)">2 hours ago</div>
                                    </div>
                                </div>
                                <div className="h-32 rounded-xl bg-(--muted)/10" />
                            </div>
                        ))}
                    </div>
                </>
            )}
        </div>
    );
}
