export const getMockPosts = () => {
    return [
        {
            ID: "1",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Sunday mornings are for slow breakfasts and even slower jazz. 🎷🥐 There's something magical about the quiet before the city wakes up.",
            PostImage: "/elon.jpeg",
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
        },
        {
            ID: "2",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Hello world! Good morning people!!",
            PostImage: null,
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
        },
        {
            ID: "3",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "AI will take our jobs. Basically, It will take YOUR jobs. I dont need to work. mouhahaha",
            PostImage: "/logos.png",
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
        },
        {
            ID: "4",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Sunday mornings are for slow breakfasts and even slower jazz. 🎷🥐 There's something magical about the quiet before the city wakes up.",
            PostImage: null,
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
        },
        {
            ID: "5",
            BasicUserInfo: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "Finally hiked the trail I've been looking at for months. The view from the top was absolutely worth the struggle. Nature has a way of resetting your perspective.",
            PostImage: null,
            CreatedAt: "3h ago",
            NumOfComments: 15,
            NumOfHearts: 6,
            IsHearted: false,
        },
        {
            ID: "6",
            BasicUserInfo: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "This guy is a legend",
            PostImage: "/trump.jpeg",
            CreatedAt: "3h ago",
            NumOfComments: 15,
            NumOfHearts: 6,
            IsHearted: false,
        },
        {
            ID: "7",
            BasicUserInfo: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "This guy is a fat pig",
            PostImage: "/elon.jpeg",
            CreatedAt: "3h ago",
            NumOfComments: 15,
            NumOfHearts: 6,
            IsHearted: false,
        },
        {
            ID: "8",
            BasicUserInfo: {
                UserID: "5",
                Username: "trumpet",
                Avatar: "/trump.jpeg"
            },
            Content: "this guy wants nuclear war with US",
            PostImage: "/kim.jpeg",
            CreatedAt: "3h ago",
            NumOfComments: 15,
            NumOfHearts: 6,
            IsHearted: false,
        },
        {
            ID: "9",
            BasicUserInfo: {
                UserID: "6",
                Username: "kimpossible",
                Avatar: "/kim.jpeg"
            },
            Content: "Does anyone else feel like time is moving exceptionally fast lately? I swear it was January just yesterday.",
            PostImage: null,
            CreatedAt: "5h ago",
            NumOfComments: 42,
            NumOfHearts: 145,
            IsHearted: false,
        },
        {
            ID: "10",
            BasicUserInfo: {
                UserID: "7",
                Username: "Xi_aomi",
                Avatar: "/xi.jpeg"
            },
            Content: "Small wins matter. Fixed a bug that's been bugging me (pun intended) for a week. Celebrating with a donut.",
            PostImage: null,
            CreatedAt: "1d ago",
            NumOfComments: 7,
            NumOfHearts: 12,
            IsHearted: false,
        },
        {
            ID: "11",
            BasicUserInfo: {
                UserID: "1",
                Username: "ychaniot",
                Avatar: "/putin.jpeg"
            },
            Content: "What the fuck is going on??",
            PostImage: null,
            CreatedAt: "1d ago",
            NumOfComments: 7,
            NumOfHearts: 12,
            IsHearted: false,
        }
    ];
}

export const GetPostsByUserId = (userId) => {
    return getMockPosts().filter(post => post.BasicUserInfo.UserID === userId);
}