'use client';

import type { Metadata } from "next";
import "./globals.css";
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '@/lib/query-client';
import { ErrorBoundary } from '@/components/common/error-boundary';

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className="antialiased">
        <ErrorBoundary>
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
