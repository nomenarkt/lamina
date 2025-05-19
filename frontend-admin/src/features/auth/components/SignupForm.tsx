'use client';

import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { signup } from '../api';

const schema = z
  .object({
    email: z
      .string({ required_error: 'Email is required' })
      .email('Invalid email address'),
    password: z
      .string({ required_error: 'Password is required' })
      .min(8, 'Password must be at least 8 characters'),
    confirm_password: z.string({ required_error: 'Confirm your password' }),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: 'Passwords do not match',
    path: ['confirm_password'],
  });

type SignupFormData = z.infer<typeof schema>;

export function SignupForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignupFormData>({
    resolver: zodResolver(schema),
    mode: 'onBlur', // ✅ Enables field-level validation on blur
  });

  const [serverError, setServerError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (data: SignupFormData) => {
    if (!data.email.endsWith('@madagascarairlines.com')) {
      setServerError('Email must be from @madagascarairlines.com');
      return;
    }

    setLoading(true);
    setServerError('');
    try {
      await signup({ email: data.email, password: data.password });
      setSuccess(true);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : 'Signup failed';
      setServerError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <p className="text-green-600">
        Check your email to confirm your account.
      </p>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <input
          type="email"
          placeholder="Email"
          aria-describedby="email-error"
          aria-invalid={!!errors.email}
          {...register('email')}
        />
        {errors.email && (
          <p role="alert" className="text-red-500" id="email-error">
            {errors.email.message}
          </p>
        )}
      </div>

      <div>
        <input
          type="password"
          placeholder="Password"
          aria-describedby="password-error"
          aria-invalid={!!errors.password}
          {...register('password')}
        />
        {errors.password && (
          <p role="alert" className="text-red-500" id="password-error">
            {errors.password.message}
          </p>
        )}
      </div>

      <div>
        <input
          type="password"
          placeholder="Confirm Password"
          aria-describedby="confirm-error"
          aria-invalid={!!errors.confirm_password}
          {...register('confirm_password')}
        />
        {errors.confirm_password && (
          <p role="alert" className="text-red-500" id="confirm-error">
            {errors.confirm_password.message}
          </p>
        )}
      </div>

      <button type="submit" disabled={loading}>
        {loading ? 'Signing up...' : 'Signup'}
      </button>

      {serverError && <p className="text-red-600">{serverError}</p>}
    </form>
  );
}
