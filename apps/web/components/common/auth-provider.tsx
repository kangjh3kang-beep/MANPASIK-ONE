'use client';

import React from 'react';
import { SessionProvider } from 'next-auth/react';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  return (
    // @ts-expect-error — React 19 + next-auth beta children type mismatch
    <SessionProvider>{children}</SessionProvider>
  );
}
