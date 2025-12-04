"use client";

import { createContext, useContext, useState, useEffect } from "react";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // Function to fetch user profile via proxy API
    const fetchUserProfile = async (userId) => {
        try {
            // Call Next.js proxy API route instead of directly calling backend
            const response = await fetch(`/api/auth/profile/${userId}`, {
                method: "GET",
                credentials: "include", // Important: include cookies for authentication
            });

            if (response.ok) {
                const profileData = await response.json();
                setUser(profileData);
                return profileData;
            } else {
                console.error("Failed to fetch user profile");
                return null;
            }
        } catch (error) {
            console.error("Error fetching user profile:", error);
            return null;
        }
    };

    // Function to update user profile in context
    const updateUser = (userData) => {
        setUser(userData);
    };

    // Function to clear user on logout
    const clearUser = () => {
        setUser(null);
    };

    // Check authentication status on mount
    useEffect(() => {
        // You can add logic here to check if user is authenticated
        // For now, we'll just set loading to false
        setLoading(false);
    }, []);

    return (
        <AuthContext.Provider
            value={{
                user,
                loading,
                fetchUserProfile,
                updateUser,
                clearUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

// Custom hook to use auth context
export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
