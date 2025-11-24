// src/providers/AuthProvider.js
"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api/client";
import { AUTH_ENDPOINTS } from "@/lib/api/endpoints";

const AuthContext = createContext(undefined);

export function AuthProvider({ children }) {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check session on mount
  useEffect(() => {
    let isMounted = true;

    async function fetchMe() {
      try {
        const res = await apiClient.get(AUTH_ENDPOINTS.me);
        if (isMounted) {
          setUser(res.data.user || null);
        }
      } catch (err) {
        if (isMounted) {
          setUser(null);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    fetchMe();

    return () => {
      isMounted = false;
    };
  }, []);

  async function login(credentials) {
    // credentials: { identifier, password } (or email/password depending on backend)
    await apiClient.post(AUTH_ENDPOINTS.login, credentials);
    // After successful login, refresh "me"
    const res = await apiClient.get(AUTH_ENDPOINTS.me);
    setUser(res.data.user || null);
    router.push("/feed/public");
  }

  async function logout() {
    try {
      await apiClient.post(AUTH_ENDPOINTS.logout);
    } catch (_) {
      // ignore
    }
    setUser(null);
    router.push("/login");
  }

  async function register(data) {
    await apiClient.post(AUTH_ENDPOINTS.register, data);
    const res = await apiClient.get(AUTH_ENDPOINTS.me);
    setUser(res.data.user || null);
    router.push("/feed/public");
  }

  const value = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    logout,
    register,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside <AuthProvider>");
  }
  return ctx;
}
