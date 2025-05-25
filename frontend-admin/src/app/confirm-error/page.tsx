'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import Image from 'next/image';

export default function ConfirmErrorPage() {
  const searchParams = useSearchParams();
  const [message, setMessage] = useState('Something went wrong.');
  const [action, setAction] = useState({ href: '/signup', label: 'Return to Signup' });

  useEffect(() => {
    const reason = searchParams.get('reason');

    if (reason === 'expired') {
      setMessage('Your confirmation link has expired.');
      setAction({ href: '/signup', label: 'Return to Signup' });
    } else if (reason === 'already-confirmed') {
      setMessage('This account has already been confirmed.');
      setAction({ href: '/login', label: 'Return to Login' });
    } else if (reason === 'invalid') {
      setMessage('Invalid or malformed confirmation link.');
      setAction({ href: '/signup', label: 'Return to Signup' });
    }
  }, [searchParams]);

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

        <div className="text-red-600 text-5xl">❌</div>

        <h1 className="text-2xl font-semibold text-gray-800">
          Email Confirmation Failed
        </h1>

        <p className="text-gray-600">
          {message}
        </p>

        <a
          href={action.href}
          className="inline-block w-full text-center bg-brand-green hover:bg-emerald-700 text-white py-3 rounded-md font-medium transition"
        >
          {action.label}
        </a>
      </div>
    </div>
  );
}
