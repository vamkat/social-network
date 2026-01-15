"use client";

import { LiveSocketProvider } from "@/context/LiveSocketContext";

export default function LiveSocketWrapper({ children }) {
    return <LiveSocketProvider>{children}</LiveSocketProvider>;
}
