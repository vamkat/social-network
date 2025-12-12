import Navbar from "@/components/layout/navbar";

export default function MainLayout({ children }) {
    return (
        <div className="min-h-screen flex flex-col bg-(--muted)/6">
            <Navbar />
            <main className="flex-1 w-full">
                {children}
            </main>
        </div>
    );
}
