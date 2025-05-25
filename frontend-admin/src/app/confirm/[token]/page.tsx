'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

type Props = {
  params: {
    token: string;
  };
};

export default function ConfirmPage({ params }: Props) {
  const router = useRouter();
  const { token } = params;
  const [status, setStatus] = useState<'pending' | 'success' | 'error'>('pending');

  useEffect(() => {
    const confirmEmail = async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/confirm/${token}`);
        if (res.ok) {
          setStatus('success');
          setTimeout(() => router.push('/email-confirmed'), 500);
        } else {
          setStatus('error');
          setTimeout(() => router.push('/confirm-error'), 1000);
        }
      } catch {
        setStatus('error');
        setTimeout(() => router.push('/confirm-error'), 1000);
      }
    };

    confirmEmail();
  }, [token, router]);

  return (
    <div className="min-h-screen flex items-center justify-center text-gray-600">
      {status === 'pending' && 'Confirming your email…'}
      {status === 'success' && 'Email confirmed. Redirecting…'}
      {status === 'error' && 'There was a problem confirming your email.'}
    </div>
  );
}
