const normalizeVisibility = (post) =>
    (post.visibility ?? post.Visibility ?? "public").toString().toLowerCase();

export const getMockPosts = (offset = 0, limit = 10, visibilityFilter = null, groupId = null) => {
    const posts = [
        {
            ID: "1",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Sunday mornings are for slow breakfasts and even slower jazz. There's something magical about the quiet before the city wakes up.",
            PostImage: "/elon.jpeg",
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
            LastComment: {
                CommentId: "101",
                Body: "Love this vibe.",
                Creator: {
                    UserID: "2",
                    Username: "sunny_side",
                    Avatar: "/elon.jpeg"
                },
                ReactionsCount: 2,
                LikedByUser: false,
                Image: null,
                CreatedAt: "2m ago"
            }
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
            LastComment: {
                CommentId: "102",
                Body: "Good morning!",
                Creator: {
                    UserID: "3",
                    Username: "earlybird",
                    Avatar: "/elon.jpeg"
                },
                ReactionsCount: 1,
                LikedByUser: false,
                Image: null,
                CreatedAt: "5m ago"
            }
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
            LastComment: {
                CommentId: "103",
                Body: "Spicy take.",
                Creator: {
                    UserID: "5",
                    Username: "trumpet",
                    Avatar: "/trump.jpeg"
                },
                ReactionsCount: 4,
                LikedByUser: false,
                Image: null,
                CreatedAt: "8m ago"
            }
        },
        {
            ID: "4",
            BasicUserInfo: {
                UserID: "4",
                Username: "watermelon_musk",
                Avatar: "/elon.jpeg"
            },
            Content: "Sunday mornings are for slow breakfasts and even slower jazz. There's something magical about the quiet before the city wakes up.",
            PostImage: null,
            CreatedAt: "10m ago",
            NumOfComments: 3,
            NumOfHearts: 1345,
            IsHearted: false,
            LastComment: {
                CommentId: "104",
                Body: "Craving brunch now.",
                Creator: {
                    UserID: "6",
                    Username: "kimpossible",
                    Avatar: "/kim.jpeg"
                },
                ReactionsCount: 0,
                LikedByUser: false,
                Image: null,
                CreatedAt: "12m ago"
            }
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
            LastComment: {
                CommentId: "105",
                Body: "Need to see pics!",
                Creator: {
                    UserID: "7",
                    Username: "Xi_aomi",
                    Avatar: "/xi.jpeg"
                },
                ReactionsCount: 3,
                LikedByUser: false,
                Image: null,
                CreatedAt: "2h ago"
            }
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
            LastComment: {
                CommentId: "106",
                Body: "Legend indeed.",
                Creator: {
                    UserID: "1",
                    Username: "ychaniot",
                    Avatar: "/putin.jpeg"
                },
                ReactionsCount: 2,
                LikedByUser: false,
                Image: null,
                CreatedAt: "1h ago"
            }
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
            LastComment: {
                CommentId: "107",
                Body: "Whoa, harsh.",
                Creator: {
                    UserID: "2",
                    Username: "sunny_side",
                    Avatar: "/elon.jpeg"
                },
                ReactionsCount: 5,
                LikedByUser: false,
                Image: null,
                CreatedAt: "1h ago"
            }
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
            LastComment: {
                CommentId: "108",
                Body: "This escalated fast.",
                Creator: {
                    UserID: "6",
                    Username: "kimpossible",
                    Avatar: "/kim.jpeg"
                },
                ReactionsCount: 1,
                LikedByUser: false,
                Image: null,
                CreatedAt: "1h ago"
            }
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
            LastComment: {
                CommentId: "109",
                Body: "Same here!",
                Creator: {
                    UserID: "4",
                    Username: "watermelon_musk",
                    Avatar: "/elon.jpeg"
                },
                ReactionsCount: 7,
                LikedByUser: false,
                Image: null,
                CreatedAt: "3h ago"
            }
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
            LastComment: {
                CommentId: "110",
                Body: "Donut well deserved.",
                Creator: {
                    UserID: "5",
                    Username: "trumpet",
                    Avatar: "/trump.jpeg"
                },
                ReactionsCount: 0,
                LikedByUser: false,
                Image: null,
                CreatedAt: "12h ago"
            }
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
            LastComment: null,
            ReactionsCount: 1,
            LikedByUser: false,
            Image: null,
            CreatedAt: "20h ago"
        }
    ];

    let filtered = posts;

    if (visibilityFilter?.length) {
        const allowed = visibilityFilter.map((v) => v.toLowerCase());
        filtered = filtered.filter((post) => allowed.includes(normalizeVisibility(post)));
    }

    if (groupId) {
        const scoped = filtered.filter((post) => (post.GroupID ?? post.group_id) === groupId);
        filtered = scoped.length > 0 ? scoped : filtered;
    }

    return filtered.slice(offset, offset + limit);
};

export const GetPostsByUserId = (userId, offset = 0, limit = 10) => {
    const allPosts = getMockPosts(0, 1000);
    const userPosts = allPosts.filter(post => post.BasicUserInfo.UserID === userId);
    return userPosts.slice(offset, offset + limit);
};
