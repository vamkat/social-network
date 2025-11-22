export const getMockUser = (username) => {
    return {
        Username: username,
        firstName: username === "ychaniot" ? "Ypatios" : username === "gtoaka" ? "Georgia" : username === "privatefollowed" ? "Private" : "Private",
        lastName: username === "ychaniot" ? "Chaniotakos" : username === "gtoaka" ? "Toaka" : username === "privatefollowed" ? "Friend" : "User",
        AboutMe: "Digital explorer & coffee enthusiast. Building things that matter.",
        location: "San Francisco, CA",
        DateOfBirth: "1990-01-01",
        CreatedAt: "2023-01-15",
        Avatar: null,
        isFollower: username === "privatefollowed",
        publicProf: username !== "privateuser" && username !== "privatefollowed",
        FollowersNum: 1234,
        FollowingNum: 567,
        GroupsNum: 12,
    };
};
