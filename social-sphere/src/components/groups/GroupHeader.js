"use client";

import { useState } from "react";
import { Users, UserPlus, Settings, LogOut, Clock, UserRoundPlus } from "lucide-react";
import Modal from "@/components/ui/Modal";
import Container from "@/components/layout/Container";
import { requestJoinGroup } from "@/actions/groups/request-join-group";
import { leaveGroup } from "@/actions/groups/leave-group";
import Tooltip from "../ui/Tooltip";
import UpdateGroupModal from "./UpdateGroupModal";
import { useRouter } from "next/navigation";

export function GroupHeader({ group }) {
    const router = useRouter();
    const [isMember, setIsMember] = useState(group.is_member);
    const [isOwner] = useState(group.is_owner);
    const [isPending, setIsPending] = useState(group.is_pending);
    const [isLoading, setIsLoading] = useState(false);
    const [showLeaveModal, setShowLeaveModal] = useState(false);
    const [showInviteModal, setShowInviteModal] = useState(false);
    const [showUpdateModal, setShowUpdateModal] = useState(false);

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

    const handleInviteMembers = () => {
        setShowInviteModal(true);
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

            {/* Invite Members Modal - Placeholder */}
            {showInviteModal && (
                <Modal
                    isOpen={showInviteModal}
                    onClose={() => setShowInviteModal(false)}
                    title="Invite Members"
                    description="Select friends to invite to this group."
                    onConfirm={() => {
                        // TODO: Implement invite logic
                        setShowInviteModal(false);
                    }}
                    confirmText="Send Invites"
                    cancelText="Cancel"
                >
                    <div className="p-4 text-center text-(--muted)">
                        <p>Invite functionality coming soon...</p>
                    </div>
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
