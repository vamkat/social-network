import Navbar from "@/components/layout/navbar";

export default function MainLayout({ children }) {
    return (
        <div className="min-h-screen flex flex-col bg-(--background)">
            <Navbar />
            <main className="flex-1 w-full max-w-[1000px] mx-auto px-6 py-8">
                {children}
            </main>
        </div>
    );
}
