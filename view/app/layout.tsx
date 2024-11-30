import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Blaze",
  description: "A Project to never see the day of light",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
