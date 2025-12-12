import { getServerSession } from "next-auth";
import { authOptions } from "../../[...nextauth]/route";
import { NextResponse } from 'next/server';

export async function GET(request, { params }) {
    try {
        const { userId } = await params;
        const session = await getServerSession(authOptions);
        const cookieHeader = session?.backendCookie;

        const apiBase = process.env.API_BASE || "http://api-gateway:8081";

        const headers = {};
        if (cookieHeader) {
            headers['Cookie'] = cookieHeader;
        }

        const backendResponse = await fetch(`${apiBase}/profile/${userId}`, {
            method: "GET",
            headers: headers,
            next: { revalidate: 180}
        });

        if (!backendResponse.ok) {
            const errorData = await backendResponse.json().catch(() => null);
            return NextResponse.json(
                errorData || { error: "Failed to fetch profile" },
                { status: backendResponse.status }
            );
        }

        const profileData = await backendResponse.json();

        return NextResponse.json(profileData, { status: 200 });
    } catch (error) {
        console.error("Profile API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}
