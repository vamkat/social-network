"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { Home, Users, MessageCircle, Bell, User, LogOut, Settings, Menu, X, Boxes, Search } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import Tooltip from "@/components/ui/tooltip";
import { getUserByID } from "@/mock-data/users";

export default function Navbar() {
  const pathname = usePathname();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);

  // Mock user data - replace with actual auth context later
  const user = getUserByID("1");

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen(false);
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
      icon: Home,
    },
    {
      label: "Friends",
      href: "/feed/friends",
      icon: Users,
    },
    {
      label: "Groups",
      href: "/groups",
      icon: Boxes,
    },
  ];

  const isActive = (path) => pathname === path;

  return (
    <nav className="sticky top-0 z-50 w-full border-b border-black/5 dark:border-white/5 bg-(--background)/80 backdrop-blur-md">
      <div className="max-w px-6 h-16 flex items-center justify-between gap-4">
        {/* Left: Logo */}
        <div className="shrink-0 w-48">
          <Link href="/feed/public" className="flex items-center gap-2 font-semibold text-lg tracking-tight hover:opacity-80 transition-opacity">
            <Image
              src="/logo.png"
              alt="SocialSphere Logo"
              width={32}
              height={32}
              className="w-8 h-8"
            />
            <span className="hidden md:block">SocialSphere</span>
          </Link>
        </div>

        {/* Center: Search Bar */}
        <div className="hidden lg:flex flex-1 max-w-md mx-4 min-w-0">
          <div className="relative w-full group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-4 w-4 text-(--muted) group-focus-within:text-(--foreground) transition-colors" />
            </div>
            <input
              type="text"
              className="block w-full pl-10 pr-3 py-2 border border-transparent rounded-full leading-5 bg-(--muted)/10 text-(--foreground) placeholder-(--muted) focus:outline-none focus:bg-(--background) focus:ring-2 focus:ring-(--foreground)/10 focus:border-(--muted)/20 sm:text-sm transition-all duration-200"
              placeholder="Search users..."
            />
          </div>
        </div>

        {/* Right: Nav Links & Actions */}
        <div className="flex items-center justify-end gap-2 w-48 shrink-0">

          {/* Desktop Nav Links */}
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.href);
            return (
              <Tooltip key={item.href} content={item.label}>
                <Link
                  href={item.href}
                  className={`relative p-2.5 rounded-xl flex items-center justify-center transition-all duration-200 group ${active
                    ? "text-(--foreground) bg-(--muted)/10"
                    : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/5"
                    }`}
                >
                  <Icon className={`w-6 h-6 ${active ? "stroke-[2.5px]" : "stroke-2"}`} />
                  {active && (
                    <span className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full bg-(--foreground)" />
                  )}
                </Link>
              </Tooltip>
            );
          })}


          {/* Messages */}
          <Tooltip content="Messages">
            <button
              className={`p-2.5 rounded-xl transition-all duration-200 relative ${isActive('/messages')
                ? "text-(--foreground) bg-(--muted)/10"
                : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10"}`}
            >
              <MessageCircle className={`w-6 h-6 ${isActive('/messages') ? "stroke-[2.5px]" : "stroke-2"}`} />
              <span className="absolute top-0 right-0.5 text-xs text-red-500 font-extrabold flex items-center justify-center rounded-full">1</span>
            </button>
          </Tooltip>

          {/* Notifications */}
          <Tooltip content="Notifications">
            <button
              className={`p-2.5 rounded-xl transition-all duration-200 relative ${isActive('/notifications')
                ? "text-(--foreground) bg-(--muted)/10"
                : "text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10"}`}
            >
              <Bell className={`w-6 h-6 ${isActive('/notifications') ? "stroke-[2.5px]" : "stroke-2"}`} />
              <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-(--background)" />
            </button>
          </Tooltip>

          {/* User Menu (Desktop) */}
          <div className="hidden md:block relative ml-2" ref={dropdownRef}>
            <button
              onClick={() => setIsDropdownOpen(!isDropdownOpen)}
              className="flex items-center gap-2 pl-2 border-l border-(--muted)/20 hover:opacity-80 transition-opacity"
            >
              <div className="w-8 h-8 rounded-full bg-(--muted)/20 flex items-center justify-center overflow-hidden">
                {user.Avatar ? (
                  <img src={user.Avatar} alt={user.Username} className="w-full h-full object-cover" />
                ) : (
                  <User className="w-4 h-4 text-(--muted)" />
                )}
              </div>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
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
              <div className="absolute right-0 top-full mt-2 w-48 rounded-xl border border-black/5 dark:border-white/5 bg-(--background) shadow-lg shadow-black/5 animate-in fade-in zoom-in-95 duration-200">
                <div className="p-1">
                  <Link
                    href={`/profile/${user.ID}`}
                    className="flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-lg hover:bg-(--muted)/10 transition-colors"
                    onClick={() => setIsDropdownOpen(false)}
                  >
                    <User className="w-4 h-4 text-(--muted)" />
                    Profile
                  </Link>
                  <Link
                    href={`/profile/${user.ID}/settings`}
                    className="flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-lg hover:bg-(--muted)/10 transition-colors"
                    onClick={() => setIsDropdownOpen(false)}
                  >
                    <Settings className="w-4 h-4 text-(--muted)" />
                    Settings
                  </Link>
                  <div className="h-px bg-(--muted)/10 my-1" />
                  <button
                    className="w-full flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-lg text-red-500 hover:bg-red-500/5 transition-colors text-left"
                    onClick={() => setIsDropdownOpen(false)}
                  >
                    <LogOut className="w-4 h-4" />
                    Sign Out
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* Mobile Menu Toggle */}
          <button
            className="md:hidden p-2 text-(--foreground)"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Navigation Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-black/5 dark:border-white/5 bg-(--background) animate-in slide-in-from-top-2">
          <div className="p-4 space-y-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-colors ${active
                    ? "bg-(--muted)/10 text-(--foreground) font-medium"
                    : "text-(--muted) hover:bg-(--muted)/5 hover:text-(--foreground)"
                    }`}
                >
                  <Icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}
            <div className="h-px bg-(--muted)/10 my-2" />
            <Link
              href={`/profile/${user.ID}`}
              onClick={() => setIsMobileMenuOpen(false)}
              className="flex items-center gap-3 px-4 py-3 rounded-xl text-(--muted) hover:bg-(--muted)/5 hover:text-(--foreground) transition-colors"
            >
              <User className="w-5 h-5" />
              Profile
            </Link>
            <Link
              href={`/profile/${user.ID}/settings`}
              onClick={() => setIsMobileMenuOpen(false)}
              className="flex items-center gap-3 px-4 py-3 rounded-xl text-(--muted) hover:bg-(--muted)/5 hover:text-(--foreground) transition-colors"
            >
              <Settings className="w-5 h-5" />
              Settings
            </Link>
            <button
              className="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-red-500 hover:bg-red-500/5 transition-colors text-left"
            >
              <LogOut className="w-5 h-5" />
              Sign Out
            </button>
          </div>
        </div>
      )}
    </nav>
  );
}
