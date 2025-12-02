import Link from "next/link";
import LoginForm from "@/components/forms/LoginForm";

export default function LoginPage() {
    return (
        <div className="page-container justify-center items-center">
            <div className="auth-card">
                <div className="text-center mb-10">
                    <h1 className="auth-title">Welcome Back</h1>
                    <p className="auth-subtitle">
                        Happy to see you again
                    </p>
                </div>

                <LoginForm />

                <p className="auth-footer">
                    Don't have an account?{" "}
                    <Link
                        href="/register"
                        className="link-primary"
                    >
                        Sign up
                    </Link>
                </p>
            </div>
        </div>
    );
}
