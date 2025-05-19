import { InviteUserPayload } from './types';

export async function inviteUser(payload: InviteUserPayload) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE}/api/v1/admin/create-user`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
    });


  if (!res.ok) throw new Error('Failed to invite user');
  return res.json();
}
