'use client';

import Image from 'next/image';
import { ReactNode } from 'react';

type Props = {
  title: string;
  children: ReactNode;
  footer: ReactNode;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void | Promise<void>;
};

export function AuthCardLayout({ title, children, footer, onSubmit }: Props) {
  return (
    <form
      onSubmit={onSubmit}
      role="form"
      className="w-full max-w-md p-8 bg-white rounded-2xl shadow-xl space-y-6 min-h-[500px]"
    >
      <div className="flex justify-center">
        <Image
          src="/logo.webp"
          alt="Madagascar Airlines"
          width={160}
          height={60}
          priority
        />
      </div>

      <h1 className="text-2xl font-semibold text-center text-brand-grey-dark">
        {title}
      </h1>

      <div className="grid gap-4">{children}</div>

      {footer}
    </form>
  );
}
