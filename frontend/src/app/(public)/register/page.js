import RegistrationForm from "@/components/ui/RegistrationForm";

const RegisterPage = () => {
  return (
    <section className="min-h-screen bg-linear-to-b from-(--color-bg) via-white to-(--color-bg)">
      <div className="mx-auto flex w-full max-w-5xl flex-col items-center px-6 py-16 text-center">
        <div className="space-y-4">
          <p className="text-sm font-semibold tracking-wide text-(--color-border-accent)">
            Join SocialSphere
          </p>
          <div className="space-y-2">
            <h1 className="text-4xl font-extrabold leading-tight text-(--color-text)">
              Your <span className="text-(--color-accent)">community</span> starts here.
            </h1>
            <p className="text-lg text-(--color-text-muted)">
              Create a profile and unlock circles, chats, and events curated around what inspires you.
            </p>
          </div>
        </div>

        <div className="mt-12 w-full max-w-3xl">
          <RegistrationForm />
        </div>
      </div>
    </section>
  );
};

export default RegisterPage;
