import Link from "next/link";
import RegisterForm from "@/components/forms/RegisterForm";

export default function RegisterPage() {
    return (
        <div className="page-container justify-center items-center py-12">
            <div className="auth-card max-w-lg">
                <div className="text-center mb-8">
                    <h1 className="auth-title">Create Account</h1>
                    <p className="auth-subtitle">
                        Join SocialSphere and start connecting
                    </p>
                </div>

                <RegisterForm />

                <p className="auth-footer">
                    Already have an account?{" "}
                    <Link
                        href="/login"
                        className="link-primary"
                    >
                        Sign in
                    </Link>
                </p>
            </div>
        </div>
    );
}
