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
    <html lang="en" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var stored=localStorage.getItem('theme');var system=window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light';document.documentElement.setAttribute('data-theme',stored||system);}catch(e){document.documentElement.setAttribute('data-theme','light');}})();`,
          }}
        />
      </head>
      <body suppressHydrationWarning={true} className="app-shell">
        <UserProvider>
          <div className="app-backdrop" aria-hidden="true" />
          <div className="relative flex min-h-screen flex-col">
            <Navbar />
            <main className="container mx-auto max-w-7xl flex-grow px-4 sm:px-6">
              {children}
            </main>
            <footer className="w-full flex items-center justify-center gap-2 px-4 py-6 text-sm theme-faint">
              <span>Powered by</span>
              <a
                className="link link-hover flex items-center gap-1 theme-text"
                href="https://github.com/TimothyYe/godns"
                title="GoDNS project homepage"
                target="_blank"
              >
                <p className="font-medium theme-badge-sky px-2 py-1 rounded-full">GoDNS</p>
              </a>
            </footer>
          </div>
        </UserProvider>
      </body>
    </html>
  );
}
