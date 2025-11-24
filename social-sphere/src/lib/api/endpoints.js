// src/lib/api/endpoints.js

export const AUTH_ENDPOINTS = {
  register: "/api/v1/auth/register",
  login: "/api/v1/auth/login",
  logout: "/api/v1/auth/logout",
  me: "/api/v1/auth/me",
};

export const POST_ENDPOINTS = {
  feed: (type) => `/api/v1/posts/${type}`, // "public" | "friends"
  create: "/api/v1/posts",
};
