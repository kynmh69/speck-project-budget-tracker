import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "プロジェクト予算管理システム",
  description: "ソフトウェア開発プロジェクトの工数予実管理と収支管理",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.Node;
}>) {
  return (
    <html lang="ja">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
