import { getProfileInfo } from "@/actions/profile/get-profile-info";
import { redirect } from "next/navigation";
import SettingsClient from "./SettingsClient";

export default async function SettingsPage({ params }) {
    const { id } = await params;

    // get user's info
    const user = await getProfileInfo(id);

    if (user.user_id !== id) { redirect(`/profile/${id}`); }

    return <SettingsClient user={user} />;
}
