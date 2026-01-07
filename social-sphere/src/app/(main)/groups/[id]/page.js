import { getGroup } from "@/actions/groups/get-group";
import { GroupHeader } from "@/components/groups/GroupHeader";
import GroupPageContent from "@/components/groups/GroupPageContent";
import { redirect } from "next/navigation";
import { getGroupPosts } from "@/actions/groups/get-group-posts";

export default async function GroupPage({ params }) {
  const { id } = await params;
  const result = await getGroup(id);
  const response = await getGroupPosts({groupId: id, limit: 10});

  if (!result.success) {
    redirect("/groups");
  }

  const group = result.data;
  const posts = response.data;

  console.log("POSTSSSSS: ", posts);

  return (
    <div className="min-h-screen">
      <GroupHeader group={group} />
      <GroupPageContent group={group} firstPosts={posts} />
    </div>
  );
}