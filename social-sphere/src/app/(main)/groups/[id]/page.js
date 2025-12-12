import { getGroupById, getGroupPosts, getGroupMembers, getGroupEvents } from "@/services/groups/group-actions";
import PostCard from "@/components/ui/post-card";

export default async function GroupDetailPage({ params }) {
    const { id } = await params;
    const group = await getGroupById(id);

    const isMember = group.IsMember || group.IsOwner;

    // Fetch data only if member
    const posts = isMember ? await getGroupPosts(id) : [];
    const members = isMember ? await getGroupMembers(id) : [];
    const events = isMember ? await getGroupEvents(id) : [];

    return (
        <div className="space-y-8">
            {/* Group Header */}
            <div className="bg-(--background) border-b border-(--muted)/10 pb-8">
                <div className="aspect-3/1 w-full bg-(--muted)/10 rounded-xl relative overflow-hidden mb-6">
                    {group.Image ? (
                        <img
                            src={group.Image}
                            alt={group.Title}
                            className="w-full h-full object-cover"
                        />
                    ) : (
                        <div className="absolute inset-0 flex items-center justify-center text-(--muted)/20 text-6xl font-bold bg-linear-to-br from-(--muted)/5 to-(--muted)/20">
                            {group.Title.charAt(0)}
                        </div>
                    )}
                </div>

                <div className="flex items-start justify-between">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight mb-2">{group.Title}</h1>
                        <p className="text-(--muted) text-lg max-w-2xl leading-relaxed mb-4">
                            {group.Description}
                        </p>
                        <div className="flex items-center gap-4 text-sm text-(--muted)">
                            <div className="flex items-center gap-1.5">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                    <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                                    <circle cx="9" cy="7" r="4" />
                                    <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
                                    <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                                </svg>
                                {group.MembersNum} members
                            </div>
                            <span>•</span>
                            <div className="flex items-center gap-1.5">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                    <rect width="18" height="18" x="3" y="4" rx="2" ry="2" />
                                    <line x1="16" x2="16" y1="2" y2="6" />
                                    <line x1="8" x2="8" y1="2" y2="6" />
                                    <line x1="3" x2="21" y1="10" y2="10" />
                                </svg>
                                Created {group.CreatedAt}
                            </div>
                        </div>
                    </div>

                    {!isMember && (
                        <button className="btn btn-primary">
                            Join Group
                        </button>
                    )}
                    {isMember && (
                        <button className="btn btn-secondary border border-(--muted)/20">
                            Joined
                        </button>
                    )}
                </div>
            </div>

            {/* Content Section */}
            {!isMember ? (
                <div className="text-center py-12 bg-(--muted)/5 rounded-xl border border-(--muted)/10">
                    <div className="w-16 h-16 bg-(--muted)/10 rounded-full flex items-center justify-center mx-auto mb-4">
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-(--muted)">
                            <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                        </svg>
                    </div>
                    <h3 className="text-xl font-semibold mb-2">This group is private</h3>
                    <p className="text-(--muted) max-w-md mx-auto mb-6">
                        Join this group to view posts, events, and connect with other members.
                    </p>
                    <button className="btn btn-primary">
                        Join Group
                    </button>
                </div>
            ) : (
                <div className="space-y-8">
                    {/* Tabs (Simple implementation for now) */}
                    <div className="border-b border-(--muted)/10">
                        <div className="flex gap-8">
                            <button className="pb-4 border-b-2 border-(--foreground) font-medium text-(--foreground)">
                                Posts
                            </button>
                            <button className="pb-4 border-b-2 border-transparent font-medium text-(--muted) hover:text-(--foreground) transition-colors">
                                Members
                            </button>
                            <button className="pb-4 border-b-2 border-transparent font-medium text-(--muted) hover:text-(--foreground) transition-colors">
                                Events
                            </button>
                        </div>
                    </div>

                    {/* Posts Feed */}
                    <div className="space-y-6">
                        {posts.map((post) => (
                            <PostCard key={post.ID} post={post} />
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
