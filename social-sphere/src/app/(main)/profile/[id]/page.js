import { getProfileInfo } from "@/actions/profile/get-profile-info";
import ProfileContent from "@/components/profile/ProfileContent";
import { getUserPosts } from "@/actions/posts/get-user-posts";

async function getUserProfile(userId) {
    try {
        const user = await getProfileInfo(userId);
        return { success: true, user };
    } catch (error) {
        return { success: false, error: error.message };
    }
}

export async function generateMetadata({ params }) {
    const { id } = await params;
    const result = await getUserProfile(id);

    if (!result.success || !result.user) {
        return { title: "Profile" };
    }

    return {
        title: `${result.user.username}'s Profile`,
        description: `View ${result.user.first_name} ${result.user.last_name}'s profile`,
    };
}

export default async function ProfilePage({ params }) {
    const { id } = await params;
    const result = await getUserProfile(id);
    // Fetch initial 10 posts for infinite scrolling
    const posts = await getUserPosts({ creatorId: id, limit: 10, offset: 0 });

    // Pass the result object (contains success, user, or error)
    return <ProfileContent result={result} posts={posts} />;
}