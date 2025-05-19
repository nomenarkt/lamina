'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { login } from '../api';
import { useState } from 'react';

const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type LoginFormData = z.infer<typeof loginSchema>;

export function LoginForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const [serverError, setServerError] = useState('');
  const [isLoading, setLoading] = useState(false);

  const onSubmit = async (data: LoginFormData) => {
    setLoading(true);
    setServerError('');
    try {
      const res = await login(data);
      localStorage.setItem('access_token', res.access_token);
      localStorage.setItem('refresh_token', res.refresh_token);
      // TODO: Redirect based on role once decoded
      console.log('Login success!', res);
  } catch (err) {
    if (err instanceof Error) {
      setServerError(err.message);
    } else {
      setServerError('Unknown error occurred');
    }
  } finally {
    setLoading(false);
  }
};

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <input type="email" {...register('email')} placeholder="Email" className="input" />
      {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}

      <input type="password" {...register('password')} placeholder="Password" className="input" />
      {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}

      <button type="submit" className="btn" disabled={isLoading}>
        {isLoading ? 'Logging in...' : 'Login'}
      </button>

      {serverError && <p className="text-red-600">{serverError}</p>}
    </form>
  );
}
