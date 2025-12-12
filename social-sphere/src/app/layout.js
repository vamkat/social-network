import "./globals.css";
import NextAuthProvider from "@/components/providers/NextAuthProvider";

export const metadata = {
  title: "Social Sphere - Home",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <NextAuthProvider>
          {children}
        </NextAuthProvider>
      </body>
    </html>
  );
}
