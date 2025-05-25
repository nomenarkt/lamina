'use client';

import Image from 'next/image';

export default function EmailConfirmed() {
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

        <div className="text-brand-green text-5xl">✅</div>

        <h1 className="text-2xl font-semibold text-gray-800">
          Your email has been confirmed!
        </h1>

        <p className="text-gray-600">
          You can now log in to your account using your credentials.
        </p>

        <a
          href="/login"
          className="inline-block w-full text-center bg-brand-green hover:bg-emerald-700 text-white py-3 rounded-md font-medium transition"
        >
          Log In
        </a>
      </div>
    </div>
  );
}
