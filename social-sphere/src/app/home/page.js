// app/home/page.tsx
import Image from "next/image";
import Link from "next/link";

export default function HomePage() {
  return (
    <div className="page-container">
      {/* Top Navigation */}
      <nav className="nav">
        <div className="nav-content flex items-center justify-between gap-4">
          {/* Logo */}
          <Link href="/" className="nav-logo flex items-center gap-2">
            <Image
              src="/logos.png"
              alt="SocialSphere Logo"
              width={32}
              height={32}
              className="w-8 h-8"
            />
            <span className="font-semibold tracking-tight">SocialSphere</span>
          </Link>

          {/* Search (desktop only) */}
          <div className="hidden md:flex flex-1 max-w-md mx-4">
            <div className="w-full">
              <input
                type="text"
                className="w-full rounded-full border border-slate-200 dark:border-slate-800 px-4 py-2 text-sm bg-white/80 dark:bg-slate-900/80 focus:outline-none focus:ring-2 focus:ring-purple-500/40"
                placeholder="Search people, circles, posts…"
              />
            </div>
          </div>

          {/* Right actions */}
          <div className="nav-links flex items-center gap-3">
            <Link
              href="/home"
              className="hidden sm:inline-flex text-sm px-3 py-1 rounded-full border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/60 transition-colors"
            >
              Home
            </Link>
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 rounded-full border border-slate-200 dark:border-slate-800 flex items-center justify-center text-xs">
                📨
              </div>
              <div className="w-8 h-8 rounded-full border border-slate-200 dark:border-slate-800 flex items-center justify-center text-xs">
                🔔
              </div>
              <div className="w-9 h-9 rounded-full bg-purple-500/20 border border-purple-500/40 flex items-center justify-center text-xs font-semibold">
                SM
              </div>
            </div>
          </div>
        </div>
      </nav>

      {/* Main layout */}
      <main className="main-content mt-6">
        {/* FLEX LAYOUT: sidebar (left) + center + active members (right) */}
        <div className="max-w-6xl mx-auto w-full flex flex-col lg:flex-row gap-6">
          {/* LEFT SIDEBAR */}
          <aside className="w-full lg:w-64 shrink-0 space-y-4 animate-fade-in">
            {/* Profile card */}
            <section className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-purple-500/20 border border-purple-500/40 flex items-center justify-center text-sm font-semibold">
                  SM
                </div>
                <div>
                  <div className="font-semibold text-sm">Stam Manousis</div>
                  <div className="text-xs text-slate-500">@stam</div>
                </div>
              </div>
              <div className="mt-4 text-xs text-slate-500">
                Keeping it calm today ✨
              </div>
              <div className="mt-4 grid grid-cols-3 text-center text-xs">
                <div>
                  <div className="font-semibold">128</div>
                  <div className="text-slate-500">Posts</div>
                </div>
                <div>
                  <div className="font-semibold">342</div>
                  <div className="text-slate-500">Followers</div>
                </div>
                <div>
                  <div className="font-semibold">201</div>
                  <div className="text-slate-500">Following</div>
                </div>
              </div>
            </section>

            {/* Navigation / circles */}
            <section className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4 space-y-3 text-sm">
              <div className="font-semibold text-xs uppercase tracking-wide text-slate-500">
                Navigation
              </div>
              <div className="flex flex-col gap-1">
                <Link href="/home" className="px-2 py-1 rounded-md hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Home
                </Link>
                <Link href="/circles" className="px-2 py-1 rounded-md hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  My circles
                </Link>
                <Link href="/messages" className="px-2 py-1 rounded-md hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Messages
                </Link>
                <Link href="/notifications" className="px-2 py-1 rounded-md hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Notifications
                </Link>
                <Link href="/saved" className="px-2 py-1 rounded-md hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Saved
                </Link>
              </div>

              <div className="pt-3 border-t border-slate-100 dark:border-slate-800" />

              <div className="font-semibold text-xs uppercase tracking-wide text-slate-500">
                Pinned circles
              </div>
              <div className="flex flex-col gap-1">
                <button className="text-left px-2 py-1 rounded-full text-xs border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Design Circle
                </button>
                <button className="text-left px-2 py-1 rounded-full text-xs border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Close Friends
                </button>
                <button className="text-left px-2 py-1 rounded-full text-xs border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/80">
                  Study Group
                </button>
              </div>
            </section>
          </aside>

          {/* CENTER: COMPOSER + FEED */}
          <section className="flex-1 space-y-4 animate-fade-in lg:delay-100">
            {/* Composer */}
            <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4">
              <div className="flex gap-3">
                <div className="w-10 h-10 rounded-full bg-purple-500/20 border border-purple-500/40 flex items-center justify-center text-xs font-semibold">
                  SM
                </div>
                <div className="flex-1">
                  <div className="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-900/80 px-3 py-2 text-sm text-slate-500">
                    Share something with your circle…
                  </div>
                  <div className="mt-3 flex flex-wrap items-center gap-3 text-xs text-slate-500">
                    <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                      📷 <span>Add photo</span>
                    </div>
                    <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                      📊 <span>Create poll</span>
                    </div>
                    <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-slate-50 dark:bg-slate-900/80">
                      📅 <span>Start event</span>
                    </div>
                    <div className="ml-auto flex items-center gap-1 px-2 py-1 rounded-full bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300">
                      🔒 <span>Close friends</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Feed */}
            <div className="space-y-4">
              <div className="text-xs uppercase tracking-wide text-slate-500 px-1">
                For you
              </div>

              {/* Post 1 */}
              <article className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4 space-y-3">
                <header className="flex gap-3 items-start">
                  <div className="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-[10px]">
                    MK
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between gap-2">
                      <div>
                        <div className="text-sm font-semibold">Maria K.</div>
                        <div className="flex items-center gap-2 text-xs text-slate-500">
                          <span>@maria</span>
                          <span>·</span>
                          <span>2h ago</span>
                          <span className="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] bg-slate-50 dark:bg-slate-900/80 border border-slate-200/70 dark:border-slate-800/80">
                            Circle: Close Friends
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </header>

                <div className="text-sm text-slate-700 dark:text-slate-200">
                  Took a social detox this weekend and it felt amazing. Grateful
                  for this quiet space to share without the noise.
                </div>
                <div className="mt-2 rounded-xl bg-slate-100 dark:bg-slate-800 h-40 flex items-center justify-center text-xs text-slate-500">
                  Image placeholder
                </div>

                <footer className="flex items-center gap-6 text-xs text-slate-500 pt-2">
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    ❤️ <span>24</span>
                  </button>
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    💬 <span>8</span>
                  </button>
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    🔁 <span>Share</span>
                  </button>
                  <button className="ml-auto flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    📎 <span>Save</span>
                  </button>
                </footer>
              </article>

              {/* Post 2 */}
              <article className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4 space-y-3">
                <header className="flex gap-3 items-start">
                  <div className="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-[10px]">
                    SG
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between gap-2">
                      <div>
                        <div className="text-sm font-semibold">
                          Study Group
                        </div>
                        <div className="flex items-center gap-2 text-xs text-slate-500">
                          <span>Circle</span>
                          <span>·</span>
                          <span>5h ago</span>
                          <span className="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-200">
                            Upcoming meetup
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </header>

                <div className="text-sm text-slate-700 dark:text-slate-200">
                  Planning a quiet co-working evening on Friday. Cameras off,
                  lo-fi on, just focused energy. React and Go welcome.
                </div>

                <footer className="flex items-center gap-6 text-xs text-slate-500 pt-2">
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    ❤️ <span>12</span>
                  </button>
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    💬 <span>3</span>
                  </button>
                  <button className="flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    🔁 <span>Share</span>
                  </button>
                  <button className="ml-auto flex items-center gap-1 hover:text-slate-700 dark:hover:text-slate-200">
                    📎 <span>Save</span>
                  </button>
                </footer>
              </article>
            </div>
          </section>

          {/* RIGHT: ACTIVE MEMBERS (like FB chat sidebar) */}
          <aside className="hidden lg:flex w-56 shrink-0 flex-col space-y-4 animate-fade-in lg:delay-150">
            <section className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/70 p-4">
              <div className="flex items-center justify-between mb-3">
                <h2 className="text-sm font-semibold">Active now</h2>
                <span className="text-[10px] text-slate-500">Chat</span>
              </div>

              <div className="space-y-2 text-sm">
                {["Maria", "John", "Eleni", "Alex", "Nikos", "Dimitra"].map(
                  (name, idx) => (
                    <button
                      key={idx}
                      className="w-full flex items-center justify-between gap-2 px-2 py-1.5 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-900/70"
                    >
                      <div className="flex items-center gap-2">
                        <div className="relative w-8 h-8 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-[10px]">
                          <span>{name[0]}</span>
                          <span className="absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full bg-emerald-400 ring-2 ring-white dark:ring-slate-900" />
                        </div>
                        <span className="text-xs">{name}</span>
                      </div>
                      <span className="text-[10px] text-slate-400">
                        Online
                      </span>
                    </button>
                  )
                )}
              </div>
            </section>
          </aside>
        </div>
      </main>
    </div>
  );
}
