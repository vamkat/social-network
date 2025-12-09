"use client";

import { usePathname } from "next/navigation";
import { useSession } from "next-auth/react";
import { Activity, Users, Send, Bell, User, LogOut, Settings, Menu, X, HeartPulse, Search } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import Tooltip from "@/components/ui/tooltip";
import Link from "next/link";
import { logoutClient } from "@/services/auth/logout-client";

export default function Navbar() {
  const pathname = usePathname();
  const { data: session } = useSession();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(null);
  const dropdownRef = useRef(null);

  const user = session?.user;
  console.log(user);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen((prev) => prev ? false : prev);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

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
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center gap-4">
        {/* Left: Logo */}
        <div className="flex items-center justify-start">
          <Link href="/feed/public" className="flex items-center gap-2.5 group w-fit">
            <span className="hidden md:block text-base font-medium tracking-tight text-(--foreground) group-hover:text-(--accent) transition-colors">
              SocialSphere
            </span>
          </Link>
        </div>

        {/* Center: Search Bar */}
        <div className="pl-46 hidden lg:block w-full max-w-2xl mr-auto">
          <div className="relative w-full group">
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <Search className="h-4 w-4 text-(--muted) group-focus-within:text-(--accent) transition-colors" />
            </div>
            <input
              type="text"
              className="block w-full pl-11 pr-4 py-2.5 border border-(--border) rounded-full text-sm bg-(--muted)/5 text-(--foreground) placeholder-(--muted) hover:border-(--foreground) focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all"
              placeholder="Search users..."
            />
          </div>
        </div>

        {/* Right: Nav Links & Actions */}
        <div className="flex items-center gap-1.5">

          {/* Desktop Nav Links */}
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.href);
            return (
              <Tooltip key={item.href} content={item.label}>
                <Link
                  href={item.href}
                  className={`relative p-2.5 rounded-full flex items-center justify-center transition-all group ${active
                    ? "text-(--accent) bg-(--accent)/10"
                    : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10"
                    }`}
                >
                  <Icon className={`w-5 h-5 transition-transform group-hover:scale-110 ${active ? "stroke-[2.5px]" : "stroke-2"}`} />
                </Link>
              </Tooltip>
            );
          })}


          {/* Messages */}
          <Tooltip content="Messages">
            <button
              className={`p-2.5 rounded-full transition-all relative group ${isActive('/messages')
                ? "text-(--accent) bg-(--accent)/10"
                : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10"}`}
            >
              <Send className={`w-5 h-5 transition-transform group-hover:scale-110 ${isActive('/messages') ? "stroke-[2.5px]" : "stroke-2"}`} />
              <span className="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 text-[10px] font-bold text-white bg-red-500 rounded-full flex items-center justify-center border-2 border-(--background)">1</span>
            </button>
          </Tooltip>

          {/* Notifications */}
          <Tooltip content="Notifications">
            <button
              className={`p-2.5 rounded-full transition-all relative group ${isActive('/notifications')
                ? "text-(--accent) bg-(--accent)/10"
                : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10"}`}
            >
              <Bell className={`w-5 h-5 transition-transform group-hover:scale-110 ${isActive('/notifications') ? "stroke-[2.5px]" : "stroke-2"}`} />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 border-(--background)" />
            </button>
          </Tooltip>

          {/* User Menu (Desktop) */}
          {user && (
            <div className="hidden md:block relative ml-1.5 pl-3 border-l border-(--border)" ref={dropdownRef}>
              <button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className="flex items-center gap-2 hover:opacity-70 transition-opacity"
              >
                <div className="w-8 h-8 rounded-full bg-(--muted)/10 border border-(--border) flex items-center justify-center overflow-hidden transition-all hover:border-(--accent)">
                  {user.avatar ? (
                    <img src={user.avatar} alt={user.username} className="w-full h-full object-cover" />
                  ) : (
                    <User className="w-4 h-4 text-(--muted)" />
                  )}
                </div>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className={`text-(--muted) transition-transform duration-200 ${isDropdownOpen ? "rotate-180" : ""}`}
                >
                  <path d="m6 9 6 6 6-6" />
                </svg>
              </button>

              {/* Dropdown Menu */}
              {isDropdownOpen && (
                <div className="absolute right-0 top-full mt-3 w-52 rounded-2xl border border-(--border) bg-(--background) shadow-xl shadow-black/5 animate-in fade-in zoom-in-95 duration-200">
                  <div className="p-1.5">
                    <Link
                      href={`/profile/${user.user_id}`}
                      className="flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl hover:bg-(--muted)/10 transition-colors text-(--foreground)"
                      onClick={() => setIsDropdownOpen(false)}
                    >
                      <User className="w-4 h-4 text-(--muted)" />
                      Profile
                    </Link>
                    <Link
                      href={`/profile/${user.user_id}/settings`}
                      className="flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl hover:bg-(--muted)/10 transition-colors text-(--foreground)"
                      onClick={() => setIsDropdownOpen(false)}
                    >
                      <Settings className="w-4 h-4 text-(--muted)" />
                      Settings
                    </Link>
                    <div className="h-px bg-(--border) my-1.5" />
                    <button
                      className="w-full flex items-center gap-3 px-3.5 py-2.5 text-sm font-medium rounded-xl text-red-500 hover:bg-red-500/10 transition-colors text-left"
                      onClick={() => {
                        logoutClient();
                        setIsDropdownOpen(false);
                      }}
                    >
                      <LogOut className="w-4 h-4" />
                      Sign Out
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
          {/* Mobile Menu Toggle */}
          <button
            className="md:hidden p-2 text-(--foreground) hover:bg-(--muted)/10 rounded-full transition-colors"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Navigation Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-(--border) bg-(--background) animate-in slide-in-from-top-2">
          <div className="p-4 space-y-1.5">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-colors ${active
                    ? "bg-(--accent)/10 text-(--accent) font-medium"
                    : "text-(--muted) hover:bg-(--muted)/10 hover:text-(--foreground)"
                    }`}
                >
                  <Icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}

            {user && (
              <>
                <div className="h-px bg-(--border) my-2" />
                <Link
                  href={`/profile/${user.user_id}`}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-xl text-(--muted) hover:bg-(--muted)/10 hover:text-(--foreground) transition-colors"
                >
                  <User className="w-5 h-5" />
                  Profile
                </Link>
                <Link
                  href={`/profile/${user.user_id}/settings`}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-xl text-(--muted) hover:bg-(--muted)/10 hover:text-(--foreground) transition-colors"
                >
                  <Settings className="w-5 h-5" />
                  Settings
                </Link>
                <button
                  className="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-red-500 hover:bg-red-500/10 transition-colors text-left"
                  onClick={() => logoutClient()}
                >
                  <LogOut className="w-5 h-5" />
                  Sign Out
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
