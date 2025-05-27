'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

function getMessageColor(message: string) {
  if (message.startsWith('✅')) return 'text-green-600';
  if (message.startsWith('⚠️')) return 'text-yellow-600';
  if (message.startsWith('❌')) return 'text-red-600';
  return 'text-gray-600';
}

interface ResendConfirmationProps {
  email: string;
}

export default function ResendConfirmation({ email }: ResendConfirmationProps) {
  const router = useRouter();

  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [disabled, setDisabled] = useState(false);

  const resend = async () => {
    if (!email) {
      setMessage('❌ Missing email. Cannot resend confirmation.');
      return;
    }

    setIsLoading(true);
    setMessage('');

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/resend-confirmation`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });

      if (res.ok) {
        setMessage('✅ We\'ve sent you a new confirmation email.');
        setDisabled(true);
        setTimeout(() => setDisabled(false), 60000);
      } else if (res.status === 400) {
        router.push('/login');
      } else if (res.status === 403) {
        setMessage('⚠️ This option is available only for Madagascar Airlines crew accounts.');
      } else if (res.status === 404) {
        router.push('/signup?resend=failed');
      } else {
        setMessage('❌ Could not resend. Please try again later.');
      }
    } catch {
      setMessage('❌ Something went wrong. Please try again later.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="mt-4 text-sm text-gray-600 space-y-2">
      <p>Didn’t get it?</p>
      <button
        onClick={resend}
        disabled={disabled || isLoading}
        className="mt-2 font-medium text-brand-green hover:underline focus:outline-none focus-visible:ring-2 focus-visible:ring-brand-green disabled:opacity-50"
      >
        {isLoading ? '⏳ Resending...' : 'Resend confirmation'}
      </button>
      {message && (
        <p
          className={`mt-2 text-sm ${getMessageColor(message)}`}
          role="alert"
          aria-live="polite"
        >
          {message}
        </p>
      )}
      <p className="mt-1 text-xs text-gray-400">Tip: It might be in your spam or junk folder.</p>
    </div>
  );
}
