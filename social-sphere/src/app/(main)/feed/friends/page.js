import { LogoutButton } from "@/components/LogoutButton";
// import { createGroup } from "@/services/groups/createGroup";
// import { getAllGroups } from "@/services/groups/getAllGroups";
// import { getUserGroups } from "@/services/groups/getUserGroups";


export const metadata = {
    title: "Friends Feed",
}

// export async function createaGroup() {
//     const resp = await createGroup({
//         group_title: "Friendsdv sgfsdvbsdbbsgfbsfv sbv d ",
//         group_description: "Frierbgsdfdsbfgbends fv sfv gsfv sfbroup",
//         group_image: "hello",
//     });
//     return resp;
// }

export default async function FriendsFeedPage() {
    // const resp = await createaGroup();
    // console.log(resp);

    // const resp = await getUserGroups({ limit: 10, offset: 0 });
    // console.log(resp);

    return (
        <div>
            <LogoutButton />
        </div>
    );
}