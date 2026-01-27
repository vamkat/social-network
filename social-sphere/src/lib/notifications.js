export function constructLiveNotif(notif) {
    console.log("constructing!", notif)
    
    // FOLLOWERS 
    if (notif.type === "new_follower") {
        return {
            who: notif.payload.follower_name,
            whoID: notif.payload.follower_id,
            message: "followed you"
        };
    }
    if (notif.type === "follow_request") {
        return {
            who: notif.payload.requester_id,
            whoID: notif.payload.requester_name,
            message: "wants to follow you"
        };

    }
    if (notif.type === "follow_request_accepted") {
        return {
            who: notif.payload.taarget_name,
            whoID: notif.payload.target_id,
            message: "accepted your follow request"
         };
    }
    
    // POSTS 
    if (notif.type === "post_reply") {
        return {
            who: notif.payload.commenter_name,
            whoID: notif.payload.commenter_id,
            message: "commented on your",
            where: "post",
            whereID: notif.payload.post_id,
            extra: notif.payload.post_content
        };
    }
    
    if (notif.type === "like") {
        return {
            who: notif.payload.liker_name, 
            whoID: notif.payload.liker_id, 
            message: "liked your", 
            where: "post", 
            whereID: notif.payload.post_id
        };
    }

    // GROUPS
    if (notif.type === "group_invite") {
        return { };
    }

    if (notif.type === "group_join_request") {
        return { };
    }

    if (notif.type === "group_join_request_accepted") {
        return { };
    }

    if (notif.type === "group_invite_accepted") {
        return { };
    }
}