import { useState } from "react";
import { Calendar, Link as LinkIcon, Lock, Globe, UserPlus, UserCheck, UserMinus, MoreHorizontal } from "lucide-react";
import Image from "next/image";
import ProfileStats from "./profile-stats";
import { toggleFollowUser, togglePrivacy } from "@/actions/profile/profile-actions";
import Modal from "@/components/ui/modal";

export default function ProfileHeader({ user, isOwnProfile }) {
    const [isFollowing, setIsFollowing] = useState(user.isFollower);
    const [isPublic, setIsPublic] = useState(user.publicProf);
    const [isHovering, setIsHovering] = useState(false);
    const [isPrivacyHovering, setIsPrivacyHovering] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [showPrivacyModal, setShowPrivacyModal] = useState(false);

    const handleFollow = async () => {
        if (isLoading) return;
        setIsLoading(true);
        try {
            // Optimistic update
            setIsFollowing(!isFollowing);
            await toggleFollowUser(user.Username);
        } catch (error) {
            // Revert on error
            setIsFollowing(!isFollowing);
            console.error("Failed to toggle follow status:", error);
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
        setShowPrivacyModal(false);
        try {
            // Optimistic update
            setIsPublic(!isPublic);
            await togglePrivacy(user.Username);
        } catch (error) {
            // Revert on error
            setIsPublic(!isPublic);
            console.error("Failed to toggle privacy:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const showStats = isOwnProfile || (isPublic || isFollowing);

    return (
        <div className="relative mb-8">
            <div className="p-6 bg-(--background)">
                <div className="flex flex-col md:flex-row gap-6">
                    {/* Avatar */}
                    <div className="relative w-32 h-32 rounded-full border-4 border-(--background) bg-(--muted)/20 overflow-hidden shrink-0">
                        {user.Avatar ? (
                            <Image
                                src={user.Avatar}
                                alt={user.Username}
                                fill
                                className="object-cover"
                            />
                        ) : (
                            <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-900 text-4xl font-bold text-(--muted)">
                                {user.firstName?.[0]}
                            </div>
                        )}
                    </div>

                    {/* Main Content Column */}
                    <div className="flex-1 min-w-0 flex flex-col">
                        {/* Top Row: Info & Stats */}
                        <div className="flex flex-col md:flex-row justify-between items-start gap-4 mb-4">
                            <div>
                                <h1 className="text-2xl font-bold flex items-center gap-2">
                                    {user.firstName} {user.lastName}
                                </h1>
                                <p className="text-(--muted) font-medium">@{user.Username}</p>
                            </div>
                            {showStats && (
                                <ProfileStats stats={{
                                    followers: user.FollowersNum,
                                    following: user.FollowingNum,
                                    groups: user.GroupsNum
                                }} />
                            )}
                        </div>

                        {/* Middle: Bio */}
                        {user.AboutMe && (
                            <div className="mb-6 max-w-2xl">
                                <p className="text-(--foreground)/80 leading-relaxed whitespace-pre-wrap">
                                    {user.AboutMe}
                                </p>
                            </div>
                        )}

                        {/* Bottom Row: Meta & Actions */}
                        <div className="flex flex-col md:flex-row justify-between items-end gap-4 mt-auto">
                            {/* Meta Info */}
                            <div className="flex flex-wrap items-center gap-6 text-sm text-(--muted)">
                                <div className="flex items-center gap-2">
                                    <Calendar className="w-4 h-4" />
                                    Joined {new Date(user.CreatedAt).toLocaleDateString("en-US", { month: "long", year: "numeric" })}
                                </div>
                            </div>

                            {/* Actions */}
                            <div className="flex items-center gap-3">
                                {isOwnProfile ? (
                                    <>
                                        <button
                                            onClick={handlePrivacyToggle}
                                            onMouseEnter={() => setIsPrivacyHovering(true)}
                                            onMouseLeave={() => setIsPrivacyHovering(false)}
                                            className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-300 overflow-hidden cursor-pointer ${isPublic
                                                ? "bg-green-500/10 text-green-600 hover:bg-green-500/20"
                                                : "bg-(--muted)/10 text-(--muted) hover:bg-(--muted)/20"
                                                }`}
                                            style={{ maxWidth: isPrivacyHovering ? '200px' : '48px' }}
                                        >
                                            {isPublic ? <Globe className="w-4 h-4 shrink-0" /> : <Lock className="w-4 h-4 shrink-0" />}
                                            <span className={`whitespace-nowrap transition-opacity duration-300 ${isPrivacyHovering ? 'opacity-100' : 'opacity-0 w-0'}`}>
                                                {isPublic ? "Public Profile" : "Private Profile"}
                                            </span>
                                        </button>
                                    </>
                                ) : (
                                    <>
                                        <button
                                            onClick={handleFollow}
                                            onMouseEnter={() => setIsHovering(true)}
                                            onMouseLeave={() => setIsHovering(false)}
                                            className={`flex items-center gap-2 px-6 py-2 rounded-full text-sm font-medium transition-all cursor-pointer ${isFollowing
                                                ? "bg-(--muted)/10 text-(--foreground) hover:bg-red-500/10 hover:text-red-500 border border-transparent cursor-pointer"
                                                : "bg-(--foreground) text-(--background) hover:opacity-90 shadow-lg shadow-black/5 cursor-pointer"
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
                                            ) : (
                                                <>
                                                    <UserPlus className="w-4 h-4" />
                                                    Follow
                                                </>
                                            )}
                                        </button>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div className="section-divider" />

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
                            className="px-4 py-2 rounded-full text-sm font-medium bg-(--foreground) text-(--background) hover:opacity-90 transition-opacity cursor-pointer"
                        >
                            {isPublic ? "Switch to Private" : "Switch to Public"}
                        </button>
                    </>
                }
            />
        </div>
    );
}
