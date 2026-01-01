import { HeroHeader } from "@/components/HeroHeader";
import Link from "next/link";

export default function LandingPage() {
    return (
        <div className="page-container ">
            {/* Navigation */}
            <nav className="section-border">
                <div className="max-w-7xl mx-auto px-6 py-6">
                    <div className="flex items-center justify-between">
                        <Link href="/" className="flex items-center gap-3">
                            <span className="text-base font-medium tracking-tight">SocialSphere</span>
                        </Link>
                        <div className="flex items-center gap-8">
                            <Link
                                href="/login"
                                className="text-sm link-muted"
                                prefetch={false}
                            >
                                Sign In
                            </Link>
                            <Link
                                href="/register"
                                className="link-primary text-sm"
                            >
                                Get Started
                            </Link>
                            <Link
                                href="/about"
                                className="link-primary text-sm"
                            >
                                About Us
                            </Link>
                        </div>
                    </div>
                </div>
            </nav>

            {/* Hero Section */}
            <section className="section-border bg-(--muted)/10">
                {/* <div className="max-w-7xl mx-auto px-6 py-32 md:py-26">
                    <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 xl:gap-20 items-center">
 
                        <div className="relative max-w-2xl">


                            <h1 className="heading-xl mb-12">
                                Connect<br />
                                <span className="relative inline-block text-[120px] font-black">
                                    <span className="absolute inset-0 blur-sm opacity-40 translate-y-2">without</span>
                                    <span className="text-image-fill relative">without</span>
                                </span> <br />
                                the noise
                            </h1>

                            <div className="flex flex-col sm:flex-row items-start gap-8 sm:gap-10">
                                <p className="text-xl text-muted leading-relaxed max-w-md">
                                    A social network designed for meaningful conversations.
                                    No algorithms. No ads. Just real connections.
                                </p>
                                <div className="flex flex-col gap-3 pt-1 shrink-0">
                                    <Link
                                        href="/register"
                                        className="text-sm link-accent whitespace-nowrap"
                                    >
                                        Start your journey →
                                    </Link>
                                    <Link
                                        href="/login"
                                        className="text-sm link-muted whitespace-nowrap"
                                    >
                                        Learn more
                                    </Link>
                                </div>
                            </div>
                        </div>

                        <div className="hidden lg:flex justify-center items-center">
                            <img
                                src="/check.png"
                                alt="SocialSphere"
                                className="w-full max-w-lg h-auto"
                            />
                        </div>
                    </div>
                </div> */}
                <HeroHeader/>
            </section>


            {/* <section className="section-border">
                <div className="max-w-7xl mx-auto px-6 py-24">
                    <div className="grid md:grid-cols-3 gap-x-12 gap-y-20">

                        <div>
                            <div className="text-label mb-4">
                                01 / Real-time
                            </div>
                            <h3 className="heading-sm mb-4">
                                Conversations that flow
                            </h3>
                            <p className="text-muted leading-relaxed">
                                See when friends are active and share moments instantly.
                                No lag, no clutter, just pure connection.
                            </p>
                        </div>


                        <div>
                            <div className="text-label mb-4">
                                02 / Privacy
                            </div>
                            <h3 className="heading-sm mb-4">
                                Your data, your rules
                            </h3>
                            <p className="text-muted leading-relaxed">
                                You decide who sees what. Share with the world or just close friends—always your choice.
                            </p>
                        </div>

                        <div>
                            <div className="text-label mb-4">
                                03 / Groups
                            </div>
                            <h3 className="heading-sm mb-4">
                                Your inner circle
                            </h3>
                            <p className="text-muted leading-relaxed">
                                Create intimate spaces for the people who matter most. Plan events, share updates, stay close.
                            </p>
                        </div>
                    </div>
                </div>
            </section>
            

            <section className="section-border bg-(--muted)/10">
                <div className="max-w-7xl mx-auto px-6 py-32">
                    <div className="max-w-4xl">
                        <h2 className="heading-lg mb-8">
                            Ready to connect<br />differently?
                        </h2>
                        <p className="text-xl text-muted mb-10 max-w-2xl">
                            Join people who value real connections over vanity metrics
                        </p>
                        <Link
                            href="/register"
                            className="inline-block text-base link-accent underline underline-offset-8 decoration-2"
                        >
                            Create your account →
                        </Link>
                    </div>
                </div>
            </section>


            <footer>
                <div className="max-w-7xl mx-auto px-6 py-12">
                    <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
                        <div className="text-sm text-neutral-500">
                            © 2025 SocialSphere
                        </div>
                        <div className="flex gap-8 text-sm font-medium">
                            <Link href="#" className="link-muted">Privacy</Link>
                            <Link href="#" className="link-muted">Terms</Link>
                            <Link href="#" className="link-muted">Contact</Link>
                        </div>
                    </div>
                </div>
            </footer> */}
        </div>
    );
}
