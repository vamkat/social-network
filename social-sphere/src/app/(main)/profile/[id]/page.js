"use client";

import { use, useState, useEffect, useCallback } from "react";
import ProfileHeader from "@/components/features/profile/profile-header";
import { Lock } from "lucide-react";
import { fetchUserProfile } from "@/actions/profile/profile-actions";
import { fetchUserPosts } from "@/actions/posts/posts";
import { getUserByID } from "@/mock-data/users";
import FeedList from "@/components/feed/feed-list";

export default function ProfilePage({ params }) {
    const { id } = use(params);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState(null);
    const [initialPosts, setInitialPosts] = useState([]);
    const [postsLoaded, setPostsLoaded] = useState(false);


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

    useEffect(() => {
        if (user) {
            const loadPosts = async () => {
                try {
                    const posts = await fetchUserPosts(user.ID, 0, 5);
                    setInitialPosts(posts);
                } catch (error) {
                    console.error("Failed to fetch posts:", error);
                } finally {
                    setPostsLoaded(true);
                }
            };
            loadPosts();
        }
    }, [user]);

    const fetchPosts = useCallback(async (offset, limit) => {
        if (!user) return [];
        return await fetchUserPosts(user.ID, offset, limit);
    }, [user]);

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
                    {postsLoaded ? (
                        <FeedList initialPosts={initialPosts} fetchPosts={fetchPosts} />
                    ) : (
                        <div className="flex justify-center p-4">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                        </div>
                    )}
                </div>
            )
            }
        </div >
    );
}
