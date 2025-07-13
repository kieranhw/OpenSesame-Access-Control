import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Image from "next/image";
import { ThemeProvider } from "next-themes";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "OpenSesame Management",
  description: "OpenSesame Access Control Interface",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className="h-screen w-screen overflow-hidden"
    >
      <body
        className={` ${geistSans.variable} ${geistMono.variable} flex h-full w-full flex-col antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <header className="border-divider bg-card flex-none border-b px-4 py-2">
            <div className="flex items-center gap-3 h-8">
              <Image
                src="/sesame.png"
                alt="OpenSesame Logo"
                width={30}
                height={30}
                draggable="false"
              />
              <h1 className="text-foreground font-semibold">
                OpenSesame Access Control
              </h1>
            </div>
          </header>

          <main className="flex-1 overflow-auto">{children}</main>
        </ThemeProvider>
      </body>
    </html>
  );
}
