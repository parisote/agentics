import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Agentics Configuration Generator",
  description: "Generate JSON configurations for Agentics - The Go LLM agent framework",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="bg-gradient-to-b from-blue-50 to-white min-h-screen">
          {children}
        </div>
      </body>
    </html>
  );
}
