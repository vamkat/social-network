"use client";

import { useEffect, useState } from "react";
import { getAllGroups } from "@/actions/groups/get-all-groups";
import { getUserGroups } from "@/actions/groups/get-user-groups";
import { useStore } from "@/store/store";
import { Users, Plus, ArrowRight } from "lucide-react";
// import Image from "next/image";
import Link from "next/link";
import { Globe } from "lucide-react";

export default function GroupsPage() {
    const { user } = useStore();
    const [userGroups, setUserGroups] = useState([]);
    const [allGroups, setAllGroups] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            try {
                const [userGroupsRes, allGroupsRes] = await Promise.all([
                    getUserGroups({ limit: 10, offset: 0 }),
                    getAllGroups({ limit: 20, offset: 0 })
                ]);

                if (userGroupsRes.success && userGroupsRes.data) {
                    // Check structure: userGroupsRes.data might be { group_arr: [...] } or just [...] depending on endpoint
                    // Based on previous files, getUserGroups returns the response directly.
                    // endpoints_group.go L56 for GetAllGroups returns `out` which is GroupArr.
                    // So data.group_arr is likely.
                    setUserGroups(userGroupsRes.data.group_arr || []);
                }

                if (allGroupsRes.success && allGroupsRes.data) {
                    setAllGroups(allGroupsRes.data.group_arr || []);
                }

            } catch (error) {
                console.error("Error fetching groups:", error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const GroupCard = ({ group }) => (
        <Link href={`/groups/${group.group_id}`} className="block group">
            <div className="bg-background border border-(--border) rounded-xl overflow-hidden hover:shadow-lg transition-all duration-300 hover:border-(--accent)/20 h-full flex flex-col">
                <div className="relative h-32 bg-(--muted)/10">
                    {group.group_image ? (
                        <div className="w-full h-full flex items-center justify-center text-(--muted)">
                            <Users className="w-8 h-8 opacity-20" />
                        </div>
                        // <Image
                        //     src={group.group_image}
                        //     alt={group.group_title}
                        //     fill
                        //     className="object-cover transition-transform duration-500 group-hover:scale-105"
                        // />
                    ) : (
                        <div className="w-full h-full flex items-center justify-center text-(--muted)">
                            <Users className="w-8 h-8 opacity-20" />
                        </div>
                    )}
                    <div className="absolute inset-0 bg-linear-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                </div>
                <div className="p-4 flex-1 flex flex-col">
                    <h3 className="font-semibold text-lg text-foreground mb-1 line-clamp-1 group-hover:text-(--accent) transition-colors">
                        {group.group_title}
                    </h3>
                    <p className="text-sm text-(--muted) line-clamp-2 mb-4 flex-1">
                        {group.group_description}
                    </p>
                    <div className="flex items-center justify-between mt-auto">
                        <span className="text-xs font-medium bg-(--muted)/10 text-(--muted) px-2 py-1 rounded-md">
                            {group.members_count || 0} Members
                        </span>
                        <span className="text-xs text-(--accent) flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all duration-300 translate-x-2 group-hover:translate-x-0">
                            View <ArrowRight className="w-3 h-3" />
                        </span>
                    </div>
                </div>
            </div>
        </Link>
    );

    return (
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-12">

            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-foreground tracking-tight">Communities</h1>
                    <p className="text-(--muted) mt-1">Discover specific interest groups and connect with like-minded people.</p>
                </div>
                {/* Placeholder for Create Group Button - logic to be added later or now if simple */}
                <button className="flex items-center gap-2 bg-(--accent) text-background px-5 py-2.5 rounded-full font-medium text-sm hover:bg-(--accent-hover) transition-all shadow-lg shadow-black/5 cursor-pointer">
                    <Plus className="w-4 h-4" />
                    Create Group
                </button>
            </div>

            {isLoading ? (
                <div className="flex items-center justify-center py-20">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-(--accent)"></div>
                </div>
            ) : (
                <>
                    {/* My Groups Section */}
                    {userGroups.length > 0 && (
                        <div className="space-y-6">
                            <div className="flex items-center gap-2">
                                <Users className="w-5 h-5 text-(--accent)" />
                                <h2 className="text-xl font-bold text-foreground">My Groups</h2>
                            </div>
                            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                                {userGroups.map((group) => (
                                    <GroupCard key={group.group_id} group={group} />
                                ))}
                            </div>
                        </div>
                    )}

                    {/* All Groups Section */}
                    <div className="space-y-6">
                        <div className="flex items-center gap-2">
                            <Globe className="w-5 h-5 text-(--accent)" />
                            <h2 className="text-xl font-bold text-foreground">Discover</h2>
                        </div>
                        {allGroups.length > 0 ? (
                            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                                {allGroups.map((group) => (
                                    <GroupCard key={group.group_id} group={group} />
                                ))}
                            </div>
                        ) : (
                            <div className="text-center py-12 bg-(--muted)/5 rounded-2xl border border-dashed border-(--border)">
                                <Users className="w-12 h-12 text-(--muted) mx-auto mb-3 opacity-20" />
                                <p className="text-(--muted)">No groups found. Be the first to create one!</p>
                            </div>
                        )}
                    </div>
                </>
            )}
        </div>
    );
}


