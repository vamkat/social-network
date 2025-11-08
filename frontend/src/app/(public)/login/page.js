import LoginForm from "@/components/ui/LoginForm";

const LoginPage = () => {
  return (
    <section className="min-h-screen bg-linear-to-b from-(--color-bg) via-white to-(--color-bg)">
      <div className="mx-auto flex w-full max-w-5xl flex-col items-center px-6 py-16 text-center">
        <div className="space-y-4">
          <p className="text-sm font-semibold tracking-wide text-(--color-border-accent)">
            Welcome back
          </p>
          <div className="space-y-2">
            <h1 className="text-4xl font-extrabold leading-tight text-(--color-text)">
              Continue the <span className="text-(--color-accent)">conversation</span>.
            </h1>
            <p className="text-lg text-(--color-text-muted)">
              Log in to reconnect with your circles, pick up messages, and join today's events.
            </p>
          </div>
        </div>

        <div className="mt-12 w-full max-w-md">
          <LoginForm />
        </div>
      </div>
    </section>
  );
};

export default LoginPage;
