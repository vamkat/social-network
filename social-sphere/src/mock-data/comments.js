export const getMockComments = () => {
    return [
        {
            CommentId: "1",
            PostID: "8",
            Creator: {
                UserID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Body: "Ti les re malaka?!?",
            Image: null,
            CreatedAt: "10 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        },
        {
            CommentId: "2",
            PostID: "8",
            Creator: {
                UserID: "1",
                Username: "ychaniot",
                Avatar: "/putin.jpeg"
            },
            Body: "HAHAHAHHAHA",
            Image: null,
            CreatedAt: "21 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        },
        {
            CommentId: "3",
            PostID: "8",
            Creator: {
                UserID: "7",
                Username: "Xi_aomi",
                Avatar: "/xi.jpeg"
            },
            Body: "Confirmed",
            Image: null,
            CreatedAt: "21 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        },
        {
            CommentId: "4",
            PostID: "8",
            Creator: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Body: "This guy also told me so",
            Image: "/logos.png",
            CreatedAt: "21 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        },
        {
            CommentId: "5",
            PostID: "8",
            Creator: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Body: "lasjnljkndfhaoisuhxkskjdix oiasuhdk jkqzhk kjhsdiqkw jd",
            Image: null,
            CreatedAt: "21 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        },
        {
            CommentId: "6",
            PostID: "8",
            Creator: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Body: "yoyoyoyyoyoyoy",
            Image: null,
            CreatedAt: "21 minutes ago",
            ReactionsCount: 372847,
            LikedByUser: false
        }
    ];
}

export const getCommentsForPost = (postID, offset = 0, limit = 2) => {
    // Return comments in reverse order (newest first) to facilitate "load previous"
    // But wait, usually comments are displayed oldest to newest (top to bottom).
    // If we want "load previous", we want to fetch the ones *before* the current ones.
    // Let's assume the mock data is sorted by creation time (oldest first).
    // So "last" is the newest.
    // If we have [1, 2, 3, 4, 5, 6] (6 is newest).
    // Initial view: 6.
    // Load previous: want 4, 5.
    // Load previous again: want 2, 3.
    // So we need to reverse the array to slice from the end?
    // Reversed: [6, 5, 4, 3, 2, 1]
    // Offset 0, Limit 1 -> [6] (Initial)
    // Offset 1, Limit 2 -> [5, 4] (Next batch)
    // Offset 3, Limit 2 -> [3, 2] (Next batch)

    const comments = getMockComments().filter(comment => comment.PostID === postID);
    const reversedComments = [...comments].reverse();
    return reversedComments.slice(offset, offset + limit);
}
