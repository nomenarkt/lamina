'use client';

import Image from 'next/image';

export default function CheckEmail() {
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

        <div className="text-brand-gold text-5xl">📩</div>

        <h1 className="text-2xl font-semibold text-gray-800">
          Check your email to confirm your account
        </h1>

        <p className="text-gray-600">
          We&apos;ve sent a confirmation link to <strong>your email address</strong>.
          <br />
          The link is <strong>valid for 24 hours</strong> &mdash; please click it to activate your account.
        </p>
      </div>
    </div>
  );
}
