import type { Metadata } from "next";
import { siteConfig } from "@/config/site";
import { Navbar } from "@/components/navbar";
import "./globals.css";
import { UserProvider } from "@/components/user";
import { ToastProvider } from "@/components/toast-provider";

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
      <body suppressHydrationWarning={true} className="app-shell">
        <UserProvider>
          <ToastProvider />
          <div className="app-background" aria-hidden="true">
            <div className="app-orb app-orb-left" />
            <div className="app-orb app-orb-right" />
          </div>
          <div className="relative flex min-h-screen flex-col">
            <Navbar />
            <main className="mx-auto w-full max-w-7xl flex-grow px-4 pb-10 pt-6 sm:px-6 lg:px-8">
              {children}
            </main>
            <footer className="mx-auto flex w-full max-w-7xl items-center justify-center gap-2 px-4 py-6 text-sm text-base-content/60 sm:px-6 lg:px-8">
              <span>Powered by</span>
              <a
                className="link link-hover flex items-center gap-1 font-semibold text-current"
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
