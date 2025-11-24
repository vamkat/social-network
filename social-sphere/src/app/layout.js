import "./globals.css";
import { AuthProvider } from "@/providers/AuthProvider";

export const metadata = {
  title: "Social Sphere - Home",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className="bg-slate-50 text-slate-900">
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
