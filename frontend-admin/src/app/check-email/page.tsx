'use client';

import Image from 'next/image';
import { useSearchParams } from 'next/navigation';
import ResendConfirmation from '@/components/ui/ResendConfirmation';

export default function CheckEmail() {
  const searchParams = useSearchParams();
  const email = searchParams.get('email') || '';

  const hasEmail = Boolean(email);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-white px-4">
      <div className="max-w-md w-full text-center space-y-6">
        <Image
          src="/logo.webp"
          alt="Madagascar Airlines"
          width={160}
          height={80}
          className="mx-auto h-auto"
          priority
        />

        <div className="text-brand-gold text-5xl">📨</div>

        <h1 className="text-2xl font-semibold text-gray-800">Check your inbox</h1>

        <p className="text-gray-600">
          {hasEmail ? (
            <>
              We&apos;ve sent a confirmation email to <strong>{email}</strong>.
              <br />
              It’s valid for <strong>24 hours</strong> — please click the link to activate your account.
            </>
          ) : (
            <>
              ❌ <strong>Email not found.</strong> Please return to{' '}
              <a href="/signup" className="underline text-brand-green">signup</a>.
            </>
          )}
        </p>

        {hasEmail && <ResendConfirmation email={email} />}
      </div>
    </div>
  );
}
