import { getGroup } from "@/actions/groups/get-group";
import { GroupHeader } from "@/components/groups/GroupHeader";
import { redirect } from "next/navigation";

export default async function GroupPage({ params }) {
  const { id } = await params;
  const result = await getGroup(id);

  if (!result.success) {
    redirect("/groups");
  }

  const group = result.data;

  return (
    <div className="min-h-screen">
      <GroupHeader group={group} />

      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-(--muted)">
          <p>Group content area - Posts, Events, and Members sections coming soon...</p>
        </div>
      </div>
    </div>
  );
}