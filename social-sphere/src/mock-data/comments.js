const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export const getMockComments = () => {
    return [
        {
            ID: "1",
            PostID: "1",
            BasicUserInfo: {
                ID: "2",
                Username: "gtoaka",
                Avatar: "/logo.png"
            },
            Content: "This is such a vibe. Need the playlist.",
            CommentImage: null,
            CreatedAt: "2m ago",
            NumOfHearts: 23
        },
        {
            ID: "2",
            PostID: "1",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "Sunday jazz is undefeated.",
            CommentImage: null,
            CreatedAt: "5m ago",
            NumOfHearts: 6
        },
        {
            ID: "3",
            PostID: "1",
            BasicUserInfo: {
                ID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Drop a link to your favorite track?",
            CommentImage: null,
            CreatedAt: "8m ago",
            NumOfHearts: 14
        },
        {
            ID: "4",
            PostID: "1",
            BasicUserInfo: {
                ID: "1",
                Username: "ychaniot",
                Avatar: "/putin.jpeg"
            },
            Content: "Adding this to my morning routine.",
            CommentImage: null,
            CreatedAt: "12m ago",
            NumOfHearts: 9
        },
        {
            ID: "5",
            PostID: "1",
            BasicUserInfo: {
                ID: "7",
                Username: "xi_aomi",
                Avatar: "/xi.jpeg"
            },
            Content: "Coffee + jazz + sunshine = perfection.",
            CommentImage: null,
            CreatedAt: "15m ago",
            NumOfHearts: 18
        },
        {
            ID: "14",
            PostID: "1",
            BasicUserInfo: {
                ID: "3",
                Username: "privateuser",
                Avatar: "/logos.png"
            },
            Content: "This thread is making me crave a croissant.",
            CommentImage: null,
            CreatedAt: "18m ago",
            NumOfHearts: 5
        },
        {
            ID: "15",
            PostID: "1",
            BasicUserInfo: {
                ID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Dropping my go-to playlist in DMs.",
            CommentImage: null,
            CreatedAt: "20m ago",
            NumOfHearts: 16
        },
        {
            ID: "16",
            PostID: "1",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "The calm before the inbox storm.",
            CommentImage: null,
            CreatedAt: "23m ago",
            NumOfHearts: 4
        },
        {
            ID: "17",
            PostID: "1",
            BasicUserInfo: {
                ID: "2",
                Username: "gtoaka",
                Avatar: "/logo.png"
            },
            Content: "Need a photo of that breakfast spread!",
            CommentImage: null,
            CreatedAt: "27m ago",
            NumOfHearts: 10
        },
        {
            ID: "6",
            PostID: "2",
            BasicUserInfo: {
                ID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Energy levels high today!",
            CommentImage: null,
            CreatedAt: "3m ago",
            NumOfHearts: 7
        },
        {
            ID: "7",
            PostID: "2",
            BasicUserInfo: {
                ID: "3",
                Username: "privateuser",
                Avatar: "/logos.png"
            },
            Content: "Good morning to everyone except my deadlines.",
            CommentImage: null,
            CreatedAt: "9m ago",
            NumOfHearts: 2
        },
        {
            ID: "8",
            PostID: "2",
            BasicUserInfo: {
                ID: "2",
                Username: "gtoaka",
                Avatar: "/logo.png"
            },
            Content: "Let's grab brunch later?",
            CommentImage: null,
            CreatedAt: "14m ago",
            NumOfHearts: 4
        },
        {
            ID: "9",
            PostID: "3",
            BasicUserInfo: {
                ID: "1",
                Username: "ychaniot",
                Avatar: "/putin.jpeg"
            },
            Content: "Spicy take but you might be right.",
            CommentImage: null,
            CreatedAt: "30m ago",
            NumOfHearts: 31
        },
        {
            ID: "10",
            PostID: "3",
            BasicUserInfo: {
                ID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "Robots can't steal my charisma.",
            CommentImage: null,
            CreatedAt: "32m ago",
            NumOfHearts: 11
        },
        {
            ID: "18",
            PostID: "3",
            BasicUserInfo: {
                ID: "7",
                Username: "xi_aomi",
                Avatar: "/xi.jpeg"
            },
            Content: "Watching this from the sidelines with popcorn.",
            CommentImage: null,
            CreatedAt: "35m ago",
            NumOfHearts: 19
        },
        {
            ID: "11",
            PostID: "5",
            BasicUserInfo: {
                ID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Need the route name! That view though.",
            CommentImage: null,
            CreatedAt: "1h ago",
            NumOfHearts: 3
        },
        {
            ID: "12",
            PostID: "5",
            BasicUserInfo: {
                ID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Props for finishing it. I bailed halfway last time.",
            CommentImage: null,
            CreatedAt: "1h ago",
            NumOfHearts: 8
        },
        {
            ID: "13",
            PostID: "8",
            BasicUserInfo: {
                ID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Ti les re malaka?!?",
            CommentImage: null,
            CreatedAt: "10m ago",
            NumOfHearts: 372847
        }
    ];
};

export const getCommentsForPost = (postID) => {
    return getMockComments().filter((comment) => comment.PostID === postID);
};

export const getLastCommentForPostID = (postID) => {
    const comments = getCommentsForPost(postID);
    return comments[comments.length - 1];
};

export const getPaginatedComments = (postID, cursor = 0, limit = 3) => {
    // newest-first ordering
    const comments = getCommentsForPost(postID).slice().reverse();
    const slice = comments.slice(cursor, cursor + limit);
    return {
        comments: slice,
        hasMore: cursor + slice.length < comments.length,
        nextCursor: cursor + slice.length,
        total: comments.length,
    };
};

export const fetchPaginatedComments = async (postID, cursor = 0, limit = 3) => {
    await delay(450);
    return getPaginatedComments(postID, cursor, limit);
};

export const countCommentsForPost = (postID) => getCommentsForPost(postID).length;
