"use client";

import { useState } from "react";
import { Calendar, Globe, UserPlus, UserCheck, UserMinus, Clock, Lock, Check, X } from "lucide-react";
import Modal from "@/components/ui/Modal";
import Container from "@/components/layout/Container";
import { followUser } from "@/actions/requests/follow-user";
import { unfollowUser } from "@/actions/requests/unfollow-user";
import { handleFollowRequest } from "@/actions/requests/handle-request";
import { updatePrivacyAction } from "@/actions/profile/settings";
import Tooltip from "../ui/Tooltip";

export function ProfileHeader({ user, onUnfollow=null }) {
    const [isFollowing, setIsFollowing] = useState(user.viewer_is_following);
    const [isPublic, setIsPublic] = useState(user.public);
    const [isHovering, setIsHovering] = useState(false);
    const [isPending, setIsPending] = useState(user.is_pending);
    const [isLoading, setIsLoading] = useState(false);
    const [showPrivacyModal, setShowPrivacyModal] = useState(false);

    const [userAskedToFollow, setUserAskedToFollow] = useState(user.follow_request_from_profile_owner);
    const [requestLoading, setRequestLoading] = useState(false);

    const handleAcceptRequest = async () => {
        if (requestLoading) return;
        setRequestLoading(true);
        try {
            const response = await handleFollowRequest({ requesterId: user.user_id, accept: true });
            if (response.success) {
                setUserAskedToFollow(false);
            } else {
                console.error("Error accepting follow request:", response.error);
            }
        } catch (error) {
            console.error("Unexpected error accepting follow request:", error);
        } finally {
            setRequestLoading(false);
        }
    };

    const handleDeclineRequest = async () => {
        if (requestLoading) return;
        setRequestLoading(true);
        try {
            const response = await handleFollowRequest({ requesterId: user.user_id, accept: false });
            if (response.success) {
                setUserAskedToFollow(false);
            } else {
                console.error("Error declining follow request:", response.error);
            }
        } catch (error) {
            console.error("Unexpected error declining follow request:", error);
        } finally {
            setRequestLoading(false);
        }
    };

    const handleFollow = async () => {
        if (isLoading) return;
        setIsLoading(true);

        try {
            if (isFollowing) {
                // Handle Unfollow
                const response = await unfollowUser(user.user_id);
                if (response.success) {
                    setIsFollowing(false);
                    if (onUnfollow) onUnfollow({isPublic, isFollowing: false});
                } else {
                    console.error("Error unfollowing user:", response.error);
                }
            } else if (isPending) {
                // If pending, maybe we want to cancel request?
                // For now, let's treat clicking on Pending as Unfollow/Cancel Request which usually hits the same unfollow endpoint or a specifics cancel endpoint.
                // Assuming unfollowUser handles cancelling requests too or we need a specific cancel action.
                // Assuming unfollow works for cancelling pending requests as well for now based on typical implementations, 
                // but strictly speaking we might need a cancelRequest action if the backend distinguishes.
                // Since I didn't create cancel-request, I'll try unfollowUser.
                const response = await unfollowUser(user.user_id);
                if (response.success) {
                    setIsPending(false);
                } else {
                    console.error("Error cancelling follow request:", response.error);
                }
            } else {
                // Handle Follow
                const response = await followUser(user.user_id);
                if (response.success) {
                    // Use the actual backend response to determine state
                    if (response.data.is_pending) {
                        setIsPending(true);
                        setIsFollowing(false);
                    } else if (response.data.viewer_is_following) {
                        setIsFollowing(true);
                        setIsPending(false);
                    } else {
                        // Fallback logic if needed, or error state
                        console.error("Unexpected follow state:", response.data);
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
        if (isLoading) return;
        setIsLoading(true);

        try {
            const newPrivacyState = !isPublic;
            const response = await updatePrivacyAction(newPrivacyState);

            if (response.success) {
                setIsPublic(newPrivacyState);
                setShowPrivacyModal(false);
            } else {
                console.error("Failed to update privacy settings");
            }
        } catch (error) {
            console.error("Error updating privacy:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const canViewProfile = user.own_profile || isPublic || isFollowing;

    return (
        <>
            <div className="w-full border-b border-(--border)">
                <Container>
                    <div className="py-8">
                        {/* Top Section: Avatar, Name, Actions */}
                        <div className="flex flex-col sm:flex-row gap-6 items-start sm:items-center mb-6">
                            {/* Avatar */}
                            <div className="relative">
                                <div className="w-24 h-24 sm:w-28 sm:h-28 rounded-full overflow-hidden bg-(--muted)/10 border-2 border-(--border) ring-4 ring-background shadow-lg">
                                    {user.avatar_url ? (
                                        <img
                                            src={user.avatar_url}
                                            alt={user.username}
                                            className="w-full h-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-gray-100 to-gray-200 text-4xl font-bold text-(--muted)">
                                            {user.first_name?.[0]?.toUpperCase()}
                                        </div>
                                    )}
                                </div>
                            </div>

                            {/* Name & Actions */}
                            <div className="flex-1 min-w-0 flex flex-col sm:flex-row justify-between items-start gap-4">
                                <div className="flex-1 min-w-0">
                                    <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight mb-1">
                                        {user.first_name} {user.last_name}
                                    </h1>
                                    <p className="text-(--muted) text-base">@{user.username}</p>
                                </div>


                                {/* Action Buttons */}
                                <div className="flex flex-col items-end gap-2 shrink-0">
                                    <div className="flex items-center gap-2">
                                    {user.own_profile ? (
                                        <>
                                            {/* Privacy Toggle */}
                                            <Tooltip content="Privacy">
                                            <button
                                                onClick={handlePrivacyToggle}
                                                className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium border transition-colors cursor-pointer ${
                                                    isPublic
                                                        ? "border-(--accent)/30 bg-(--accent)/5 text-(--accent) hover:bg-(--accent)/10"
                                                        : "border-(--border) text-(--muted) hover:bg-(--muted)/5"
                                                }`}
                                            >
                                                {isPublic ? (
                                                    <>
                                                        <Globe className="w-4 h-4" />
                                                        <span className="hidden sm:inline">Public</span>
                                                    </>
                                                ) : (
                                                    <>
                                                        <Lock className="w-4 h-4" />
                                                        <span className="hidden sm:inline">Private</span>
                                                    </>
                                                )}
                                            </button>
                                            </Tooltip>
                                        </>
                                    ) : (
                                        <button
                                            onClick={handleFollow}
                                            disabled={isLoading}
                                            onMouseEnter={() => setIsHovering(true)}
                                            onMouseLeave={() => setIsHovering(false)}
                                            className={`flex items-center gap-2 px-6 py-2 rounded-full text-sm font-medium transition-all cursor-pointer ${
                                                isLoading
                                                    ? "opacity-70 cursor-wait"
                                                    : isFollowing
                                                    ? "bg-(--muted)/10 text-foreground border border-(--border) hover:bg-red-500/10 hover:text-red-500 hover:border-red-500/20"
                                                    : isPending
                                                    ? "bg-(--muted)/10 text-(--muted) border border-(--border)"
                                                    : "bg-(--accent) text-white hover:bg-(--accent-hover) shadow-lg shadow-(--accent)/20"
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
                                    {/* Pending Follow Request from Profile Owner */}
                                    {userAskedToFollow && !user.own_profile && (
                                        <div className="flex flex-col items-end gap-1">
                                            <span className="text-[13px] text-(--muted)">Pending request from {user.first_name}</span>
                                            <div className="flex items-center gap-2">
                                                <Tooltip content="Decline">
                                                    <button
                                                        onClick={handleDeclineRequest}
                                                        disabled={requestLoading}
                                                        className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all cursor-pointer ${
                                                            requestLoading
                                                                ? "opacity-70 cursor-wait"
                                                                : "bg-(--muted)/10 text-(--muted) border border-(--border) hover:bg-red-500/10 hover:text-red-500 hover:border-red-500/20"
                                                        }`}
                                                    >
                                                        <X className="w-4 h-4" />
                                                    </button>
                                                </Tooltip>
                                                <Tooltip content="Accept">
                                                    <button
                                                        onClick={handleAcceptRequest}
                                                        disabled={requestLoading}
                                                        className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all cursor-pointer ${
                                                            requestLoading
                                                                ? "opacity-70 cursor-wait"
                                                                : "bg-(--accent) text-white hover:bg-(--accent-hover) shadow-lg shadow-(--accent)/20"
                                                        }`}
                                                    >
                                                        <Check className="w-4 h-4" />
                                                    </button>
                                                </Tooltip>
                                            </div>
                                        </div>
                                    )}
                                </div>

                            </div>
                        </div>

                        {/* Bio Section */}
                        {canViewProfile && user.about && (
                            <div className="mb-6">
                                <p className="text-(--foreground)/90 leading-relaxed whitespace-pre-wrap text-[15px]">
                                    {user.about}
                                </p>
                            </div>
                        )}

                        {/* Stats & Meta Row */}
                        {canViewProfile && (
                            <div className="flex flex-wrap items-center gap-x-6 gap-y-3 text-sm">
                                {/* Stats - Inline */}
                                <div className="flex items-center gap-4">
                                    <div className="flex items-center gap-1.5">
                                        <span className="font-semibold text-foreground">{user.followers_count || 0}</span>
                                        <span className="text-(--muted)">Followers</span>
                                    </div>
                                    <div className="flex items-center gap-1.5">
                                        <span className="font-semibold text-foreground">{user.following_count || 0}</span>
                                        <span className="text-(--muted)">Following</span>
                                    </div>
                                    <div className="flex items-center gap-1.5">
                                        <span className="font-semibold text-foreground">{user.groups_count || 0}</span>
                                        <span className="text-(--muted)">Groups</span>
                                    </div>
                                </div>

                                {/* Separator */}
                                <div className="hidden sm:block w-px h-4 bg-(--border)" />

                                {/* Joined Date */}
                                <div className="flex items-center gap-2 text-(--muted)">
                                    <Calendar className="w-4 h-4" />
                                    <span>Joined {new Date(user.created_at).toLocaleDateString("en-US", { month: "short", year: "numeric" })}</span>
                                </div>
                            </div>
                        )}
                    </div>
                </Container>
            </div>

            {/* Privacy Modal */}
            <Modal
                isOpen={showPrivacyModal}
                onClose={() => setShowPrivacyModal(false)}
                title={isPublic ? "Switch to Private Profile?" : "Switch to Public Profile?"}
                description={
                    isPublic
                        ? "Switching to a private profile means only your followers will be able to see your content and profile details. You will need to review and approve all new follow requests."
                        : "Switching to a public profile allows anyone to view your content and profile details. New users can follow you immediately without requiring approval."
                }
                onConfirm={confirmPrivacyToggle}
                confirmText={isPublic ? "Switch to Private" : "Switch to Public"}
                cancelText="Cancel"
                isLoading={isLoading}
                loadingText="Updating..."
            />
        </>
    );
}
