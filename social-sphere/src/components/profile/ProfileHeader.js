"use client";

import { useState } from "react";
import { Calendar, Link as Lock, Globe, UserPlus, UserCheck, UserMinus, Clock } from "lucide-react";
import Image from "next/image";
import ProfileStats from "./ProfileStats";
import Modal from "@/components/ui/Modal";
import { follow, unfollow } from "@/services/profile/follow";

export function ProfileHeader({ user }) {
    const [isFollowing, setIsFollowing] = useState(user.viewer_is_following);
    const [isPublic, setIsPublic] = useState(user.public);
    const [isHovering, setIsHovering] = useState(false);
    const [isPending, setIsPending] = useState(user.is_pending);
    const [isPrivacyHovering, setIsPrivacyHovering] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [showPrivacyModal, setShowPrivacyModal] = useState(false);

    const handleFollow = async () => {
        if (isLoading) return;
        setIsLoading(true);

        try {
            if (isFollowing) {
                // Handle Unfollow
                const response = await unfollow(user.user_id);
                if (response.success) {
                    setIsFollowing(false);
                    setIsPending(false);
                } else {
                    console.error("Error unfollowing user:", response.error);
                }
            } else if (isPending) {
                // If pending, maybe we want to cancel request? 
                // For now, let's treat clicking on Pending as Unfollow/Cancel Request
                const response = await unfollow(user.user_id);
                if (response.success) {
                    setIsPending(false);
                    setIsFollowing(false);
                } else {
                    console.error("Error cancelling follow request:", response.error);
                }
            } else {
                // Handle Follow
                const response = await follow(user.user_id);
                if (response.success) {
                    if (isPublic) {
                        setIsFollowing(true);
                    } else {
                        setIsPending(true);
                    }
                } else {
                    console.error("Error following user:", response.error);
                }
            }
        } catch (error) {
            console.error("Unexpected error handling follow action:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handlePrivacyToggle = () => {
        setShowPrivacyModal(true);
    };

    const confirmPrivacyToggle = async () => {
        // Placeholder for future privacy toggle implementation
        setShowPrivacyModal(false);
    };

    const canViewProfile = user.own_profile || isPublic || isFollowing;

    return (
        <>
            <div className="bg-background border border-(--border) rounded-2xl overflow-hidden mb-6 py-2">
                <div className="p-6 md:p-8">
                    <div className="flex flex-col md:flex-row gap-8">
                        {/* Avatar */}
                        <div className="relative w-32 h-32 md:w-40 md:h-40 rounded-full border-4 border-background shadow-sm shrink-0">
                            <div className="w-full h-full rounded-full overflow-hidden bg-(--muted)/10 relative">
                                {user.avatar ? (
                                    <Image
                                        src={user.avatar}
                                        alt={user.username}
                                        fill
                                        className="object-cover"
                                    />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-900 text-5xl font-bold text-(--muted)">
                                        {user.first_name?.[0]}
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Main Content Column */}
                        <div className="flex-1 min-w-0 flex flex-col pt-2">
                            {/* Top Row: Info & Stats */}
                            <div className="flex flex-col md:flex-row justify-between items-start gap-4 mb-6">
                                <div>
                                    <h1 className="text-2xl font-bold text-foreground tracking-tight mb-1">
                                        {user.first_name} {user.last_name}
                                    </h1>
                                    <p className="text-(--muted) font-medium text-base">@{user.username}</p>
                                </div>
                                {canViewProfile && (
                                    <ProfileStats stats={{
                                        followers: user.followers_count,
                                        following: user.following_count,
                                        groups: user.groups_count
                                    }} />
                                )}
                            </div>

                            {/* Middle: Bio - Only show if allowed */}
                            {canViewProfile && user.about && (
                                <div className="mb-8 max-w-2xl">
                                    <p className="text-(--foreground)/80 leading-relaxed whitespace-pre-wrap text-[15px]">
                                        {user.about}
                                    </p>
                                </div>
                            )}

                            {/* Bottom Row: Meta & Actions */}
                            <div className="flex flex-col md:flex-row justify-between items-end gap-6 mt-auto">
                                {/* Meta Info - Only show if allowed */}
                                <div className="flex flex-wrap items-center gap-6 text-sm text-(--muted)">
                                    {canViewProfile && (
                                        <div className="flex items-center gap-2">
                                            <Calendar className="w-4 h-4" />
                                            <span>Joined {new Date(user.created_at).toLocaleDateString("en-US", { month: "long", year: "numeric" })}</span>
                                        </div>
                                    )}
                                </div>

                                {/* Actions */}
                                <div className="flex items-center gap-3">
                                    {user.own_profile ? (
                                        <button
                                            onClick={handlePrivacyToggle}
                                            onMouseEnter={() => setIsPrivacyHovering(true)}
                                            onMouseLeave={() => setIsPrivacyHovering(false)}
                                            className={`flex items-center gap-2 px-4 py-2.5 rounded-full text-sm font-medium transition-all duration-300 overflow-hidden cursor-pointer border ${isPublic
                                                ? "bg-(--accent)/5 text-(--accent) border-(--accent-hover)/20 hover:bg-(--accent)/5"
                                                : "bg-(--muted)/5 text-(--muted) border-(--border) hover:bg-(--muted)/10"
                                                }`}
                                            style={{ maxWidth: isPrivacyHovering ? '200px' : '48px' }}
                                        >
                                            {isPublic ? <Globe className="w-4 h-4 shrink-0" /> : <Lock className="w-4 h-4 shrink-0" />}
                                            <span className={`whitespace-nowrap transition-opacity duration-300 ${isPrivacyHovering ? 'opacity-100' : 'opacity-0 w-0'}`}>
                                                {isPublic ? "Public Profile" : "Private Profile"}
                                            </span>
                                        </button>
                                    ) : (
                                        <button
                                            onClick={handleFollow}
                                            disabled={isLoading}
                                            onMouseEnter={() => setIsHovering(true)}
                                            onMouseLeave={() => setIsHovering(false)}
                                            className={`flex items-center gap-2 px-8 py-2.5 rounded-full text-sm font-medium transition-all cursor-pointer ${isLoading ? "opacity-70 cursor-wait" : ""
                                                } ${isFollowing
                                                    ? "bg-(--muted)/10 text-foreground hover:bg-red-500/10 hover:text-red-500 border border-transparent"
                                                    : isPending
                                                        ? "bg-(--muted)/10 text-(--muted) border border-transparent"
                                                        : "bg-foreground text-background hover:opacity-90 shadow-lg shadow-black/5"
                                                }`}
                                        >
                                            {isFollowing ? (
                                                isHovering ? (
                                                    <>
                                                        <UserMinus className="w-4 h-4" />
                                                        Unfollow
                                                    </>
                                                ) : (
                                                    <>
                                                        <UserCheck className="w-4 h-4" />
                                                        Following
                                                    </>
                                                )
                                            ) : isPending ? (
                                                <>
                                                    <Clock className="w-4 h-4" />
                                                    Pending
                                                </>
                                            ) : (
                                                <>
                                                    <UserPlus className="w-4 h-4" />
                                                    Follow
                                                </>
                                            )}
                                        </button>
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <Modal
                isOpen={showPrivacyModal}
                onClose={() => setShowPrivacyModal(false)}
                title={isPublic ? "Switch to Private Profile?" : "Switch to Public Profile?"}
                description={isPublic
                    ? "Switching to a private profile means only your followers will be able to see your content and profile details. You will need to review and approve all new follow requests."
                    : "Switching to a public profile allows anyone to view your content and profile details. New users can follow you immediately without requiring approval."
                }
                footer={
                    <>
                        <button
                            onClick={() => setShowPrivacyModal(false)}
                            className="px-4 py-2 rounded-full text-sm font-medium text-(--muted) hover:bg-(--muted)/10 transition-colors cursor-pointer"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={confirmPrivacyToggle}
                            className="px-4 py-2 rounded-full text-sm font-medium bg-(--accent) text-background hover:bg-(--accent-hover) transition-opacity cursor-pointer"
                        >
                            {isPublic ? "Switch to Private" : "Switch to Public"}
                        </button>
                    </>
                }
            />
        </>
    );
}
