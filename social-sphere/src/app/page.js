import Link from "next/link";
import Image from "next/image";

export default function LandingPage() {
  return (
    <div className="page-container">
      {/* Navigation */}
      <nav className="nav">
        <div className="nav-content">
          <Link href="/" className="nav-logo flex items-center gap-2">
            <Image
              src="/logos.png"
              alt="SocialSphere Logo"
              width={32}
              height={32}
              className="w-8 h-8"
            />
            SocialSphere
          </Link>
          <div className="nav-links">
            <Link href="/login" className="nav-link">
              Sign In
            </Link>
            <Link href="/register" className="btn btn-primary">
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      <main className="main-content">
        {/* Hero Section */}
        <section className="hero animate-fade-in">
          <h1 className="hero-title">
            Connect without<br />the noise
          </h1>
          <p className="hero-subtitle">
            A sanctuary for your social life. Meaningful conversations,
            private circles, and a space that respects your peace of mind.
          </p>
          <Link href="/register" className="btn btn-primary text-base px-8 py-4">
            Start your journey
          </Link>
        </section>

        <div className="section-divider" />

        {/* Feature 1: Real-time */}
        <section className="section animate-fade-in delay-100">
          <div className="section-content">
            <span className="feature-label">Presence</span>
            <h2 className="feature-title">Feel the Connection</h2>
            <p className="feature-desc">
              Conversations that flow naturally, just like in real life.
              See when friends are active and share moments instantly,
              without the lag or the clutter.
            </p>
          </div>
          <div className="section-visual">
            {/* Abstract visual representation */}
            <div className="w-64 h-64 rounded-full bg-purple-50 dark:bg-purple-900/10 flex items-center justify-center">
              <div className="flex gap-2">
                <div className="w-8 h-8 rounded-full bg-purple-500/20" />
                <div className="w-8 h-8 rounded-full bg-purple-500/20 translate-y-4" />
                <div className="w-8 h-8 rounded-full bg-purple-500/20" />
              </div>
            </div>
          </div>
        </section>

        <div className="section-divider" />

        {/* Feature 2: Groups */}
        <section className="section flex-col-reverse md:flex-row-reverse animate-fade-in delay-200">
          <div className="section-content">
            <span className="feature-label">Circles</span>
            <h2 className="feature-title">Your Inner Circle</h2>
            <p className="feature-desc">
              Create intimate spaces for the people who matter most.
              Plan get-togethers, share private updates, and keep your
              close communities truly close.
            </p>
          </div>
          <div className="section-visual">
            <div className="w-64 h-64 rounded-full bg-blue-50 dark:bg-blue-900/10 flex items-center justify-center">
              <div className="w-32 h-32 rounded-full bg-blue-500/10 animate-pulse" />
            </div>
          </div>
        </section>

        <div className="section-divider" />

        {/* Feature 3: Privacy */}
        <section className="section animate-fade-in delay-300">
          <div className="section-content">
            <span className="feature-label">Privacy</span>
            <h2 className="feature-title">On Your Terms</h2>
            <p className="feature-desc">
              You decide who sees what. Whether you want to share with the world
              or just a few close friends, your privacy is always in your control.
            </p>
          </div>
          <div className="section-visual">
            <div className="w-64 h-64 rounded-full bg-emerald-50 dark:bg-emerald-900/10 flex items-center justify-center">
              <div className="w-24 h-32 rounded-2xl border-2 border-emerald-500/10 flex flex-col p-4 gap-2">
                <div className="w-8 h-8 rounded-full bg-emerald-500/10" />
                <div className="w-full h-2 rounded-full bg-emerald-500/10" />
                <div className="w-2/3 h-2 rounded-full bg-emerald-500/10" />
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="footer">
        <div className="footer-content">
          <div>© 2025 SocialSphere. All rights reserved.</div>
          <div className="flex gap-6">
            <Link href="#" className="hover:text-(--foreground) transition-colors">Privacy</Link>
            <Link href="#" className="hover:text-(--foreground) transition-colors">Terms</Link>
            <Link href="#" className="hover:text-(--foreground) transition-colors">Contact</Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
