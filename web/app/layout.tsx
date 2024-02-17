import type { Metadata } from "next";
import { siteConfig } from "@/config/site";
import { Navbar } from "@/components/navbar";
import "./globals.css";
import { UserProvider } from "@/components/user";

export const metadata: Metadata = {
  title: siteConfig.name,
  description: siteConfig.description,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head />
      <body>
        <UserProvider>
          <div className="relative flex flex-col h-screen">
            <Navbar />
            <main className="container mx-auto max-w-7xl px-6 flex-grow">
              {children}
            </main>
            <footer className="w-full flex items-center justify-center py-3 gap-2">
              <span className="text-default-600">Powered by</span>
              <a
                className="link link-hover flex items-center gap-1 text-current"
                href="https://github.com/TimothyYe/godns"
                title="GoDNS project homepage"
                target="_blank"
              >
                <p className="text-primary">GoDNS</p>
              </a>
            </footer>
          </div>
        </UserProvider>
      </body>
    </html>
  );
}
