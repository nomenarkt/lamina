'use client';

import { ReactQueryProvider } from './queryClient';

export function ClientWrapper({ children }: { children: React.ReactNode }) {
  return <ReactQueryProvider>{children}</ReactQueryProvider>;
}
