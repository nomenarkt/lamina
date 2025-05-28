'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { AuthCardLayout } from './AuthCardLayout';
import ResendConfirmation from '@/components/ui/ResendConfirmation';

export function LoginCard() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showResend, setShowResend] = useState(false);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setShowResend(false);
    setLoading(true);

    if (!email || !password) {
      setError('Please fill in all fields');
      setLoading(false);
      return;
    }

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const result = await res.json().catch(() => ({}));

      if (res.status === 403) {
        setError('Account not confirmed');
        setShowResend(true);
        return;
      }

      if (res.status === 401) {
        setError('Invalid email or password');
        return;
      }

      if (!res.ok) {
        throw new Error(result.error || 'Login failed');
      }

      localStorage.setItem('access_token', result.access_token);
      router.push('/dashboard');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'An unknown error occurred';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthCardLayout
      title="Login to Your Account"
      onSubmit={handleLogin}
      footer={
        <p className="text-center text-sm text-gray-600">
          Not a member?{' '}
          <Link href="/signup" className="text-brand-green underline">
            Sign up
          </Link>
        </p>
      }
    >
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          placeholder="you@example.com"
          onChange={(e) => setEmail(e.target.value)}
          className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-brand-green"
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium text-gray-700">
          Password
        </label>
        <input
          id="password"
          type="password"
          value={password}
          placeholder="••••••••"
          onChange={(e) => setPassword(e.target.value)}
          className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-brand-green"
        />
      </div>

      <div role="alert" aria-live="assertive" className="text-sm text-red-600 min-h-[20px] mt-1">
        {error || ''}
      </div>

      {showResend && (
        <div className="mt-4 border-t border-gray-200 pt-4">
          <ResendConfirmation email={email} />
        </div>
      )}

      <button
        type="submit"
        disabled={loading}
        aria-disabled={loading}
        className={`w-full h-12 mt-2 rounded-md font-medium text-white transition ${
          loading ? 'bg-brand-grey-dark cursor-not-allowed' : 'bg-brand-green hover:bg-emerald-700'
        }`}
      >
        {loading ? 'Loading…' : 'Log In'}
      </button>
    </AuthCardLayout>
  );
}
