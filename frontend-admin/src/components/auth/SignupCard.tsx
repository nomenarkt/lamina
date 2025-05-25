'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { AuthCardLayout } from './AuthCardLayout';

export function SignupCard() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const isEmailValid = email.endsWith('@madagascarairlines.com');

  const handleSignup = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (!email || !password || !confirmPassword) {
        throw new Error('Please fill in all fields');
      }

      if (!isEmailValid) {
        throw new Error('Only @madagascarairlines.com emails are accepted');
      }

      if (password !== confirmPassword) {
        throw new Error('Passwords do not match');
      }

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json().catch(() => ({}));

      if (!res.ok) {
        throw new Error(data.message || 'Signup failed');
      }

      router.push('/check-email'); // ✅ Redirect here after successful signup
    } catch (err) {
      const message =
        err instanceof Error ? err.message : typeof err === 'string' ? err : 'An unknown error occurred';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthCardLayout
      title="Create Account"
      onSubmit={handleSignup}
      footer={
        <p className="text-center text-sm text-gray-600">
          Already have an account?{' '}
          <Link href="/login" className="text-brand-green underline">
            Login
          </Link>
        </p>
      }
    >
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-brand-green"
        />
        <p className="text-sm mt-1 text-gray-500">
          Only emails ending with <strong>@madagascarairlines.com</strong> are allowed.
        </p>
        {!isEmailValid && email.length > 0 && (
          <p className="text-sm text-red-600 mt-1" role="alert">
            Only @madagascarairlines.com emails are accepted.
          </p>
        )}
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium text-gray-700">Password</label>
        <input
          id="password"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-brand-green"
        />
      </div>

      <div>
        <label htmlFor="confirm-password" className="block text-sm font-medium text-gray-700">Confirm Password</label>
        <input
          id="confirm-password"
          type="password"
          placeholder="••••••••"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-brand-green"
        />
      </div>

      <div role="alert" aria-live="assertive" className="text-sm text-red-600 min-h-[20px] mt-1">
        {error || ''}
      </div>

      <button
        type="submit"
        disabled={loading}
        aria-disabled={loading}
        className={`w-full h-12 mt-2 rounded-md font-medium text-white transition ${
          loading ? 'bg-brand-grey-dark cursor-not-allowed' : 'bg-brand-green hover:bg-emerald-700'
        }`}
      >
        {loading ? 'Loading…' : 'Sign Up'}
      </button>
    </AuthCardLayout>
  );
}
