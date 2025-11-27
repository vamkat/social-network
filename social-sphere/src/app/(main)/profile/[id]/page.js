"use client";

import { use, useState, useEffect } from "react";
import ProfileHeader from "@/components/features/profile/profile-header";
import { Lock } from "lucide-react";
import { fetchUserProfile } from "@/actions/profile/profile-actions";
import { GetPostsByUserId } from "@/mock-data/posts";
import { getUserByID } from "@/mock-data/users";
import PostCard from "@/components/ui/post-card";

export default function ProfilePage({ params }) {
    const { id } = use(params);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState(null);


    // mock data
    const currentUser = getUserByID("1");

    // Data Fetching
    useEffect(() => {
        const loadUser = async () => {
            try {
                const data = await fetchUserProfile(id);
                setUser(data);
            } catch (error) {
                console.error("Failed to fetch user:", error);
            } finally {
                setLoading(false);
            }
        };

        loadUser();
    }, [id]);

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <div className="w-8 h-8 border-4 border-(--foreground) border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    if (!user) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <div className="w-8 h-8 border-4 border-(--foreground) border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    // Check if profile is private and viewer is not following (and not owner)
    const isOwnProfile = user.ID === currentUser.ID; // Mock check
    const isPrivateView = !user.publicProf && !user.isFollower && !isOwnProfile;

    // Mock Posts Data
    const userPosts = GetPostsByUserId(user.ID);

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
                <div className="mt-8">
                    <div className="flex flex-col">
                        {userPosts.map((post, i) => (
                            <PostCard key={i} post={post} />
                        ))}
                    </div>
                </div>
            )
            }
        </div >
    );
}
