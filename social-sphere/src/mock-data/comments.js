export const getMockComments = () => {
    return [
        {
            ID: "1",
            PostID: "8",
            BasicUserInfo: {
                ID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Ti les re malaka?!?",
            CommentImage: null,
            CreatedAt: "10 minutes ago",
            NumOfHearts: 372847
        },
        {
            ID: "2",
            PostID: "8",
            BasicUserInfo: {
                ID: "1",
                Username: "ychaniot",
                Avatar: "/putin.jpeg"
            },
            Content: "HAHAHAHHAHA",
            CommentImage: null,
            CreatedAt: "21 minutes ago",
            NumOfHearts: 372847
        },
        {
            ID: "3",
            PostID: "8",
            BasicUserInfo: {
                ID: "7",
                Username: "Xi_aomi",
                Avatar: "/xi.jpeg"
            },
            Content: "Confirmed",
            CommentImage: null,
            CreatedAt: "21 minutes ago",
            NumOfHearts: 372847
        },
        {
            ID: "4",
            PostID: "8",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "This guy also told me so",
            CommentImage: "/logos.png",
            CreatedAt: "21 minutes ago",
            NumOfHearts: 372847
        },
        {
            ID: "5",
            PostID: "8",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "lasjnljkndfhaoisuhxkskjdix oiasuhdk jkqzhk kjhsdiqkw jd",
            CommentImage: null,
            CreatedAt: "21 minutes ago",
            NumOfHearts: 372847
        },
        {
            ID: "6",
            PostID: "8",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "yoyoyoyyoyoyoy",
            CommentImage: null,
            CreatedAt: "21 minutes ago",
            NumOfHearts: 372847
        }
    ];
}

export const getLastCommentForPostID = (postID) => {
    const comments = getMockComments().filter(comment => comment.PostID === postID);
    return comments[comments.length - 1];
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