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
        }
    ];
}

export const getLastCommentForPostID = (postID) => {
    const comments = getMockComments().filter(comment => comment.PostID === postID);
    return comments[comments.length - 1];
}