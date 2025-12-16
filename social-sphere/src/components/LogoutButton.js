'use client';

import { logout } from "@/actions/auth/logout";
import { useState } from "react";
import { useStore } from "@/store/store";
import { useRouter } from "next/navigation";

export function LogoutButton() {
    const [isLoading, setIsLoading] = useState(false)
    const clearUser = useStore((state) => state.clearUser);
    const router = useRouter();

    const handleLogout = async () => {
        setIsLoading(true)

        try {
            // logout
            const resp = await logout();

            if (!resp.success) {
                console.error('error:', resp.error);
            }

            // clear user from state and local storage
            clearUser();

            // Redirect to login
            router.push("/login");

        } catch (error) {
            console.error('Logout error:', error)
        }
    }

    return (
        <div className="flex items-center justify-center min-h-screen">
            <button
                className="w-1/2 flex justify-center btn btn-primary"
                onClick={handleLogout}
                disabled={isLoading}
            >
                {isLoading ? 'Logging out...' : 'Logout'}
            </button>
        </div>
    )
}