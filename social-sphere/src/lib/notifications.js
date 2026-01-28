import nextConfig from "../../next.config.mjs";

export function constructLiveNotif(notif) {
    console.log("constructing!", notif)
    
    // FOLLOWERS 
    if (notif.type === "new_follower") {
        return {
            who: notif.payload.follower_name,
            whoID: notif.payload.follower_id,
            message: " is now following you"
        };
    }
    if (notif.type === "follow_request") {
        return {
            who: notif.payload.requester_id,
            whoID: notif.payload.requester_name,
            message: " wants to follow you"
        };

    }
    if (notif.type === "follow_request_accepted") {
        return {
            who: notif.payload.taarget_name,
            whoID: notif.payload.target_id,
            message: " accepted your follow request"
         };
    }
    
    // POSTS 
    if (notif.type === "post_reply") {
        return {
            who: notif.payload.commenter_name,
            whoID: notif.payload.commenter_id,
            message: " commented on your post: ",
            wherePost: notif.payload.post_content,
            whereID: notif.payload.post_id,
        };
    }
    
    if (notif.type === "like") {
        return {
            who: notif.payload.liker_name, 
            whoID: notif.payload.liker_id, 
            message: " liked your ", 
            wherePost: "post", 
            whereID: notif.payload.post_id
        };
    }

    // GROUPS
    if (notif.type === "group_invite") {
        return {
            who: notif.payload.inviter_name,
            whoID: notif.payload.inviter_id,
            message: " invited you to join group: ",
            whereGroup: notif.payload.group_name,
            whereID: notif.payload.group_id
         };
    }

    if (notif.type === "group_join_request") {
        return {
            who: notif.payload.requester_name,
            whoID: notif.payload.requester_id,
            message: " wants to join your group: ",
            whereGroup: notif.payload.group_name,
            whereID: notif.payload.group_id
         };
    }

    if (notif.type === "group_join_request_accepted") {
        return {
            message: "You were accepted to join group: ",
            whereGroup: notif.payload.group_name,
            whereID: notif.payload.group_id
         };
    }

    if (notif.type === "group_invite_accepted") {
        return {
            who: notif.payload.invited_name,
            whoID: notif.payload.invited_id,
            message: " accepted your invitation to join group: ",
            whereGroup: notif.payload.group_name,
            whereID: notif.payload.group_id
         };
    }
}