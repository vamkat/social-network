export const getMockNotifications = () => {
  return [
    {
      id: 1,
      type: "follower",
      title: "@john wants to follow you",
      ctaLabels: ["Accept", "Decline"],
      createdAt: "2h",
      isUnread: true,
    },

    {
      id: 2,
      type: "event",
      title: "New event in Group Name",
      ctaLabels: ["View"],
      createdAt: "1h",
      isUnread: false,
    },
  ];
};

export const getNotificationsId = (notificationId) => {
  return getMockNotifications().find((notification) => notification.id === notificationId);
};
