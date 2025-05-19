'use client';

import { ReactQueryProvider } from '@/app/providers/queryClient';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <ReactQueryProvider>{children}</ReactQueryProvider>;
}
