import { getAllGroups, getMyGroups } from "@/actions/groups/group-actions";
import Link from "next/link";
import GroupCard from "@/components/ui/group-card";

export default async function GroupsPage() {
    const allGroups = await getAllGroups();
    const myGroups = await getMyGroups();

    return (

        <div className="space-y-8">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Groups</h1>
                    <p className="text-muted mt-1">Discover and join communities</p>
                </div>
                <Link href="/groups/create" className="btn btn-primary">
                    Create Group
                </Link>
            </div>

            {/* Search Bar */}
            <div className="relative">
                <input
                    type="text"
                    placeholder="Search groups..."
                    className="w-full bg-muted/5 border border-muted/20 rounded-full px-6 py-3 pl-12 text-foreground placeholder:text-muted focus:outline-none focus:border-foreground transition-colors"
                />
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="absolute left-4 top-1/2 -translate-y-1/2 text-muted"
                >
                    <circle cx="11" cy="11" r="8" />
                    <path d="m21 21-4.3-4.3" />
                </svg>
            </div>

            {/* My Groups Section */}
            {myGroups.length > 0 && (
                <div>
                    <div className="section-divider" />

                    <h2 className="pt-6 text-xl font-semibold mb-6 text-center">My Groups</h2>
                    <div className="section-divider" />

                    <div className="pt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {myGroups.map((group) => (
                            <GroupCard key={group.ID} group={group} />
                        ))}
                    </div>
                </div>
            )}
            <div className="section-divider" />

            {/* All Groups Section */}
            <div>
                <h2 className="pt-6 text-xl font-semibold mb-6 text-center">Browse All Groups</h2>
                <div className="section-divider" />

                <div className="pt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {allGroups.map((group) => (
                        <GroupCard key={group.ID} group={group} />
                    ))}
                </div>
            </div>
        </div>
    );
}


