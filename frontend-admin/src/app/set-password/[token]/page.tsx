'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useState } from 'react';

const schema = z
  .object({
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

type FormData = z.infer<typeof schema>;

export default function SetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');
  const [serverError, setServerError] = useState('');
  const [success, setSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    setServerError('');
    try {
      const res = await fetch('/api/auth/set-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          token,
          password: data.password,
        }),
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Something went wrong');
      }

      setSuccess(true);
      setTimeout(() => {
        router.push('/login');
      }, 2000);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setServerError(err.message);
      } else {
        setServerError('An unknown error occurred');
      }
    }
  };

  if (!token) {
    return <p className="text-red-500">Invalid or missing token.</p>;
  }

  if (success) {
    return <p className="text-green-600">Password set successfully. Redirecting to login...</p>;
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="max-w-md mx-auto space-y-4 p-6">
      <h1 className="text-2xl font-bold">Set Your Password</h1>

      <input
        type="password"
        placeholder="New Password"
        {...register('password')}
        className="w-full p-2 border rounded"
      />
      {errors.password && <p className="text-red-500">{errors.password.message}</p>}

      <input
        type="password"
        placeholder="Confirm Password"
        {...register('confirmPassword')}
        className="w-full p-2 border rounded"
      />
      {errors.confirmPassword && (
        <p className="text-red-500">{errors.confirmPassword.message}</p>
      )}

      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full bg-blue-600 text-white p-2 rounded"
      >
        {isSubmitting ? 'Submitting...' : 'Set Password'}
      </button>

      {serverError && <p className="text-red-600">{serverError}</p>}
    </form>
  );
}
