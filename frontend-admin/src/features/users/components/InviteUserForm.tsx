'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { inviteUserSchema } from '../schemas/inviteUserSchema';
import { useInviteUser } from '../hooks/useInviteUser';
import { InviteUserPayload } from '../types';
import { z } from 'zod';

type FormValues = z.infer<typeof inviteUserSchema>;

export function InviteUserForm() {
  const { mutate, status, isSuccess, error } = useInviteUser();

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(inviteUserSchema),
    defaultValues: { role: 'viewer' },
  });

  const isExternal = watch('role') === 'external';

  const onSubmit = (data: FormValues) => {
    mutate(data as InviteUserPayload);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <input {...register('email')} placeholder="Email" />
      {errors.email && <p>{errors.email.message}</p>}

      <select {...register('role')}>
        <option value="admin">Admin</option>
        <option value="viewer">Viewer</option>
        <option value="external">External</option>
      </select>

      {isExternal && (
        <>
          <input {...register('company')} placeholder="Company" />
          <div className="flex gap-2">
            <input type="datetime-local" {...register('accessDuration.from')} />
            <input type="datetime-local" {...register('accessDuration.to')} />
          </div>
        </>
      )}

      <button type="submit" disabled={status === 'pending'}>
        {status === 'pending' ? 'Sending...' : 'Send Invitation'}
      </button>

      {isSuccess && <p>Invitation sent!</p>}
      {error && (
        <p className="text-red-500">
          {error instanceof Error ? error.message : 'Something went wrong'}
        </p>
      )}
    </form>
  );
}
