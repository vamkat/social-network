"use client";

import { usePathname } from "next/navigation";
import { Activity, Users, Send, Bell, User, LogOut, Settings, HeartPulse, Search, Loader2 } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import Tooltip from "@/components/ui/Tooltip";
import Link from "next/link";
import { useStore } from "@/store/store";
import { logout } from "@/services/auth/logout";
import { useRouter } from "next/navigation";
import { SearchUsers } from "@/services/search/searchUsers";

export default function Navbar() {
    const pathname = usePathname();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const dropdownRef = useRef(null);
    const searchRef = useRef(null);
    const clearUser = useStore((state) => state.clearUser);
    const router = useRouter();

    const user = useStore((state) => state.user);

    // Search State
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState([]);
    const [isSearching, setIsSearching] = useState(false);
    const [showSearchResults, setShowSearchResults] = useState(false);

    // Close dropdowns when clicking outside
    useEffect(() => {
        function handleClickOutside(event) {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsDropdownOpen(false);
            }
            if (searchRef.current && !searchRef.current.contains(event.target)) {
                setShowSearchResults(false);
            }
        }

        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    // Debounced Search
    useEffect(() => {
        const timer = setTimeout(async () => {
            if (searchQuery.trim().length >= 2) {
                setIsSearching(true);
                try {
                    const response = await SearchUsers({ query: searchQuery, limit: 5 });
                    if (response && response.users) {
                        setSearchResults(response.users);
                        setShowSearchResults(true);
                    } else {
                        setSearchResults([]);
                    }
                } catch (error) {
                    console.error("Search error:", error);
                    setSearchResults([]);
                } finally {
                    setIsSearching(false);
                }
            } else {
                setSearchResults([]);
                setShowSearchResults(false);
            }
        }, 300); // 300ms debounce

        return () => clearTimeout(timer);
    }, [searchQuery]);

    const handleLogout = async () => {
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

    if (!user) {
        return <div></div>;
    }

    const navItems = [
        {
            label: "Public",
            href: "/feed/public",
            icon: Activity,
        },
        {
            label: "Friends",
            href: "/feed/friends",
            icon: HeartPulse,
        },
        {
            label: "Groups",
            href: "/groups",
            icon: Users,
        },
    ];

    const isActive = (path) => pathname === path;

    return (
        <nav className="sticky top-0 z-50 w-full border-b border-(--border) bg-(--background)/95 backdrop-blur-md">
            <div className="w-full px-3 sm:px-4 md:px-6">
                <div className="flex items-center justify-between h-16 gap-3">
                    {/* Left Section: Logo + Search */}
                    <div className="flex items-center gap-70 flex-1">
                        {/* Logo - Hidden on smallest screens */}
                        <Link
                            href="/feed/public"
                            className="hidden sm:flex items-center shrink-0"
                        >
                            <span className="hidden md:block text-base font-medium tracking-tight text-foreground hover:text-(--muted) transition-colors">
                                SocialSphere
                            </span>
                        </Link>

                        {/* Desktop Search Bar */}
                        <div className="hidden lg:flex flex-1 max-w-md" ref={searchRef}>
                            <div className="relative w-full group">
                                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                    {isSearching ? (
                                        <Loader2 className="h-4 w-4 text-(--muted) animate-spin" />
                                    ) : (
                                        <Search className="h-4 w-4 text-(--muted) group-focus-within:text-(--accent) transition-colors" />
                                    )}
                                </div>
                                <input
                                    type="text"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    onFocus={() => {
                                        if (searchResults.length > 0) setShowSearchResults(true);
                                    }}
                                    className="block w-full pl-11 pr-4 py-2.5 border border-(--border) rounded-full text-sm bg-(--muted)/5 text-foreground placeholder-(--muted) hover:border-foreground focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all"
                                    placeholder="Search users..."
                                />

                                {/* Search Results Dropdown */}
                                {showSearchResults && (
                                    <div className="absolute top-full left-0 right-0 mt-2 bg-background border border-(--border) rounded-2xl shadow-xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                                        {searchResults.length > 0 ? (
                                            <div className="py-2">
                                                {searchResults.map((result) => (
                                                    <Link
                                                        key={result.user_id}
                                                        href={`/profile/${result.user_id}`}
                                                        onClick={() => {
                                                            setShowSearchResults(false);
                                                            setSearchQuery("");
                                                        }}
                                                        className="flex items-center gap-3 px-4 py-3 hover:bg-(--muted)/5 transition-colors"
                                                    >
                                                        <div className="w-10 h-10 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden shrink-0">
                                                            {/* {result.avatar_id ? (
                                                                <img src={result.avatar_id} alt={result.username} className="w-full h-full object-cover" />
                                                            ) : (
                                                                <User className="w-5 h-5 text-(--muted)" />
                                                            )} */}
                                                        </div>
                                                        <div className="flex-1 min-w-0">
                                                            <p className="text-sm font-medium text-foreground truncate">
                                                                {result.username}
                                                            </p>
                                                        </div>
                                                    </Link>
                                                ))}
                                            </div>
                                        ) : (
                                            <div className="p-4 text-center text-sm text-(--muted)">
                                                No users found
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Right Section: Nav + Actions */}
                    <div className="flex items-center gap-1.5 shrink-0">
                        {/* Desktop Navigation */}
                        <div className="hidden md:flex items-center gap-1">
                            {navItems.map((item) => {
                                const Icon = item.icon;
                                const active = isActive(item.href);
                                return (
                                    <Tooltip key={item.href} content={item.label}>
                                        <Link
                                            href={item.href}
                                            className={`flex items-center gap-2 px-3 py-2 rounded-full text-sm font-medium transition-all ${active
                                                ? "bg-(--accent)/10 text-(--accent)"
                                                : "text-(--muted) hover:text-foreground hover:bg-(--muted)/10"
                                                }`}
                                        >
                                            <Icon className="w-[18px] h-[18px]" strokeWidth={active ? 2.5 : 2} />
                                        </Link>
                                    </Tooltip>
                                );
                            })}
                        </div>

                        {/* Mobile Navigation - Icon only */}
                        <div className="flex md:hidden items-center gap-0.5">
                            {navItems.map((item) => {
                                const Icon = item.icon;
                                const active = isActive(item.href);
                                return (
                                    <Tooltip key={item.href} content={item.label}>
                                        <Link
                                            href={item.href}
                                            className={`p-2.5 rounded-full transition-all ${active
                                                ? "bg-(--accent)/10 text-(--accent)"
                                                : "text-(--muted) hover:text-foreground hover:bg-(--muted)/10"
                                                }`}
                                        >
                                            <Icon className="w-5 h-5" strokeWidth={active ? 2.5 : 2} />
                                        </Link>
                                    </Tooltip>
                                );
                            })}
                        </div>

                        {/* Divider */}
                        <div className="h-6 w-px bg-(--border) mx-1" />

                        {/* Messages */}
                        <Tooltip content="Messages">
                            <Link
                                href="/messages"
                                className={`relative p-2.5 rounded-full transition-all ${isActive('/messages')
                                    ? "bg-(--accent)/10 text-(--accent)"
                                    : "text-(--muted) hover:text-foreground hover:bg-(--muted)/10"
                                    }`}
                            >
                                <Send className="w-5 h-5" strokeWidth={isActive('/messages') ? 2.5 : 2} />
                                <span className="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 text-[10px] font-bold text-white bg-red-500 rounded-full flex items-center justify-center border-2 border-background">
                                    1
                                </span>
                            </Link>
                        </Tooltip>

                        {/* Notifications */}
                        <Tooltip content="Notifications">
                            <Link
                                href="/notifications"
                                className={`relative p-2.5 rounded-full transition-all ${isActive('/notifications')
                                    ? "bg-(--accent)/10 text-(--accent)"
                                    : "text-(--muted) hover:text-foreground hover:bg-(--muted)/10"
                                    }`}
                            >
                                <Bell className="w-5 h-5" strokeWidth={isActive('/notifications') ? 2.5 : 2} />
                                <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 border-background" />
                            </Link>
                        </Tooltip>

                        {/* User Dropdown */}
                        {user && (
                            <div className="relative ml-1.5 pl-2.5 border-l border-(--border)" ref={dropdownRef}>
                                <button
                                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                                    className="flex items-center gap-2 hover:opacity-80 transition-opacity"
                                >
                                    <div className="w-8 h-8 rounded-full bg-(--muted)/10 border border-(--border) flex items-center justify-center overflow-hidden hover:border-(--accent) transition-colors">
                                        {user.avatar ? (
                                            <img src={user.avatar} alt={user.username[0]} className="w-full h-full object-cover" />
                                        ) : (
                                            <User className="w-4 h-4 text-(--muted)" />
                                        )}
                                    </div>
                                    <svg
                                        className={`hidden sm:block w-3.5 h-3.5 text-(--muted) transition-transform ${isDropdownOpen ? "rotate-180" : ""
                                            }`}
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                    </svg>
                                </button>

                                {/* Dropdown Menu */}
                                {isDropdownOpen && (
                                    <div className="absolute right-0 top-full mt-3 w-52 rounded-2xl border border-(--border) bg-background shadow-xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                                        <div className="p-1.5">
                                            <Link
                                                href={`/profile/${user.id}`}
                                                onClick={() => setIsDropdownOpen(false)}
                                                className="flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl hover:bg-(--muted)/10 transition-colors text-foreground"
                                            >
                                                <User className="w-4 h-4 text-(--muted)" />
                                                Profile
                                            </Link>
                                            <Link
                                                href={`/profile/${user.id}/settings`}
                                                onClick={() => setIsDropdownOpen(false)}
                                                className="flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl hover:bg-(--muted)/10 transition-colors text-foreground"
                                            >
                                                <Settings className="w-4 h-4 text-(--muted)" />
                                                Settings
                                            </Link>
                                            <div className="h-px bg-(--border) my-1.5" />
                                            <button
                                                onClick={() => {
                                                    setIsDropdownOpen(false);
                                                    handleLogout();
                                                }}
                                                className="w-full flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl text-red-500 hover:bg-red-500/10 transition-colors text-left"
                                            >
                                                <LogOut className="w-4 h-4" />
                                                Sign Out
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>

                {/* Mobile Search Bar - Below main nav */}
                <div className="lg:hidden pb-3">
                    {/* Placeholder for now - logic can be mirrored if needed or hidden */}
                    <div className="relative w-full group">
                        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                            <Search className="h-4 w-4 text-(--muted) group-focus-within:text-(--accent) transition-colors" />
                        </div>
                        <input
                            type="text"
                            /* Mobile search logic can be implemented here similarly to desktop if requested */
                            className="block w-full pl-11 pr-4 py-2.5 border border-(--border) rounded-full text-sm bg-(--muted)/5 text-foreground placeholder-(--muted) hover:border-foreground focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all"
                            placeholder="Search users..."
                        />
                    </div>
                </div>
            </div>
        </nav>
    );
}