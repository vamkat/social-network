"use client";

import { useState, useEffect } from "react";
import { Users, UserPlus, Settings, LogOut, Clock, UserRoundPlus, User, Check, Loader2 } from "lucide-react";
import Modal from "@/components/ui/Modal";
import Container from "@/components/layout/Container";
import { requestJoinGroup } from "@/actions/groups/request-join-group";
import { leaveGroup } from "@/actions/groups/leave-group";
import { inviteToGroup } from "@/actions/groups/invite-to-group";
import { getFollowers } from "@/actions/users/get-followers";
import Tooltip from "../ui/Tooltip";
import UpdateGroupModal from "./UpdateGroupModal";
import { useRouter } from "next/navigation";
import { useStore } from "@/store/store";

export function GroupHeader({ group }) {
    const router = useRouter();
    const user = useStore((state) => state.user);
    const [isMember, setIsMember] = useState(group.is_member);
    const [isOwner] = useState(group.is_owner);
    const [isPending, setIsPending] = useState(group.is_pending);
    const [isLoading, setIsLoading] = useState(false);
    const [showLeaveModal, setShowLeaveModal] = useState(false);
    const [showInviteModal, setShowInviteModal] = useState(false);
    const [showUpdateModal, setShowUpdateModal] = useState(false);

    // Invite modal state
    const [followers, setFollowers] = useState([]);
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [isLoadingFollowers, setIsLoadingFollowers] = useState(false);
    const [isInviting, setIsInviting] = useState(false);
    const [inviteSuccess, setInviteSuccess] = useState(false);

    const handleJoinRequest = async () => {
        if (isLoading) return;
        setIsLoading(true);

        try {
            const response = await requestJoinGroup({ groupId: group.group_id });
            if (response.success) {
                // Toggle pending state (the endpoint handles both request and cancel)
                setIsPending(!isPending);
            } else {
                console.error("Error handling join request:", response.error);
            }
        } catch (error) {
            console.error("Error handling join request:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleLeaveGroup = async () => {
        if (isLoading) return;
        setIsLoading(true);

        try {
            const response = await leaveGroup({ groupId: group.group_id });
            if (response.success) {
                setIsMember(false);
                setShowLeaveModal(false);
            } else {
                console.error("Error leaving group:", response.error);
            }
        } catch (error) {
            console.error("Error leaving group:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleInviteMembers = async () => {
        setShowInviteModal(true);
        setInviteSuccess(false);
        setSelectedUsers([]);

        if (user?.id) {
            setIsLoadingFollowers(true);
            try {
                const result = await getFollowers({ userId: user.id });
                if (Array.isArray(result)) {
                    setFollowers(result);
                } else {
                    setFollowers([]);
                }
            } catch (error) {
                console.error("Error fetching followers:", error);
                setFollowers([]);
            } finally {
                setIsLoadingFollowers(false);
            }
        }
    };

    const handleCloseInviteModal = () => {
        setShowInviteModal(false);
        setSelectedUsers([]);
        setFollowers([]);
        setInviteSuccess(false);
    };

    const toggleUserSelection = (userId) => {
        setSelectedUsers((prev) =>
            prev.includes(userId)
                ? prev.filter((id) => id !== userId)
                : [...prev, userId]
        );
    };

    const handleSendInvites = async () => {
        if (selectedUsers.length === 0 || isInviting) return;

        setIsInviting(true);
        try {
            const response = await inviteToGroup({
                groupId: group.group_id,
                invitedIds: selectedUsers,
            });

            if (response.success) {
                setInviteSuccess(true);
                setTimeout(() => {
                    handleCloseInviteModal();
                }, 1500);
            } else {
                console.error("Error inviting users:", response.error);
            }
        } catch (error) {
            console.error("Error inviting users:", error);
        } finally {
            setIsInviting(false);
        }
    };

    const handleUpdateSuccess = () => {
        // Refresh the page to get updated group data
        router.refresh();
    };

    return (
        <>
            <div className="w-full border-b border-(--border)">
                <Container>
                    <div className="py-8">
                        {/* Top Section: Group Image, Title, Actions */}
                        <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center mb-6">
                            {/* Group Image */}
                            <div className="relative">
                                <div className="w-24 h-24 sm:w-28 sm:h-28 rounded-2xl overflow-hidden bg-(--muted)/10 border-2 border-(--border) ring-4 ring-background shadow-lg">
                                    {group.group_image_url ? (
                                        <img
                                            src={group.group_image_url}
                                            alt={group.group_title}
                                            className="w-full h-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-gray-100 to-gray-200">
                                            <Users className="w-12 h-12 text-(--muted)" />
                                        </div>
                                    )}
                                </div>
                            </div>

                            {/* Title & Actions */}
                            <div className="flex-1 min-w-0 flex flex-col sm:flex-row justify-between items-start gap-4">
                                <div className="flex-1 min-w-0">
                                    <div className="flex items-center gap-3 mb-2">
                                    {isOwner && (
                                            <span className="inline-flex items-center gap-1 px-1 py-0.5 rounded-full text-[10px] bg-(--accent) text-white shadow-sm">
                                                {/* <Shield className="w-3 h-3" /> */}
                                                Owner
                                            </span>
                                        )}
                                        {!isOwner && isMember && (
                                            <span className="inline-flex items-center gap-1 px-1 py-0.5 rounded-full text-xs bg-green-500 text-white shadow-sm">
                                                Member
                                            </span>
                                        )}
                                    </div>
                                    <div className="flex items-center gap-3 mb-2">
                                        <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">
                                            {group.group_title}
                                        </h1>
                                        
                                    </div>
                                    
                                    <div className="flex items-center gap-2 text-(--muted)">
                                        <Users className="w-4 h-4" />
                                        <span className="text-base">
                                            {group.members_count || 0} {group.members_count === 1 ? "Member" : "Members"}
                                        </span>
                                    </div>
                                </div>

                                {/* Action Buttons */}
                                <div className="flex items-center gap-2 shrink-0 flex-wrap">
                                    {isOwner ? (
                                        <>
                                            {/* Invite Members */}
                                            <Tooltip content="Invite members">
                                            <button
                                                onClick={handleInviteMembers}
                                                className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium bg-(--accent) text-white hover:bg-(--accent-hover) shadow-lg shadow-(--accent)/20 transition-colors cursor-pointer"
                                            >
                                                <UserPlus className="w-4 h-4" />
                                            </button>
                                            </Tooltip>
                                            {/* Settings/Edit Group */}
                                            <Tooltip content="Settings">
                                            <button
                                                onClick={() => setShowUpdateModal(true)}
                                                className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium border border-(--border) text-foreground hover:bg-(--muted)/5 transition-colors cursor-pointer"
                                            >
                                                <Settings className="w-4 h-4" />
                                            </button>
                                            </Tooltip>
                                        </>
                                    ) : isMember ? (
                                        <>
                                            {/* Invite Members */}
                                            <Tooltip content="Invite Members">
                                            <button
                                                onClick={handleInviteMembers}
                                                className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium border border-(--accent) text-(--accent) hover:bg-(--accent)/5 transition-colors"
                                            >
                                                <UserPlus className="w-4 h-4" />
                                            </button>
                                            </Tooltip>
                                            {/* Leave Group */}
                                            <Tooltip content="Leave Group">
                                            <button
                                                onClick={() => setShowLeaveModal(true)}
                                                className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium border border-(--border) text-(--muted) hover:bg-red-500/10 hover:text-red-500 hover:border-red-500/20 transition-colors"
                                            >
                                                <LogOut className="w-4 h-4" />
                                                <span className="hidden sm:inline">Leave</span>
                                            </button>
                                            </Tooltip>
                                        </>
                                    ) : (
                                        /* Non-member: Request to Join or Cancel Request */
                                        <Tooltip content={isPending ? (
                                            "Pending Request"
                                        ) : ( "Request to Join" )} >
                                        <button
                                            onClick={handleJoinRequest}
                                            disabled={isLoading}
                                            className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all cursor-pointer ${
                                                isLoading
                                                    ? "opacity-70 cursor-wait"
                                                    : isPending
                                                    ? "bg-(--muted)/10 text-(--muted) border border-(--border) hover:bg-red-500/10 hover:text-red-500 hover:border-red-500/20"
                                                    : "bg-(--accent) text-white hover:bg-(--accent-hover) shadow-lg shadow-(--accent)/20"
                                            }`}
                                        >
                                            {isPending ? (
                                                <>
                                                    <Clock className="w-4 h-4" />
                                                </>
                                            ) : (
                                                <>
                                                    <UserRoundPlus className="w-4 h-4" />
                                                </>
                                            )}
                                        </button>
                                        </Tooltip>
                                    )}
                                </div>
                            </div>
                        </div>

                        {/* Description Section */}
                        {group.group_description && (
                            <div className="mb-6">
                                <p className="text-(--foreground)/90 leading-relaxed whitespace-pre-wrap text-[15px]">
                                    {group.group_description}
                                </p>
                            </div>
                        )}
                    </div>
                </Container>
            </div>

            {/* Leave Group Modal */}
            <Modal
                isOpen={showLeaveModal}
                onClose={() => setShowLeaveModal(false)}
                title="Leave Group?"
                description={`Are you sure you want to leave "${group.group_title}"? You will need to request to join again if you change your mind.`}
                onConfirm={handleLeaveGroup}
                confirmText="Leave Group"
                cancelText="Cancel"
                isLoading={isLoading}
                loadingText="Leaving..."
            />

            {/* Invite Members Modal */}
            {showInviteModal && (
                <Modal
                    isOpen={showInviteModal}
                    onClose={handleCloseInviteModal}
                    title="Invite Members"
                    description={inviteSuccess ? "" : "Select followers to invite to this group."}
                    onConfirm={inviteSuccess ? undefined : handleSendInvites}
                    confirmText={`Send Invites${selectedUsers.length > 0 ? ` (${selectedUsers.length})` : ""}`}
                    cancelText="Cancel"
                    isLoading={isInviting}
                    loadingText="Sending..."
                >
                    {inviteSuccess ? (
                        <div className="py-8 text-center">
                            <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-(--accent)/5 flex items-center justify-center">
                                <Check className="w-6 h-6 text-(--accent)" />
                            </div>
                            <p className="text-(--accent) font-medium">Users have been successfully invited!</p>
                        </div>
                    ) : isLoadingFollowers ? (
                        <div className="py-8 flex flex-col items-center gap-3">
                            <Loader2 className="w-6 h-6 text-(--accent) animate-spin" />
                            <p className="text-sm text-(--muted)">Loading followers...</p>
                        </div>
                    ) : followers.length === 0 ? (
                        <div className="py-8 text-center text-(--muted)">
                            <p>No followers to invite.</p>
                        </div>
                    ) : (
                        <div className="max-h-64 overflow-y-auto -mx-5 px-5">
                            <div className="space-y-1">
                                {followers.map((follower) => {
                                    const isSelected = selectedUsers.includes(follower.id);
                                    return (
                                        <button
                                            key={follower.id}
                                            type="button"
                                            onClick={() => toggleUserSelection(follower.id)}
                                            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors cursor-pointer ${
                                                isSelected
                                                    ? "bg-(--accent)/10 border border-(--accent)/30"
                                                    : "hover:bg-(--muted)/5 border border-transparent"
                                            }`}
                                        >
                                            <div className="w-10 h-10 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden shrink-0">
                                                {follower.avatar_url ? (
                                                    <img
                                                        src={follower.avatar_url}
                                                        alt={follower.username}
                                                        className="w-full h-full object-cover"
                                                    />
                                                ) : (
                                                    <User className="w-5 h-5 text-(--muted)" />
                                                )}
                                            </div>
                                            <div className="flex-1 min-w-0 text-left">
                                                <p className="text-sm font-medium text-foreground truncate">
                                                    {follower.username}
                                                </p>
                                            </div>
                                            <div className={`w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors ${
                                                isSelected
                                                    ? "bg-(--accent) border-(--accent)"
                                                    : "border-(--border)"
                                            }`}>
                                                {isSelected && <Check className="w-3 h-3 text-white" />}
                                            </div>
                                        </button>
                                    );
                                })}
                            </div>
                        </div>
                    )}
                </Modal>
            )}

            {/* Update Group Modal */}
            <UpdateGroupModal
                isOpen={showUpdateModal}
                onClose={() => setShowUpdateModal(false)}
                onSuccess={handleUpdateSuccess}
                group={group}
            />
        </>
    );
}
